package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/charmbracelet/log"
	"github.com/ovh/go-ovh/ovh"
)

var (
	logger    *log.Logger
	once      sync.Once
	GroupName string = os.Getenv("GROUP_NAME")
)

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
	Oauth2ClientIDRef         corev1.SecretKeySelector `json:"oauth2ClientIDRef"`
	Oauth2ClientSecretRef     corev1.SecretKeySelector `json:"oauth2ClientSecretRef"`
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

func main() {
	getLogger().Info("Starting OVH DNS webhook...")

	if GroupName == "" {
		getLogger().Fatal("'GROUP_NAME' environment variable must be specified")
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
	getLogger().Info("Validating provider config...")

	if allowAmbientCredentials {
		// When allowAmbientCredentials is true, OVH client can load missing config
		// values from the environment variables and the ovh.conf files.
		getLogger().Info("Ambient credentials allowed, skipping config validation.")
		return nil
	}
	if cfg.Endpoint == "" {
		getLogger().Error("No endpoint provided in OVH config")
		return errors.New("no endpoint provided in OVH config")
	}
	switch cfg.AuthenticationMethod {
	case "application":
		if cfg.ApplicationKeyRef.Name == "" {
			getLogger().Error("No application key provided in OVH config")
			return errors.New("no application key provided in OVH config")
		}
		if cfg.ApplicationSecretRef.Name == "" {
			getLogger().Error("No application secret provided in OVH config")
			return errors.New("no application secret provided in OVH config")
		}
		if cfg.ApplicationConsumerKeyRef.Name == "" {
			getLogger().Error("No application consumer key provided in OVH config")
			return errors.New("no application consumer key provided in OVH config")
		}
	case "oauth2":
		if cfg.Oauth2ClientIDRef.Name == "" {
			getLogger().Error("No oauth2 client ID provided in OVH config")
			return errors.New("no oauth2 client ID provided in OVH config")
		}
		if cfg.Oauth2ClientSecretRef.Name == "" {
			getLogger().Error("No oauth2 client secret provided in OVH config")
			return errors.New("no oauth2 client secret provided in OVH config")
		}
	default:
		getLogger().Error("Invalid value for authentifaction method, allowed values: 'application' and 'oauth2'")
		return errors.New("invalid value for authentifaction method, allowed values: 'application' and 'oauth2'")
	}
	getLogger().Info("Provider config: passed.")
	return nil
}

// ovhClientApplication creates an OVH client using Application credentials.
// These credentials consist of an Application Key, Application Secret, and
// Consumer Key, which are used to authenticate API requests.
func (s *ovhDNSProviderSolver) ovhClientApplication(ch *v1alpha1.ChallengeRequest, cfg *ovhDNSProviderConfig) (*ovh.Client, error) {
	applicationKey, err := s.secret(cfg.ApplicationKeyRef, ch.ResourceNamespace)
	if err != nil {
		getLogger().Error("Error retrieving application key from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationKeyRef.Name, "key", cfg.ApplicationKeyRef.Key, "error", err)
		return nil, err
	}

	applicationSecret, err := s.secret(cfg.ApplicationSecretRef, ch.ResourceNamespace)
	if err != nil {
		getLogger().Error("Error retrieving application secret from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationSecretRef.Name, "key", cfg.ApplicationSecretRef.Key, "error", err)
		return nil, err
	}

	applicationConsumerKey, err := s.secret(cfg.ApplicationConsumerKeyRef, ch.ResourceNamespace)
	if err != nil {
		getLogger().Error("Error retrieving application consumer key from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationConsumerKeyRef.Name, "key", cfg.ApplicationConsumerKeyRef.Key, "error", err)
		return nil, err
	}

	return ovh.NewClient(cfg.Endpoint, applicationKey, applicationSecret, applicationConsumerKey)
}

// ovhClientOAuth2 creates an OVH client using OAuth2 credentials.
// These credentials consist of a Client ID and Client Secret, which are used
// to authenticate API requests.
func (s *ovhDNSProviderSolver) ovhClientOAuth2(ch *v1alpha1.ChallengeRequest, cfg *ovhDNSProviderConfig) (*ovh.Client, error) {
	oauth2ClientID, err := s.secret(cfg.Oauth2ClientIDRef, ch.ResourceNamespace)
	if err != nil {
		getLogger().Error("Error retrieving OAuth2 client ID from secret", "namespace", ch.ResourceNamespace, "secret", cfg.Oauth2ClientIDRef.Name, "key", cfg.Oauth2ClientIDRef.Key, "error", err)
		return nil, err
	}

	oauth2ClientSecret, err := s.secret(cfg.Oauth2ClientSecretRef, ch.ResourceNamespace)
	if err != nil {
		getLogger().Error("Error retrieving OAuth2 client secret from secret", "namespace", ch.ResourceNamespace, "secret", cfg.Oauth2ClientSecretRef.Name, "key", cfg.Oauth2ClientSecretRef.Key, "error", err)
		return nil, err
	}

	return ovh.NewOAuth2Client(cfg.Endpoint, oauth2ClientID, oauth2ClientSecret)
}

// ovhClient creates and returns an OVH client based on the provided configuration.
func (s *ovhDNSProviderSolver) ovhClient(ch *v1alpha1.ChallengeRequest) (*ovh.Client, error) {
	getLogger().Info("Starting challenge request...")
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}

	getLogger().Info(fmt.Sprintf("Resource namespace: %s", ch.ResourceNamespace))

	err = s.validate(&cfg, ch.AllowAmbientCredentials)
	if err != nil {
		getLogger().Error("Failed to validate OVH config", "error", err)
		return nil, err
	}

	switch cfg.AuthenticationMethod {
	case "application":
		return s.ovhClientApplication(ch, &cfg)
	case "oauth2":
		return s.ovhClientOAuth2(ch, &cfg)
	default:
		return nil, fmt.Errorf("invalid value for authentifaction method, allowed values: application, oauth2'")
	}
}

func (s *ovhDNSProviderSolver) secret(ref corev1.SecretKeySelector, namespace string) (string, error) {
	if ref.Name == "" {
		getLogger().Error("No secret name specified in secret key selector")
		return "", fmt.Errorf("no secret name specified in secret key selector")
	}

	secret, err := s.client.CoreV1().Secrets(namespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
	if err != nil {
		getLogger().Error("Error retrieving secret", "namespace", namespace, "secret", ref.Name, "error", err)
		return "", err
	}

	bytes, ok := secret.Data[ref.Key]
	if !ok {
		getLogger().Error("Key not found in secret", "namespace", namespace, "secret", ref.Name, "key", ref.Key)
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
		getLogger().Error("Failed to create OVH client", "error", err)
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
		getLogger().Error("Failed to create OVH client", "error", err)
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
		getLogger().Error("Failed to create Kubernetes clientset", "error", err)
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
		getLogger().Error("Failed to decode OVH config", "error", err)
		return cfg, fmt.Errorf("error decoding OVH config: %v", err)
	}

	return cfg, nil
}

// getSubDomain extracts the subdomain from the fully qualified domain name (FQDN)
// by removing the main domain part. If the main domain is not found in the FQDN,
// it returns the FQDN without the trailing dot.
func getSubDomain(domain, fqdn string) string {
	if idx := strings.Index(fqdn, "."+domain); idx != -1 {
		return fqdn[:idx]
	}

	return util.UnFqdn(fqdn)
}

// addTXTRecord adds a TXT record to the OVH DNS zone for the specified domain and subdomain.
func addTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	err := validateZone(ovhClient, domain)
	if err != nil {
		getLogger().Error("Failed to validate DNS zone", "domain", domain, "subdomain", subDomain, "error", err)
		return err
	}

	_, err = createRecord(ovhClient, domain, "TXT", subDomain, target)
	if err != nil {
		getLogger().Error("Failed to create TXT record", "domain", domain, "subdomain", subDomain, "error", err)
		return err
	}
	return refreshRecords(ovhClient, domain)
}

func removeTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	ids, err := listRecords(ovhClient, domain, "TXT", subDomain)
	if err != nil {
		getLogger().Error("Failed to list TXT records", "domain", domain, "subdomain", subDomain, "error", err)
		return err
	}

	for _, id := range ids {
		record, err := getRecord(ovhClient, domain, id)
		if err != nil {
			getLogger().Error("Failed to get TXT record", "domain", domain, "subdomain", subDomain, "recordID", id, "error", err)
			return err
		}
		// Remove surrounded quotes as OVH seems to implicitly add them for a TXT record
		if strings.Trim(record.Target, "\"") != target {
			continue
		}
		err = deleteRecord(ovhClient, domain, id)
		if err != nil {
			getLogger().Error("Failed to delete TXT record", "domain", domain, "subdomain", subDomain, "recordID", id, "error", err)
			return err
		}
	}

	return refreshRecords(ovhClient, domain)
}

// validateZone checks if the DNS zone for the given domain is deployed in OVH.
func validateZone(ovhClient *ovh.Client, domain string) error {
	url := "/domain/zone/" + domain + "/status"
	zoneStatus := ovhZoneStatus{}
	err := ovhClient.Get(url, &zoneStatus)
	if err != nil {
		getLogger().Error("Failed to get DNS zone status", "domain", domain, "error", err)
		return fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	if !zoneStatus.IsDeployed {
		getLogger().Error("DNS zone not deployed", "domain", domain)
		return fmt.Errorf("DNS zone not deployed for domain %s", domain)
	}

	return nil
}

// listRecords retrieves the IDs of DNS records for the specified domain, field type, and subdomain.
func listRecords(ovhClient *ovh.Client, domain, fieldType, subDomain string) ([]int64, error) {
	url := "/domain/zone/" + domain + "/record?fieldType=" + fieldType + "&subDomain=" + subDomain
	ids := []int64{}
	err := ovhClient.Get(url, &ids)
	if err != nil {
		getLogger().Error("Failed to list DNS records", "domain", domain, "fieldType", fieldType, "subDomain", subDomain, "url", url, "error", err)
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	return ids, nil
}

// getRecord retrieves the details of a DNS record for the specified domain and record ID.
func getRecord(ovhClient *ovh.Client, domain string, id int64) (*ovhZoneRecord, error) {
	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	record := ovhZoneRecord{}
	err := ovhClient.Get(url, &record)
	if err != nil {
		getLogger().Error("Failed to get DNS record", "domain", domain, "recordID", id, "url", url, "error", err)
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	return &record, nil
}

func deleteRecord(ovhClient *ovh.Client, domain string, id int64) error {
	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	err := ovhClient.Delete(url, nil)
	if err != nil {
		getLogger().Error("Failed to delete DNS record", "domain", domain, "recordID", id, "url", url, "error", err)
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
		getLogger().Error("Failed to create DNS record", "domain", domain, "subDomain", subDomain, "url", url, "error", err)
		return nil, fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	return &record, nil
}

func refreshRecords(ovhClient *ovh.Client, domain string) error {
	url := "/domain/zone/" + domain + "/refresh"
	err := ovhClient.Post(url, nil, nil)
	if err != nil {
		getLogger().Error("Failed to refresh DNS records", "domain", domain, "url", url, "error", err)
		return fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	return nil
}

// getLogger initializes and returns a singleton logger instance
// with custom configuration options
func getLogger() *log.Logger {
	once.Do(func() {
		logger = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			Level:           log.InfoLevel,
		})
	})
	return logger
}
