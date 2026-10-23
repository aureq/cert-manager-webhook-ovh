package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cert-manager/cert-manager/pkg/acme/webhook/apis/acme/v1alpha1"
	"github.com/cert-manager/cert-manager/pkg/issuer/acme/dns/util"
	"github.com/ovh/go-ovh/ovh"
)

func (s *ovhDNSProviderSolver) validateProviderConfig(cfg *ovhDNSProviderConfig, allowAmbientCredentials bool) error {
	// log.V(logf.TraceLevel).Info("Entering validateProviderConfig() function")

	var err error

	if allowAmbientCredentials {
		// When allowAmbientCredentials is true, OVH client can load missing config
		// values from the environment variables and the ovh.conf files.
		// Ambient configuration should be located in ./ovh.conf, ~/.ovh.conf or /etc/ovh.conf
		// Note: This method isn't tested, and probably not widely used.
		// log.V(logf.WarnLevel).Info("Ambient credentials allowed, skipping config validation.")
		return nil
	}
	if cfg.Endpoint == "" {
		err = errors.New("no endpoint provided in OVH config")
		// log.V(logf.ErrorLevel).Error(err, "No endpoint provided in OVH config")
		return err
	}
	switch cfg.AuthenticationMethod {
	case "application":
		if cfg.ApplicationKeyRef.Name == "" {
			err = errors.New("no application key provided in OVH config")
			// log.V(logf.ErrorLevel).Error(err, "No application key provided in OVH config")
			return err
		}
		if cfg.ApplicationSecretRef.Name == "" {
			err = errors.New("no application secret provided in OVH config")
			// log.V(logf.ErrorLevel).Error(err, "No application secret provided in OVH config")
			return err
		}
		if cfg.ApplicationConsumerKeyRef.Name == "" {
			err = errors.New("no application consumer key provided in OVH config")
			// log.V(logf.ErrorLevel).Error(err, "No application consumer key provided in OVH config")
			return err
		}
	case "oauth2":
		if cfg.Oauth2ClientIDRef.Name == "" {
			err = errors.New("no oauth2 client ID provided in OVH config")
			// log.V(logf.ErrorLevel).Error(err, "No oauth2 client ID provided in OVH config")
			return err
		}
		if cfg.Oauth2ClientSecretRef.Name == "" {
			err = errors.New("no oauth2 client secret provided in OVH config")
			// log.V(logf.ErrorLevel).Error(err, "No oauth2 client secret provided in OVH config")
			return err
		}
	default:
		err = errors.New("invalid value for authentication method, allowed values: 'application' and 'oauth2'")
		// log.V(logf.ErrorLevel).Error(err, "Invalid value for authentication method, allowed values: 'application' and 'oauth2'", "ovhAuthenticationMethod", cfg.AuthenticationMethod)
		return err
	}

	// log.V(logf.DebugLevel).Info("Provider config: passed.")
	// log.V(logf.TraceLevel).Info("Exiting validateProviderConfig() function")
	return nil
}

// ovhClientApplication creates an OVH client using Application credentials.
// These credentials consist of an Application Key, Application Secret, and
// Consumer Key, which are used to authenticate API requests.
func (s *ovhDNSProviderSolver) ovhClientApplication(ch *v1alpha1.ChallengeRequest, cfg *ovhDNSProviderConfig) (*ovh.Client, error) {
	// log.V(logf.TraceLevel).Info("Entering ovhClientApplication() function")

	applicationKey, err := s.secret(cfg.ApplicationKeyRef, ch.ResourceNamespace)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving application key from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationKeyRef.Name, "key", cfg.ApplicationKeyRef.Key)
		return nil, err
	}

	applicationSecret, err := s.secret(cfg.ApplicationSecretRef, ch.ResourceNamespace)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving application secret from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationSecretRef.Name, "key", cfg.ApplicationSecretRef.Key)
		return nil, err
	}

	applicationConsumerKey, err := s.secret(cfg.ApplicationConsumerKeyRef, ch.ResourceNamespace)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving application consumer key from secret", "namespace", ch.ResourceNamespace, "secret", cfg.ApplicationConsumerKeyRef.Name, "key", cfg.ApplicationConsumerKeyRef.Key)
		return nil, err
	}

	client, err := ovh.NewClient(cfg.Endpoint, applicationKey, applicationSecret, applicationConsumerKey)
	// log.V(logf.TraceLevel).Info("Exiting ovhClientApplication() function")
	return client, err
}

