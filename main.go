package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	logf "github.com/cert-manager/cert-manager/pkg/logs"
	"github.com/ovh/go-ovh/ovh"
)

var GroupName = os.Getenv("GROUP_NAME")

func main() {
	logf.Log.Info("Webhook starting...")

	if GroupName == "" {
		panic("GROUP_NAME must be specified")
	}

	// This will register our ovh DNS provider with the webhook serving
	// library, making it available as an API under the provided GroupName.
	// You can register multiple DNS provider implementations with a single
	// webhook, where the Name() method will be used to disambiguate between
	// the different implementations.
	cmd.RunWebhookServer(GroupName,
		&ovhDNSProviderSolver{},
	)
}

// ovhDNSProviderSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/jetstack/cert-manager/pkg/acme/webhook.Solver`
// interface.
type ovhDNSProviderSolver struct {
	client *kubernetes.Clientset
}

// ovhDNSProviderConfig is a structure that is used to decode into when
// solving a DNS01 challenge.
// This information is provided by cert-manager, and may be a reference to
// additional configuration that's needed to solve the challenge for this
// particular certificate or issuer.
// This typically includes references to Secret resources containing DNS
// provider credentials, in cases where a 'multi-tenant' DNS solver is being
// created.
// If you do *not* require per-issuer or per-certificate configuration to be
// provided to your webhook, you can skip decoding altogether in favour of
// using CLI flags or similar to provide configuration.
// You should not include sensitive information here. If credentials need to
// be used by your provider here, you should reference a Kubernetes Secret
// resource and fetch these credentials using a Kubernetes clientset.
type ovhDNSProviderConfig struct {
	Endpoint                  string                   `json:"endpoint"`
	AuthenticationMethod      string                   `json:"authenticationMethod"`
	ApplicationKeyRef         corev1.SecretKeySelector `json:"applicationKeyRef"`
	ApplicationSecretRef      corev1.SecretKeySelector `json:"applicationSecretRef"`
	ApplicationConsumerKeyRef corev1.SecretKeySelector `json:"applicationConsumerKeyRef"`
	OAuth2ClientIdRef         corev1.SecretKeySelector `json:"oauth2ClientIdRef"`
	OAuth2ClientSecretRef     corev1.SecretKeySelector `json:"oauth2ClientSecretRef"`
}

type ovhZoneStatus struct {
	IsDeployed bool `json:"isDeployed"`
}

type ovhZoneRecord struct {
	Id        int64  `json:"id,omitempty"`
	FieldType string `json:"fieldType"`
	SubDomain string `json:"subDomain"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl,omitempty"`
}

// Name is used as the name for this DNS solver when referencing it on the ACME
// Issuer resource.
// This should be unique **within the group name**, i.e. you can have two
// solvers configured with the same Name() **so long as they do not co-exist
// within a single webhook deployment**.
// For example, `cloudflare` may be used as the name of a solver.
func (s *ovhDNSProviderSolver) Name() string {
	return "ovh"
}

func (s *ovhDNSProviderSolver) validate(cfg *ovhDNSProviderConfig, allowAmbientCredentials bool) error {
	logf.Log.Info("Validating provider config...")

	if allowAmbientCredentials {
		// When allowAmbientCredentials is true, OVH client can load missing config
		// values from the environment variables and the ovh.conf files.
		return nil
	}
	if cfg.Endpoint == "" {
		return errors.New("no endpoint provided in OVH config")
	}

	switch cfg.AuthenticationMethod {
	case "application":
		if cfg.ApplicationKeyRef.Name == "" {
			return errors.New("no application key provided in OVH config")
		}
		if cfg.ApplicationSecretRef.Name == "" {
			return errors.New("no application secret provided in OVH config")
		}
		if cfg.ApplicationConsumerKeyRef.Name == "" {
			return errors.New("no consumer key provided in OVH config")
		}
	case "oauth2":
		if cfg.OAuth2ClientIdRef.Name == "" {
			return errors.New("no client ID provided in OVH config")
		}
		if cfg.OAuth2ClientSecretRef.Name == "" {
			return errors.New("no client secret provided in OVH config")
		}
	default:
		return errors.New("invalid authentication method provided in OVH config")
	}
	logf.Log.Info("Provider config: passed.")
	return nil
}

func (s *ovhDNSProviderSolver) ovhClient(ch *v1alpha1.ChallengeRequest) (*ovh.Client, error) {
	logf.Log.Info("Starting challenge request...")
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}

	logf.Log.Info(fmt.Sprintf("Resource namespace: %s", ch.ResourceNamespace))

	err = s.validate(&cfg, ch.AllowAmbientCredentials)
	if err != nil {
		return nil, err
	}

	switch cfg.AuthenticationMethod {
	case "application":
		applicationKey, err := s.secret(cfg.ApplicationKeyRef, ch.ResourceNamespace)
		if err != nil {
			return nil, err
		}

		applicationSecret, err := s.secret(cfg.ApplicationSecretRef, ch.ResourceNamespace)
		if err != nil {
			return nil, err
		}

		applicationConsumerKey, err := s.secret(cfg.ApplicationConsumerKeyRef, ch.ResourceNamespace)
		if err != nil {
			return nil, err
		}

		return ovh.NewClient(cfg.Endpoint, applicationKey, applicationSecret, applicationConsumerKey)
	case "oauth2":
		clientId, err := s.secret(cfg.OAuth2ClientIdRef, ch.ResourceNamespace)
		if err != nil {
			return nil, err
		}
		clientSecret, err := s.secret(cfg.OAuth2ClientSecretRef, ch.ResourceNamespace)
		if err != nil {
			return nil, err
		}
		return ovh.NewOAuth2Client(cfg.Endpoint, clientId, clientSecret)
	}

	return nil, fmt.Errorf("unknown authentication method %q", cfg.AuthenticationMethod)
}

