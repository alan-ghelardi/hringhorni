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

// HydrateQueueOutput HydrateQueueOutput
//
// swagger:model HydrateQueueOutput
type HydrateQueueOutput struct {

	// dry run
	// Required: true
	DryRun *bool `json:"dryRun"`

	// executions
	// Required: true
	Executions map[string]ProcessedExecution `json:"executions"`
}

// Validate validates this hydrate queue output
func (m *HydrateQueueOutput) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDryRun(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateExecutions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HydrateQueueOutput) validateDryRun(formats strfmt.Registry) error {

	if err := validate.Required("dryRun", "body", m.DryRun); err != nil {
		return err
	}

	return nil
}

func (m *HydrateQueueOutput) validateExecutions(formats strfmt.Registry) error {

	if err := validate.Required("executions", "body", m.Executions); err != nil {
		return err
	}

	for k := range m.Executions {

		if err := validate.Required("executions"+"."+k, "body", m.Executions[k]); err != nil {
			return err
		}
		if val, ok := m.Executions[k]; ok {
			if err := val.Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("executions" + "." + k)
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("executions" + "." + k)
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this hydrate queue output based on the context it is used
func (m *HydrateQueueOutput) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateExecutions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HydrateQueueOutput) contextValidateExecutions(ctx context.Context, formats strfmt.Registry) error {

	if err := validate.Required("executions", "body", m.Executions); err != nil {
		return err
	}

	for k := range m.Executions {

		if val, ok := m.Executions[k]; ok {
			if err := val.ContextValidate(ctx, formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *HydrateQueueOutput) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HydrateQueueOutput) UnmarshalBinary(b []byte) error {
	var res HydrateQueueOutput
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
