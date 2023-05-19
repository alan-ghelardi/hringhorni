// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CanaryConfigUpdateResponse CanaryConfigUpdateResponse
//
// swagger:model CanaryConfigUpdateResponse
type CanaryConfigUpdateResponse struct {

	// canary config Id
	CanaryConfigID string `json:"canaryConfigId,omitempty"`
}

// Validate validates this canary config update response
func (m *CanaryConfigUpdateResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this canary config update response based on context it is used
func (m *CanaryConfigUpdateResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CanaryConfigUpdateResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CanaryConfigUpdateResponse) UnmarshalBinary(b []byte) error {
	var res CanaryConfigUpdateResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
