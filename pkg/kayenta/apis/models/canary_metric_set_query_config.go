// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CanaryMetricSetQueryConfig CanaryMetricSetQueryConfig
//
// swagger:model CanaryMetricSetQueryConfig
type CanaryMetricSetQueryConfig struct {

	// custom filter
	CustomFilter string `json:"customFilter,omitempty"`

	// custom filter template
	CustomFilterTemplate string `json:"customFilterTemplate,omitempty"`

	// custom inline template
	CustomInlineTemplate string `json:"customInlineTemplate,omitempty"`

	// service type
	ServiceType string `json:"serviceType,omitempty"`

	// type
	Type string `json:"type,omitempty"`
}

// Validate validates this canary metric set query config
func (m *CanaryMetricSetQueryConfig) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this canary metric set query config based on context it is used
func (m *CanaryMetricSetQueryConfig) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CanaryMetricSetQueryConfig) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CanaryMetricSetQueryConfig) UnmarshalBinary(b []byte) error {
	var res CanaryMetricSetQueryConfig
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
