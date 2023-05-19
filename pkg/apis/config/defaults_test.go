package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
)

func TestNewDefaults(t *testing.T) {
	tests := []struct {
		name string
		in   *corev1.ConfigMap
		want *Defaults
		err  error
	}{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewDefaultsFromConfigMap(test.in)
			if err != nil {
				if diff := cmp.Diff(test.err.Error(), err.Error()); diff != "" {
					t.Errorf("NewDefaultsFromConfigMap() returned an unexpected error: (-want +got):\n%s", diff)
				}
			} else if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NewDefaultsFromConfigMap: (-want +got):\n%s", diff)
			}
		})
	}
}
