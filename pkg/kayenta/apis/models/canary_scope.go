// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// CanaryScope CanaryScope
//
// swagger:model CanaryScope
type CanaryScope struct {

	// end
	// Format: date-time
	End strfmt.DateTime `json:"end,omitempty"`

	// extended scope params
	ExtendedScopeParams map[string]string `json:"extendedScopeParams,omitempty"`

	// location
	Location string `json:"location,omitempty"`

	// scope
	Scope string `json:"scope,omitempty"`

	// start
	// Format: date-time
	Start strfmt.DateTime `json:"start,omitempty"`

	// step
	Step int64 `json:"step,omitempty"`
}

// Validate validates this canary scope
func (m *CanaryScope) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEnd(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStart(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CanaryScope) validateEnd(formats strfmt.Registry) error {
	if swag.IsZero(m.End) { // not required
		return nil
	}

	if err := validate.FormatOf("end", "body", "date-time", m.End.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *CanaryScope) validateStart(formats strfmt.Registry) error {
	if swag.IsZero(m.Start) { // not required
		return nil
	}

	if err := validate.FormatOf("start", "body", "date-time", m.Start.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this canary scope based on context it is used
func (m *CanaryScope) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *CanaryScope) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CanaryScope) UnmarshalBinary(b []byte) error {
	var res CanaryScope
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
