package main

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/acme/webhook/cmd"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	logf "github.com/cert-manager/cert-manager/pkg/logs"
)

var (
	log       = logf.Log.WithName("webhook-ovh")
	GroupName = os.Getenv("GROUP_NAME")
)

// ovhDNSProviderSolver implements the provider-specific logic needed to
// 'present' an ACME challenge TXT record for your own DNS provider.
// To do so, it must implement the `github.com/jetstack/cert-manager/pkg/acme/webhook.Solver`
// interface.
type ovhDNSProviderSolver struct {
	client kubernetes.Interface
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

	if GroupName == "" {
		klog.Fatal("'GROUP_NAME' environment variable must be specified")
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

// Present is responsible for actually presenting the DNS record with the
// DNS provider.
// This method should tolerate being called multiple times with the same value.
// cert-manager itself will later perform a self check to ensure that the
// solver has correctly configured the DNS provider.
func (s *ovhDNSProviderSolver) Present(ch *v1alpha1.ChallengeRequest) error {
	log.V(logf.TraceLevel).Info("Entering Present() function")
	ovhClient, err := s.ovhClient(ch)
	if err != nil {
		log.V(logf.ErrorLevel).Error(err, "Failed to create OVH client")
		return err
	}
	domain := util.UnFqdn(ch.ResolvedZone)
	subDomain := getSubDomain(domain, ch.ResolvedFQDN)
	target := ch.Key
	log.V(logf.DebugLevel).Info("Prepared present parameters", "domain", domain, "subDomain", subDomain, "target", target)
	err = addTXTRecord(ovhClient, domain, subDomain, target)
	log.V(logf.TraceLevel).Info("Exiting Present() function")
	return err
}

// CleanUp should delete the relevant TXT record from the DNS provider console.
// If multiple TXT records exist with the same record name (e.g.
// _acme-challenge.example.com) then **only** the record with the same `key`
// value provided on the ChallengeRequest should be cleaned up.
// This is in order to facilitate multiple DNS validations for the same domain
// concurrently.
func (s *ovhDNSProviderSolver) CleanUp(ch *v1alpha1.ChallengeRequest) error {
	log.V(logf.TraceLevel).Info("Entering CleanUp() function")

	ovhClient, err := s.ovhClient(ch)
	if err != nil {
		log.V(logf.ErrorLevel).Error(err, "Failed to create OVH client")
		return err
	}

	domain := util.UnFqdn(ch.ResolvedZone)
	subDomain := getSubDomain(domain, ch.ResolvedFQDN)
	target := ch.Key

	log.V(logf.DebugLevel).Info("Prepared cleanup parameters", "domain", domain, "subDomain", subDomain, "target", target)

	err = removeTXTRecord(ovhClient, domain, subDomain, target)
	log.V(logf.TraceLevel).Info("Exiting CleanUp() function")
	return err
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
	log.V(logf.TraceLevel).Info("Entering Initialize() function")

	client, err := kubernetes.NewForConfig(kubeClientConfig)
	if err != nil {
		log.V(logf.ErrorLevel).Error(err, "Failed to create Kubernetes clientset")
		return err
	}

	s.client = client
	log.V(logf.TraceLevel).Info("Exiting Initialize() function")
	return nil
}
