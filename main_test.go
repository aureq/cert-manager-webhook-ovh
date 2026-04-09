package main

import (
	"encoding/json"
	"testing"

	// acmetest "github.com/cert-manager/cert-manager/test/acme"
	corev1 "k8s.io/api/core/v1"
)

func TestValidate(t *testing.T) {
	s := &ovhDNSProviderSolver{}

	// fixture := acmetest.NewFixture(&ovhDNSProviderSolver{},
	// 	acmetest.SetResolvedZone(zone),
	// 	acmetest.SetAllowAmbientCredentials(false),
	// 	acmetest.SetManifestPath("testdata/ovh"),
	// )

	//need to uncomment and  RunConformance delete runBasic and runExtended once https://github.com/cert-manager/cert-manager/pull/4835 is merged
	// fixture.RunConformance(t)
	// fixture.RunBasic(t)
	// fixture.RunExtended(t)

	// No endpoint provided
	err := s.validate(&ovhDNSProviderConfig{}, false)
	if err == nil || err.Error() != "no endpoint provided in OVH config" {
		t.Errorf("expected missing endpoint error, got %v", err)
	}

	// No application key
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "application",
	}, false)
	if err == nil || err.Error() != "no application key provided in OVH config" {
		t.Errorf("expected application key error, got %v", err)
	}

	// No application secret
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "application",
		ApplicationKeyRef:    corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationKey"}},
	}, false)
	if err == nil || err.Error() != "no application secret provided in OVH config" {
		t.Errorf("expected application secret error, got %v", err)
	}

	// No consumer key nor application consumer key
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "application",
		ApplicationKeyRef:    corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationKey"}},
		ApplicationSecretRef: corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationSecret"}},
	}, false)
	if err == nil || err.Error() != "no application consumer key provided in OVH config" {
		t.Errorf("expected consumer key error, got %v", err)
	}

	// No OAuth2 Client ID
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "oauth2",
	}, false)
	if err == nil || err.Error() != "no oauth2 client ID provided in OVH config" {
		t.Errorf("expected oauth2 clientID error, got %v", err)
	}

	// No OAuth2 Client Secret
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "oauth2",
		Oauth2ClientIDRef:    corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "oauth2ClientId"}},
	}, false)
	if err == nil || err.Error() != "no oauth2 client secret provided in OVH config" {
		t.Errorf("expected oauth2 client secret error, got %v", err)
	}

	// Invalid Authentication Method
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "foo",
	}, false)
	if err == nil || err.Error() != "invalid value for authentifaction method, allowed values: 'application' and 'oauth2'" {
		t.Errorf("expected invalid method error, got %v", err)
	}

	// Valid configuration for Application
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:                  "ovh-eu",
		AuthenticationMethod:      "application",
		ApplicationKeyRef:         corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationKey"}},
		ApplicationSecretRef:      corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationSecret"}},
		ApplicationConsumerKeyRef: corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "applicationCustomerKey"}},
	}, false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Valid configuration for OAuth2
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:              "ovh-eu",
		AuthenticationMethod:  "oauth2",
		Oauth2ClientIDRef:     corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "oauth2ClientId"}},
		Oauth2ClientSecretRef: corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "oauth2ClientSecret"}},
	}, false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestFindMatchingZone(t *testing.T) {
	tests := []struct {
		name      string
		zones     []string
		fqdn      string
		expected  string
		expectErr bool
	}{
		{
			name:     "single zone match",
			zones:    []string{"example.com"},
			fqdn:     "_acme-challenge.example.com.",
			expected: "example.com",
		},
		{
			name:     "deepest zone wins",
			zones:    []string{"com", "example.com", "sub.example.com"},
			fqdn:     "_acme-challenge.sub.example.com.",
			expected: "sub.example.com",
		},
		{
			name:     "fqdn equals zone exactly",
			zones:    []string{"example.com"},
			fqdn:     "example.com.",
			expected: "example.com",
		},
		{
			name:     "fqdn without trailing dot",
			zones:    []string{"example.com"},
			fqdn:     "_acme-challenge.example.com",
			expected: "example.com",
		},
		{
			name:      "no matching zone",
			zones:     []string{"other.com", "another.org"},
			fqdn:      "_acme-challenge.example.com.",
			expectErr: true,
		},
		{
			name:      "empty zones list",
			zones:     []string{},
			fqdn:      "_acme-challenge.example.com.",
			expectErr: true,
		},
		{
			name:     "multiple zones only one matches",
			zones:    []string{"other.com", "example.com", "another.org"},
			fqdn:     "_acme-challenge.test.example.com.",
			expected: "example.com",
		},
		{
			name:     "zone with trailing dot",
			zones:    []string{"example.com."},
			fqdn:     "_acme-challenge.example.com.",
			expected: "example.com",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := findMatchingZone(tc.zones, tc.fqdn)
			if tc.expectErr {
				if err == nil {
					t.Errorf("expected error, got result %q", result)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestConfigDeserializationWithUseOvhApiZoneResolution(t *testing.T) {
	// Test that useOvhApiZoneResolution defaults to false
	jsonData := []byte(`{"endpoint": "ovh-eu", "authenticationMethod": "application"}`)
	var cfg ovhDNSProviderConfig
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if cfg.UseOvhApiZoneResolution {
		t.Error("expected UseOvhApiZoneResolution to default to false")
	}

	// Test that useOvhApiZoneResolution can be set to true
	jsonData = []byte(`{"endpoint": "ovh-eu", "authenticationMethod": "application", "useOvhApiZoneResolution": true}`)
	cfg = ovhDNSProviderConfig{}
	if err := json.Unmarshal(jsonData, &cfg); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if !cfg.UseOvhApiZoneResolution {
		t.Error("expected UseOvhApiZoneResolution to be true")
	}
}