// ovhClientOAuth2 creates an OVH client using OAuth2 credentials.
// These credentials consist of a Client ID and Client Secret, which are used
// to authenticate API requests.
func (s *ovhDNSProviderSolver) ovhClientOAuth2(ch *v1alpha1.ChallengeRequest, cfg *ovhDNSProviderConfig) (*ovh.Client, error) {
	// log.V(logf.TraceLevel).Info("Entering ovhClientOAuth2() function")

	oauth2ClientID, err := s.secret(cfg.Oauth2ClientIDRef, ch.ResourceNamespace)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving OAuth2 client ID from secret", "namespace", ch.ResourceNamespace, "secret", cfg.Oauth2ClientIDRef.Name, "key", cfg.Oauth2ClientIDRef.Key)
		return nil, err
	}

	oauth2ClientSecret, err := s.secret(cfg.Oauth2ClientSecretRef, ch.ResourceNamespace)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving OAuth2 client secret from secret", "namespace", ch.ResourceNamespace, "secret", cfg.Oauth2ClientSecretRef.Name, "key", cfg.Oauth2ClientSecretRef.Key)
		return nil, err
	}

	client, err := ovh.NewOAuth2Client(cfg.Endpoint, oauth2ClientID, oauth2ClientSecret)
	// log.V(logf.TraceLevel).Info("Exiting ovhClientOAuth2() function")
	return client, err
}

// ovhClient creates and returns an OVH client based on the provided configuration.
func (s *ovhDNSProviderSolver) ovhClient(ch *v1alpha1.ChallengeRequest) (*ovh.Client, error) {
	// log.V(logf.TraceLevel).Info("Entering ovhClient() function")
	cfg, err := loadConfig(ch.Config)
	if err != nil {
		return nil, err
	}

	// log.V(logf.TraceLevel).Info("Resource namespace", "namespace", ch.ResourceNamespace)

	err = s.validateProviderConfig(&cfg, ch.AllowAmbientCredentials)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to validate OVH config")
		return nil, err
	}

	switch cfg.AuthenticationMethod {
	case "application":
		// log.V(logf.TraceLevel).Info("Exiting ovhClient() function", "method", cfg.AuthenticationMethod)
		return s.ovhClientApplication(ch, &cfg)
	case "oauth2":
		// log.V(logf.TraceLevel).Info("Exiting ovhClient() function", "method", cfg.AuthenticationMethod)
		return s.ovhClientOAuth2(ch, &cfg)
	default:
		// log.V(logf.TraceLevel).Info("Exiting ovhClient() function", "method", cfg.AuthenticationMethod)
		return nil, fmt.Errorf("invalid value for authentication method, allowed values: 'application', 'oauth2'")
	}
}

// getSubDomain extracts the subdomain from the fully qualified domain name (FQDN)
// by removing the main domain part. If the main domain is not found in the FQDN,
// it returns the FQDN without the trailing dot.
func getSubDomain(domain, fqdn string) string {
	// log.V(logf.TraceLevel).Info("Entering getSubDomain() function", "domain", domain, "fqdn", fqdn)

	// Ensure both domain and fqdn are in unfqdn format for comparison
	domain = util.UnFqdn(domain)
	fqdn = util.UnFqdn(fqdn)

	// Check if fqdn ends with the domain
	if strings.HasSuffix(fqdn, "."+domain) {
		// Remove the domain part and the dot before it
		subDomain := strings.TrimSuffix(fqdn, "."+domain)
		// log.V(logf.TraceLevel).Info("Exiting getSubDomain() function", "subdomain", subDomain)
		return subDomain
	}

	// If the fqdn equals the domain exactly, return empty string
	if fqdn == domain {
		// log.V(logf.TraceLevel).Info("Exiting getSubDomain() function", "subdomain == fqdn", fqdn)
		return ""
	}

	// log.V(logf.TraceLevel).Info("Exiting getSubDomain() function", "subdomain", fqdn)
	return fqdn
}

// addTXTRecord adds a TXT record to the OVH DNS zone for the specified domain and subdomain.
func addTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	// log.V(logf.TraceLevel).Info("Entering addTXTRecord() function", "domain", domain, "subdomain", subDomain, "target", target)

	err := validateZone(ovhClient, domain)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to validate DNS zone", "domain", domain, "subdomain", subDomain)
		return err
	}

	_, err = createRecord(ovhClient, domain, "TXT", subDomain, target)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to create TXT record", "domain", domain, "subdomain", subDomain)
		return err
	}
	err = refreshRecords(ovhClient, domain)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to refresh DNS records", "domain", domain)
		return err
	}

	// log.V(logf.TraceLevel).Info("Exiting addTXTRecord() function", "domain", domain, "subdomain", subDomain, "target", target)
	return nil
}

