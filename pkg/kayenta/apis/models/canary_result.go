// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// CanaryResult CanaryResult
//
// swagger:model CanaryResult
type CanaryResult struct {

	// canary duration
	CanaryDuration string `json:"canaryDuration,omitempty"`

	// judge result
	JudgeResult *CanaryJudgeResult `json:"judgeResult,omitempty"`
}

// Validate validates this canary result
func (m *CanaryResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateJudgeResult(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CanaryResult) validateJudgeResult(formats strfmt.Registry) error {
	if swag.IsZero(m.JudgeResult) { // not required
		return nil
	}

	if m.JudgeResult != nil {
		if err := m.JudgeResult.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("judgeResult")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("judgeResult")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this canary result based on the context it is used
func (m *CanaryResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateJudgeResult(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *CanaryResult) contextValidateJudgeResult(ctx context.Context, formats strfmt.Registry) error {

	if m.JudgeResult != nil {
		if err := m.JudgeResult.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("judgeResult")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("judgeResult")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *CanaryResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *CanaryResult) UnmarshalBinary(b []byte) error {
	var res CanaryResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
