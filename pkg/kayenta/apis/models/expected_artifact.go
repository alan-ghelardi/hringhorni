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

// ExpectedArtifact ExpectedArtifact
//
// swagger:model ExpectedArtifact
type ExpectedArtifact struct {

	// bound artifact
	BoundArtifact *Artifact `json:"boundArtifact,omitempty"`

	// default artifact
	DefaultArtifact *Artifact `json:"defaultArtifact,omitempty"`

	// id
	ID string `json:"id,omitempty"`

	// match artifact
	MatchArtifact *Artifact `json:"matchArtifact,omitempty"`

	// use default artifact
	UseDefaultArtifact bool `json:"useDefaultArtifact,omitempty"`

	// use prior artifact
	UsePriorArtifact bool `json:"usePriorArtifact,omitempty"`
}

// Validate validates this expected artifact
func (m *ExpectedArtifact) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBoundArtifact(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDefaultArtifact(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateMatchArtifact(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExpectedArtifact) validateBoundArtifact(formats strfmt.Registry) error {
	if swag.IsZero(m.BoundArtifact) { // not required
		return nil
	}

	if m.BoundArtifact != nil {
		if err := m.BoundArtifact.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("boundArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("boundArtifact")
			}
			return err
		}
	}

	return nil
}

func (m *ExpectedArtifact) validateDefaultArtifact(formats strfmt.Registry) error {
	if swag.IsZero(m.DefaultArtifact) { // not required
		return nil
	}

	if m.DefaultArtifact != nil {
		if err := m.DefaultArtifact.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("defaultArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("defaultArtifact")
			}
			return err
		}
	}

	return nil
}

func (m *ExpectedArtifact) validateMatchArtifact(formats strfmt.Registry) error {
	if swag.IsZero(m.MatchArtifact) { // not required
		return nil
	}

	if m.MatchArtifact != nil {
		if err := m.MatchArtifact.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("matchArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("matchArtifact")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this expected artifact based on the context it is used
func (m *ExpectedArtifact) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBoundArtifact(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateDefaultArtifact(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateMatchArtifact(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ExpectedArtifact) contextValidateBoundArtifact(ctx context.Context, formats strfmt.Registry) error {

	if m.BoundArtifact != nil {
		if err := m.BoundArtifact.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("boundArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("boundArtifact")
			}
			return err
		}
	}

	return nil
}

func (m *ExpectedArtifact) contextValidateDefaultArtifact(ctx context.Context, formats strfmt.Registry) error {

	if m.DefaultArtifact != nil {
		if err := m.DefaultArtifact.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("defaultArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("defaultArtifact")
			}
			return err
		}
	}

	return nil
}

func (m *ExpectedArtifact) contextValidateMatchArtifact(ctx context.Context, formats strfmt.Registry) error {

	if m.MatchArtifact != nil {
		if err := m.MatchArtifact.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("matchArtifact")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("matchArtifact")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ExpectedArtifact) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ExpectedArtifact) UnmarshalBinary(b []byte) error {
	var res ExpectedArtifact
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