func removeTXTRecord(ovhClient *ovh.Client, domain, subDomain, target string) error {
	// log.V(logf.TraceLevel).Info("Entering removeTXTRecord() function", "domain", domain, "subdomain", subDomain, "target", target)

	ids, err := listRecords(ovhClient, domain, "TXT", subDomain)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to list TXT records", "domain", domain, "subdomain", subDomain)
		return err
	}

	for _, id := range ids {
		record, err := getRecord(ovhClient, domain, id)
		if err != nil {
			// log.V(logf.ErrorLevel).Error(err, "Failed to get TXT record", "domain", domain, "subdomain", subDomain, "recordID", id)
			return err
		}
		// Remove surrounded quotes as OVH seems to implicitly add them for a TXT record
		if strings.Trim(record.Target, "\"") != target {
			// log.V(logf.WarnLevel).Info("Skipping deletion of TXT record as target does not match", "domain", domain, "subdomain", subDomain, "recordID", id, "recordTarget", record.Target, "expectedTarget", target)
			continue
		}
		err = deleteRecord(ovhClient, domain, id)
		if err != nil {
			// log.V(logf.ErrorLevel).Error(err, "Failed to delete TXT record", "domain", domain, "subdomain", subDomain, "recordID", id)
			return err
		}
	}

	err = refreshRecords(ovhClient, domain)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to refresh DNS records", "domain", domain)
		return err
	}

	// log.V(logf.TraceLevel).Info("Exiting removeTXTRecord() function", "domain", domain, "subdomain", subDomain, "target", target)
	return nil
}

// validateZone checks if the DNS zone for the given domain is deployed in OVH.
func validateZone(ovhClient *ovh.Client, domain string) error {
	// log.V(logf.TraceLevel).Info("Entering validateZone() function", "domain", domain)

	url := "/domain/zone/" + domain + "/status"
	zoneStatus := ovhZoneStatus{}

	err := ovhClient.Get(url, &zoneStatus)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to get DNS zone status", "domain", domain)
		return fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}
	if !zoneStatus.IsDeployed {
		// log.V(logf.ErrorLevel).Error(nil, "DNS zone not deployed", "domain", domain)
		return fmt.Errorf("DNS zone not deployed for domain %s", domain)
	}

	// log.V(logf.TraceLevel).Info("Exiting validateZone() function", "domain", domain)
	return nil
}

// listRecords retrieves the IDs of DNS records for the specified domain, field type, and subdomain.
func listRecords(ovhClient *ovh.Client, domain, fieldType, subDomain string) ([]int64, error) {
	// log.V(logf.TraceLevel).Info("Entering listRecords() function", "domain", domain, "fieldType", fieldType, "subDomain", subDomain)

	url := "/domain/zone/" + domain + "/record?fieldType=" + fieldType + "&subDomain=" + subDomain
	ids := []int64{}
	err := ovhClient.Get(url, &ids)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to list DNS records", "domain", domain, "fieldType", fieldType, "subDomain", subDomain, "url", url)
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}

	// log.V(logf.TraceLevel).Info("Exiting listRecords() function", "domain", domain, "fieldType", fieldType, "subDomain", subDomain, "ids", ids)
	return ids, nil
}

// getRecord retrieves the details of a DNS record for the specified domain and record ID.
func getRecord(ovhClient *ovh.Client, domain string, id int64) (*ovhZoneRecord, error) {
	// log.V(logf.TraceLevel).Info("Entering getRecord() function", "domain", domain, "recordID", id)

	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	record := ovhZoneRecord{}
	err := ovhClient.Get(url, &record)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to get DNS record", "domain", domain, "recordID", id, "url", url)
		return nil, fmt.Errorf("OVH API call failed: GET %s - %v", url, err)
	}

	// log.V(logf.TraceLevel).Info("Exiting getRecord() function", "domain", domain, "recordID", id)
	return &record, nil
}

func deleteRecord(ovhClient *ovh.Client, domain string, id int64) error {
	// log.V(logf.TraceLevel).Info("Entering deleteRecord() function", "domain", domain, "recordID", id)

	url := "/domain/zone/" + domain + "/record/" + strconv.FormatInt(id, 10)
	err := ovhClient.Delete(url, nil)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to delete DNS record", "domain", domain, "recordID", id, "url", url)
		return fmt.Errorf("OVH API call failed: DELETE %s - %v", url, err)
	}

	// log.V(logf.TraceLevel).Info("Exiting deleteRecord() function", "domain", domain, "recordID", id)
	return nil
}

func createRecord(ovhClient *ovh.Client, domain, fieldType, subDomain, target string) (*ovhZoneRecord, error) {
	// log.V(logf.TraceLevel).Info("Entering createRecord() function", "domain", domain, "fieldType", fieldType, "subDomain", subDomain, "target", target)

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
		// log.V(logf.ErrorLevel).Error(err, "Failed to create DNS record", "domain", domain, "subDomain", subDomain, "url", url)
		return nil, fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	// log.V(logf.TraceLevel).Info("Exiting createRecord() function", "domain", domain, "fieldType", fieldType, "subDomain", subDomain, "target", target)
	return &record, nil
}

func refreshRecords(ovhClient *ovh.Client, domain string) error {
	// log.V(logf.TraceLevel).Info("Entering refreshRecords() function", "domain", domain)

	url := "/domain/zone/" + domain + "/refresh"
	err := ovhClient.Post(url, nil, nil)
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Failed to refresh DNS records", "domain", domain, "url", url)
		return fmt.Errorf("OVH API call failed: POST %s - %v", url, err)
	}

	// log.V(logf.TraceLevel).Info("Exiting refreshRecords() function", "domain", domain)
	return nil
}
