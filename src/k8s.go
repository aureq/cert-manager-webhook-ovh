package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *ovhDNSProviderSolver) secret(ref corev1.SecretKeySelector, namespace string) (string, error) {
	// log.V(logf.TraceLevel).Info("Entering secret() function")
	var err error
	if ref.Name == "" {
		err = errors.New("no secret name specified in secret key selector")
		// log.V(logf.ErrorLevel).Error(err, "No secret name specified in secret key selector")
		return "", err
	}

	secret, err := s.client.CoreV1().Secrets(namespace).Get(context.TODO(), ref.Name, metav1.GetOptions{})
	if err != nil {
		// log.V(logf.ErrorLevel).Error(err, "Error retrieving secret", "namespace", namespace, "secret", ref.Name)
		return "", err
	}

	bytes, ok := secret.Data[ref.Key]
	if !ok {
		err = fmt.Errorf("key not found %q in secret '%s/%s'", ref.Key, namespace, ref.Name)
		// log.V(logf.ErrorLevel).Error(err, "Key not found in secret", "namespace", namespace, "secret", ref.Name, "key", ref.Key)
		return "", err
	}
	// log.V(logf.TraceLevel).Info("Exiting secret() function")
	return strings.TrimSuffix(string(bytes), "\n"), nil
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
		// log.V(logf.ErrorLevel).Error(err, "Failed to decode OVH config")
		return cfg, fmt.Errorf("error decoding OVH config: %v", err)
	}

	return cfg, nil
}
