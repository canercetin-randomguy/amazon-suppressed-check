// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Refinements Search refinements.
//
// swagger:model Refinements
type Refinements struct {

	// Brand search refinements.
	// Required: true
	Brands []*BrandRefinement `json:"brands"`

	// Classification search refinements.
	// Required: true
	Classifications []*ClassificationRefinement `json:"classifications"`
}

// Validate validates this refinements
func (m *Refinements) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBrands(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateClassifications(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Refinements) validateBrands(formats strfmt.Registry) error {

	if err := validate.Required("brands", "body", m.Brands); err != nil {
		return err
	}

	for i := 0; i < len(m.Brands); i++ {
		if swag.IsZero(m.Brands[i]) { // not required
			continue
		}

		if m.Brands[i] != nil {
			if err := m.Brands[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("brands" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("brands" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Refinements) validateClassifications(formats strfmt.Registry) error {

	if err := validate.Required("classifications", "body", m.Classifications); err != nil {
		return err
	}

	for i := 0; i < len(m.Classifications); i++ {
		if swag.IsZero(m.Classifications[i]) { // not required
			continue
		}

		if m.Classifications[i] != nil {
			if err := m.Classifications[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("classifications" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("classifications" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this refinements based on the context it is used
func (m *Refinements) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBrands(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateClassifications(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Refinements) contextValidateBrands(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Brands); i++ {

		if m.Brands[i] != nil {
			if err := m.Brands[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("brands" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("brands" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *Refinements) contextValidateClassifications(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Classifications); i++ {

		if m.Classifications[i] != nil {
			if err := m.Classifications[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("classifications" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("classifications" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *Refinements) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Refinements) UnmarshalBinary(b []byte) error {
	var res Refinements
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}