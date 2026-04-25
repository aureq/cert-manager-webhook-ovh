package main

import (
	"encoding/json"
	"strings"
	"testing"

	// acmetest "github.com/cert-manager/cert-manager/test/acme"
	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

// jsonRaw wraps a JSON string into an *extapi.JSON suitable for loadConfig().
func jsonRaw(s string) *extapi.JSON {
	return &extapi.JSON{Raw: []byte(s)}
}

// makeSecret creates a Kubernetes Secret object for use with the fake clientset.
func makeSecret(namespace, name string, data map[string][]byte) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

// solverWithSecrets creates an ovhDNSProviderSolver with a fake Kubernetes
// clientset pre-loaded with the given secrets.
func solverWithSecrets(secrets ...*corev1.Secret) *ovhDNSProviderSolver {
	objs := make([]runtime.Object, len(secrets))
	for i, s := range secrets {
		objs[i] = s
	}
	return &ovhDNSProviderSolver{
		client: fake.NewClientset(objs...),
	}
}

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

func TestName(t *testing.T) {
	s := &ovhDNSProviderSolver{}
	if name := s.Name(); name != "ovh" {
		t.Errorf("expected solver name 'ovh', got %q", name)
	}
}

// TestLoadConfig verifies JSON configuration parsing, covering nil input,
// valid configs, malformed JSON, type mismatches, unknown fields, and large
// payloads to ensure stability and security against malformed input.
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   *extapi.JSON
		wantErr bool
		errMsg  string
		check   func(t *testing.T, cfg ovhDNSProviderConfig)
	}{
		{
			// Stability: loadConfig must handle nil input without panicking.
			name:    "nil input returns empty config",
			input:   nil,
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.Endpoint != "" {
					t.Errorf("expected empty endpoint, got %q", cfg.Endpoint)
				}
			},
		},
		{
			// Stability: an empty JSON object should produce a zero-value config.
			name:    "empty JSON object",
			input:   jsonRaw(`{}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.Endpoint != "" || cfg.AuthenticationMethod != "" {
					t.Error("expected zero-value config for empty JSON object")
				}
			},
		},
		{
			// Stability: all fields should be correctly populated from valid JSON.
			name: "valid full config with all fields",
			input: jsonRaw(`{
				"endpoint": "ovh-eu",
				"authenticationMethod": "application",
				"applicationKeyRef": {"name": "secret1", "key": "appKey"},
				"applicationSecretRef": {"name": "secret1", "key": "appSecret"},
				"applicationConsumerKeyRef": {"name": "secret1", "key": "consumerKey"},
				"oauth2ClientIDRef": {"name": "secret2", "key": "clientId"},
				"oauth2ClientSecretRef": {"name": "secret2", "key": "clientSecret"}
			}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.Endpoint != "ovh-eu" {
					t.Errorf("expected endpoint 'ovh-eu', got %q", cfg.Endpoint)
				}
				if cfg.AuthenticationMethod != "application" {
					t.Errorf("expected authMethod 'application', got %q", cfg.AuthenticationMethod)
				}
				if cfg.ApplicationKeyRef.Name != "secret1" || cfg.ApplicationKeyRef.Key != "appKey" {
					t.Errorf("unexpected ApplicationKeyRef: %+v", cfg.ApplicationKeyRef)
				}
				if cfg.ApplicationSecretRef.Name != "secret1" || cfg.ApplicationSecretRef.Key != "appSecret" {
					t.Errorf("unexpected ApplicationSecretRef: %+v", cfg.ApplicationSecretRef)
				}
				if cfg.ApplicationConsumerKeyRef.Name != "secret1" || cfg.ApplicationConsumerKeyRef.Key != "consumerKey" {
					t.Errorf("unexpected ApplicationConsumerKeyRef: %+v", cfg.ApplicationConsumerKeyRef)
				}
				if cfg.Oauth2ClientIDRef.Name != "secret2" || cfg.Oauth2ClientIDRef.Key != "clientId" {
					t.Errorf("unexpected Oauth2ClientIDRef: %+v", cfg.Oauth2ClientIDRef)
				}
				if cfg.Oauth2ClientSecretRef.Name != "secret2" || cfg.Oauth2ClientSecretRef.Key != "clientSecret" {
					t.Errorf("unexpected Oauth2ClientSecretRef: %+v", cfg.Oauth2ClientSecretRef)
				}
			},
		},
		{
			// Stability: partial config should populate only provided fields.
			name:    "partial config with only endpoint",
			input:   jsonRaw(`{"endpoint": "ovh-ca"}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.Endpoint != "ovh-ca" {
					t.Errorf("expected endpoint 'ovh-ca', got %q", cfg.Endpoint)
				}
				if cfg.AuthenticationMethod != "" {
					t.Errorf("expected empty authMethod, got %q", cfg.AuthenticationMethod)
				}
			},
		},
		{
			// Security: syntactically invalid JSON must be rejected to prevent
			// misconfiguration from silently passing through.
			name:    "invalid JSON syntax",
			input:   jsonRaw(`{invalid`),
			wantErr: true,
			errMsg:  "error decoding OVH config",
		},
		{
			// Security: unknown fields in user-supplied config must be silently
			// ignored and not populate unexpected struct fields.
			name:    "unknown fields are ignored",
			input:   jsonRaw(`{"endpoint": "ovh-eu", "maliciousField": "evil"}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.Endpoint != "ovh-eu" {
					t.Errorf("expected endpoint 'ovh-eu', got %q", cfg.Endpoint)
				}
			},
		},
		{
			// Security: a type mismatch (number instead of string) must produce
			// an error rather than silently coercing the value.
			name:    "wrong type for string field",
			input:   jsonRaw(`{"endpoint": 12345}`),
			wantErr: true,
			errMsg:  "error decoding OVH config",
		},
		{
			// Stability: nested SecretKeySelector structs must deserialize correctly.
			name:    "nested SecretKeySelector",
			input:   jsonRaw(`{"applicationKeyRef": {"name": "my-secret", "key": "my-key"}}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if cfg.ApplicationKeyRef.Name != "my-secret" {
					t.Errorf("expected ref name 'my-secret', got %q", cfg.ApplicationKeyRef.Name)
				}
				if cfg.ApplicationKeyRef.Key != "my-key" {
					t.Errorf("expected ref key 'my-key', got %q", cfg.ApplicationKeyRef.Key)
				}
			},
		},
		{
			// Stability: empty raw bytes are not valid JSON and must error.
			name:    "empty raw bytes",
			input:   &extapi.JSON{Raw: []byte{}},
			wantErr: true,
			errMsg:  "error decoding OVH config",
		},
		{
			// Security: a very large JSON payload must not cause a panic (DoS resilience).
			name:    "large payload does not panic",
			input:   jsonRaw(`{"endpoint": "` + strings.Repeat("a", 1024*1024) + `"}`),
			wantErr: false,
			check: func(t *testing.T, cfg ovhDNSProviderConfig) {
				if len(cfg.Endpoint) != 1024*1024 {
					t.Errorf("expected 1MB endpoint string, got length %d", len(cfg.Endpoint))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := loadConfig(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.check != nil {
				tt.check(t, cfg)
			}
		})
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

// TestGetSubDomain verifies subdomain extraction from FQDNs, including
// trailing dots, empty inputs, partial suffix matches, and deeply nested
// subdomains. Incorrect extraction could cause DNS records to be created
// or deleted on the wrong subdomain, which is both a stability and security concern.
func TestGetSubDomain(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		fqdn     string
		expected string
	}{
		// Stability: standard ACME challenge subdomain extraction.
		{"standard ACME challenge", "example.com", "_acme-challenge.example.com", "_acme-challenge"},
		// Stability: multi-level subdomain must preserve all levels.
		{"nested subdomain", "example.com", "_acme-challenge.sub.example.com", "_acme-challenge.sub"},
		// Stability: cert-manager may send FQDNs with trailing dots.
		{"FQDN with trailing dot", "example.com", "_acme-challenge.example.com.", "_acme-challenge"},
		// Stability: domain with trailing dot must be normalized.
		{"domain with trailing dot", "example.com.", "_acme-challenge.example.com", "_acme-challenge"},
		// Stability: both domain and FQDN with trailing dots.
		{"both with trailing dots", "example.com.", "_acme-challenge.example.com.", "_acme-challenge"},
		// Stability: apex domain (FQDN equals domain) returns empty string.
		{"FQDN equals domain", "example.com", "example.com", ""},
		// Stability: both with dots and equal.
		{"FQDN equals domain both dotted", "example.com.", "example.com.", ""},
		// Security: when domain is not found in FQDN, returns full FQDN without panic.
		{"domain not in FQDN", "other.com", "_acme-challenge.example.com", "_acme-challenge.example.com"},
		// Security: empty domain must not panic.
		{"empty domain", "", "_acme-challenge.example.com", "_acme-challenge.example.com"},
		// Security: empty FQDN must not panic.
		{"empty fqdn", "example.com", "", ""},
		// Security: both empty must not panic.
		{"both empty", "", "", ""},
		// Security: partial suffix match must NOT match. "le.com" is a suffix of
		// "example.com" but not at a dot boundary. This prevents acting on wrong domains.
		{"partial suffix match at non-boundary", "le.com", "_acme-challenge.example.com", "_acme-challenge.example.com"},
		// Stability: deeply nested subdomains must be fully preserved.
		{"deeply nested subdomain", "example.com", "a.b.c.d.e.example.com", "a.b.c.d.e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSubDomain(tt.domain, tt.fqdn)
			if result != tt.expected {
				t.Errorf("getSubDomain(%q, %q) = %q, want %q", tt.domain, tt.fqdn, result, tt.expected)
			}
		})
	}
}

// TestValidateAmbientCredentials verifies that when allowAmbientCredentials
// is true, all config validation is bypassed. This is security-relevant because
// it documents an intentional bypass -- OVH client will load credentials from
// environment variables and ovh.conf files instead of Kubernetes secrets.
func TestValidateAmbientCredentials(t *testing.T) {
	s := &ovhDNSProviderSolver{}

	// Security: completely empty config with ambient=true must pass validation,
	// since credentials will come from the environment.
	err := s.validate(&ovhDNSProviderConfig{}, true)
	if err != nil {
		t.Errorf("expected no error with ambient credentials and empty config, got %v", err)
	}

	// Security: even an invalid auth method with ambient=true must pass,
	// confirming ambient mode skips all checks.
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "invalid-method",
	}, true)
	if err != nil {
		t.Errorf("expected no error with ambient credentials and invalid auth method, got %v", err)
	}
}

// TestValidateEdgeCases covers authentication method edge cases to ensure
// no bypass via case variation, whitespace, or missing values.
func TestValidateEdgeCases(t *testing.T) {
	s := &ovhDNSProviderSolver{}

	// Security: empty authentication method must be rejected (falls into default case).
	err := s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "",
	}, false)
	if err == nil {
		t.Error("expected error for empty authentication method, got nil")
	}

	// Security: uppercase "Application" must be rejected. The auth method check
	// is case-sensitive, preventing bypass via case variation.
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "Application",
	}, false)
	if err == nil {
		t.Error("expected error for uppercase 'Application', got nil")
	}

	// Security: whitespace-padded auth method must be rejected. No trimming
	// means " application " is not the same as "application".
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: " application ",
	}, false)
	if err == nil {
		t.Error("expected error for whitespace-padded auth method, got nil")
	}

	// Security: ApplicationKeyRef with Key set but Name empty must be rejected.
	// Having only a Key without a secret Name is insufficient to retrieve credentials.
	err = s.validate(&ovhDNSProviderConfig{
		Endpoint:             "ovh-eu",
		AuthenticationMethod: "application",
		ApplicationKeyRef:    corev1.SecretKeySelector{Key: "some-key"},
	}, false)
	if err == nil {
		t.Error("expected error when ApplicationKeyRef has Key but no Name")
	}
}

// TestSecret verifies Kubernetes secret retrieval, including value sanitization
// (trailing newline trimming), error handling for missing secrets/keys,
// and namespace isolation to prevent cross-namespace secret access.
func TestSecret(t *testing.T) {
	tests := []struct {
		name      string
		secrets   []*corev1.Secret
		ref       corev1.SecretKeySelector
		namespace string
		wantVal   string
		wantErr   bool
		errMsg    string
	}{
		{
			// Stability: basic secret retrieval must work correctly.
			name: "happy path retrieves correct value",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc123")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "abc123",
		},
		{
			// Security: credentials often have trailing newlines from file-based
			// secret creation. TrimSuffix ensures they are removed.
			name: "trailing newline is trimmed",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc123\n")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "abc123",
		},
		{
			// Stability: values without trailing newlines must not be altered.
			name: "no trailing newline unchanged",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc123")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "abc123",
		},
		{
			// Security: a secret containing only a newline should produce an empty
			// string after trimming, not a newline character.
			name: "value with only newline",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("\n")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "",
		},
		{
			// Security: TrimSuffix only removes one trailing "\n". Multiple trailing
			// newlines should leave all but the last.
			name: "multiple trailing newlines trims only last",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc\n\n")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "abc\n",
		},
		{
			// Security: embedded newlines within a value must be preserved.
			name: "embedded newline preserved",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc\ndef")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantVal:   "abc\ndef",
		},
		{
			// Security: empty ref.Name must be rejected early without making any
			// Kubernetes API call.
			name:      "empty ref name returns error",
			secrets:   nil,
			ref:       corev1.SecretKeySelector{Key: "token"},
			namespace: "default",
			wantErr:   true,
			errMsg:    "no secret name specified",
		},
		{
			// Stability: requesting a non-existent secret must return an error.
			name:      "secret not found",
			secrets:   nil,
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "nonexistent"}, Key: "token"},
			namespace: "default",
			wantErr:   true,
		},
		{
			// Stability: requesting a key that does not exist in the secret data
			// must return a descriptive error.
			name: "key not found in secret",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{"token": []byte("abc")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "wrong-key"},
			namespace: "default",
			wantErr:   true,
			errMsg:    "key not found",
		},
		{
			// Security: a secret in namespace "other" must not be accessible from
			// namespace "default". This verifies namespace isolation.
			name: "wrong namespace enforces isolation",
			secrets: []*corev1.Secret{
				makeSecret("other", "my-secret", map[string][]byte{"token": []byte("abc")}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantErr:   true,
		},
		{
			// Stability: a secret with an empty Data map must return a key-not-found error.
			name: "empty data map",
			secrets: []*corev1.Secret{
				makeSecret("default", "my-secret", map[string][]byte{}),
			},
			ref:       corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "my-secret"}, Key: "token"},
			namespace: "default",
			wantErr:   true,
			errMsg:    "key not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := solverWithSecrets(tt.secrets...)
			val, err := s.secret(tt.ref, tt.namespace)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if val != tt.wantVal {
				t.Errorf("expected value %q, got %q", tt.wantVal, val)
			}
		})
	}
}

// TestConfigJSONFieldNames verifies that JSON serialization field names match
// the expected camelCase names, and that a roundtrip marshal/unmarshal produces
// an identical config. This prevents silent misconfiguration if JSON tags are
// accidentally changed.
func TestConfigJSONFieldNames(t *testing.T) {
	// Stability: roundtrip marshal/unmarshal must produce identical config.
	t.Run("roundtrip marshal unmarshal", func(t *testing.T) {
		original := ovhDNSProviderConfig{
			Endpoint:                  "ovh-eu",
			AuthenticationMethod:      "application",
			ApplicationKeyRef:         corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s1"}, Key: "k1"},
			ApplicationSecretRef:      corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s2"}, Key: "k2"},
			ApplicationConsumerKeyRef: corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s3"}, Key: "k3"},
			Oauth2ClientIDRef:         corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s4"}, Key: "k4"},
			Oauth2ClientSecretRef:     corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s5"}, Key: "k5"},
		}

		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("failed to marshal config: %v", err)
		}

		var restored ovhDNSProviderConfig
		if err := json.Unmarshal(data, &restored); err != nil {
			t.Fatalf("failed to unmarshal config: %v", err)
		}

		if original != restored {
			t.Errorf("roundtrip mismatch:\n  original: %+v\n  restored: %+v", original, restored)
		}
	})

	// Security: verify expected JSON field names appear in marshaled output.
	// Incorrect tags would cause user-supplied config to silently fail to populate.
	t.Run("expected JSON field names present", func(t *testing.T) {
		cfg := ovhDNSProviderConfig{
			Endpoint:             "ovh-eu",
			AuthenticationMethod: "application",
			ApplicationKeyRef:    corev1.SecretKeySelector{LocalObjectReference: corev1.LocalObjectReference{Name: "s1"}},
		}

		data, err := json.Marshal(cfg)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}
		jsonStr := string(data)

		expectedFields := []string{
			`"endpoint"`,
			`"authenticationMethod"`,
			`"applicationKeyRef"`,
			`"applicationSecretRef"`,
			`"applicationConsumerKeyRef"`,
			`"oauth2ClientIDRef"`,
			`"oauth2ClientSecretRef"`,
		}
		for _, field := range expectedFields {
			if !strings.Contains(jsonStr, field) {
				t.Errorf("expected JSON field %s not found in output: %s", field, jsonStr)
			}
		}
	})

	// Security: snake_case field names must NOT populate the struct. This ensures
	// only the correct camelCase JSON tags are accepted.
	t.Run("snake_case field names rejected", func(t *testing.T) {
		input := `{"authentication_method": "application", "endpoint": "ovh-eu"}`
		var cfg ovhDNSProviderConfig
		if err := json.Unmarshal([]byte(input), &cfg); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if cfg.AuthenticationMethod != "" {
			t.Errorf("snake_case 'authentication_method' should not populate AuthenticationMethod, got %q", cfg.AuthenticationMethod)
		}
		// The camelCase endpoint should still work.
		if cfg.Endpoint != "ovh-eu" {
			t.Errorf("expected endpoint 'ovh-eu', got %q", cfg.Endpoint)
		}
	})
}
