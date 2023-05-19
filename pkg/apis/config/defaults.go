package config

import (
	corev1 "k8s.io/api/core/v1"
)

const (

	// DefaultsConfigName is the name of config map for the defaults.
	DefaultsConfigName = "config-defaults"
)

// +k8s:deepcopy-gen=true

// Defaults holds default configurations of the system.
type Defaults struct {
}

var defaultsParsers = map[string]Parser{}

// NewDefaultsFromConfigMap takes a ConfigMap and returns a Defaults object.
func NewDefaultsFromConfigMap(configMap *corev1.ConfigMap) (*Defaults, error) {
	defaults := &Defaults{}
	if err := Unmarshal(defaults, configMap, defaultsParsers); err != nil {
		return nil, err
	}

	return defaults, nil
}