func (s *ovhDNSProviderSolver) secret(ref corev1.SecretKeySelector, namespace string) (string, error) {
	if ref.Name == "" {
		return "", nil
	}

	secret, err := s.client.CoreV1().Secrets(namespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	bytes, ok := secret.Data[ref.Key]
	if !ok {
		return "", fmt.Errorf("key not found %q in secret '%s/%s'", ref.Key, namespace, ref.Name)
	}
	return strings.TrimSuffix(string(bytes), "\n"), nil
}

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (s *ovhDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	ovhClient, err := s.ovhClient(ch)
	if err != nil {
		return err
	}
	domain := util.UnFqdn(ch.ResolvedZone)
	subDomain := getSubDomain(domain, ch.ResolvedFQDN)
	target := ch.Key
	return addTXTRecord(ovhClient, domain, subDomain, target)
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (s *ovhDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	ovhClient, err := s.ovhClient(ch)
	if err != nil {
		return err
	}
	domain := util.UnFqdn(ch.ResolvedZone)
	subDomain := getSubDomain(domain, ch.ResolvedFQDN)
	target := ch.Key
	return removeTXTRecord(ovhClient, domain, subDomain, target)
}

// Initialize will be called when the webhook first starts.
// This method can be used to instantiate the webhook, i.e. initialising
// connections or warming up caches.
// Typically, the kubeClientConfig parameter is used to build a Kubernetes
// client that can be used to fetch resources from the Kubernetes API, e.g.
// Secret resources containing credentials used to authenticate with DNS
// provider accounts.
// The stopCh can be used to handle early termination of the webhook, in cases
// where a SIGTERM or similar signal is sent to the webhook process.
func (s *ovhDNSProviderSolver) Initialize(kubeClientConfig *rest.Config, stopCh <-chan struct{}) error {
	client, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		return err
	}

	s.client = client
	return nil
}

// loadConfig is a small helper function that decodes JSON configuration into
// the typed config struct.
func loadConfig(cfgJSON *extapi.JSON) (ovhDNSProviderConfig, error) {
	cfg := ovhDNSProviderConfig{}
	// handle the 'base case' where no configuration has been provided
	if cfgJSON == nil {
		return cfg, nil
	}
	if err := json.Unmarshal(cfgJSON.Raw, &cfg); err != nil {
		return cfg, fmt.Errorf("error decoding OVH config: %v", err)
	}

	return cfg, nil
}

func getSubDomain(domain, fqdn string) string {
	if idx := strings.Index(fqdn, "."+domain); idx != -1 {
		return fqdn[:idx]
	}

	return util.UnFqdn(fqdn)
}

func addTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	err := validateZone(ovhClient, domain)
	if err != nil {
		return err
	}

	_, err = createRecord(ovhClient, domain, "TXT", subDomain, target)
	if err != nil {
		return err
	}
	return refreshRecords(ovhClient, domain)
}

func removeTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	ids, err := listRecords(ovhClient, domain, "TXT", subDomain)
	if err != nil {
		return err
	}

	for _, id := range ids {
		record, err := getRecord(ovhClient, domain, id)
		if err != nil {
			return err
		}
		if record.Target != target {
			continue
		}
		err = deleteRecord(ovhClient, domain, id)
		if err != nil {
			return err
		}
	}

	return refreshRecords(ovhClient, domain)
}

func validateZone(ovhClient *ovh.Client, domain string) error {
	url := "/domain/zone/" + domain + "/status"
	zoneStatus := ovhZoneStatus{}
	err := ovhClient.Get(url, &zoneStatus)
	if err != nil {
		return fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	if !zoneStatus.IsDeployed {
		return fmt.Errorf("OVH zone not deployed for domain %s", domain)
	}

	return nil
}

func listRecords(ovhClient *ovh.Client, domain, fieldType, subDomain string) ([]int64, error) {
	url := "/domain/zone/" + domain + "/record?fieldType=" + fieldType + "&subDomain=" + subDomain
	ids := []int64{}
	err := ovhClient.Get(url, &ids)
	if err != nil {
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	return ids, nil
}

func getRecord(ovhClient *ovh.Client, domain string, id int64) (*ovhZoneRecord, error) {
	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	record := ovhZoneRecord{}
	err := ovhClient.Get(url, &record)
	if err != nil {
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	return &record, nil
}

func deleteRecord(ovhClient *ovh.Client, domain string, id int64) error {
	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	err := ovhClient.Delete(url, nil)
	if err != nil {
		return fmt.Errorf("OVH API call failed: DELETE %s - %v", url, err)
	}
	return nil
}

func createRecord(ovhClient *ovh.Client, domain, fieldType, subDomain, target string) (*ovhZoneRecord, error) {
	url := "/domain/zone/" + domain + "/record"
	params := ovhZoneRecord{
		FieldType: fieldType,
		SubDomain: subDomain,
		Target:    target,
		TTL:       60,
	}
	record := ovhZoneRecord{}
	err := ovhClient.Post(url, &params, &record)
	if err != nil {
		return nil, fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	return &record, nil
}

func refreshRecords(ovhClient *ovh.Client, domain string) error {
	url := "/domain/zone/" + domain + "/refresh"
	err := ovhClient.Post(url, nil, nil)
	if err != nil {
		return fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	return nil
}
