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

// ItemRelationshipsByMarketplace Relationship details for the Amazon catalog item for the indicated Amazon marketplace.
//
// swagger:model ItemRelationshipsByMarketplace
type ItemRelationshipsByMarketplace struct {

	// Amazon marketplace identifier.
	// Required: true
	MarketplaceID *string `json:"marketplaceId"`

	// Relationships for the item.
	// Required: true
	Relationships []*ItemRelationship `json:"relationships"`
}

// Validate validates this item relationships by marketplace
func (m *ItemRelationshipsByMarketplace) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateMarketplaceID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRelationships(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ItemRelationshipsByMarketplace) validateMarketplaceID(formats strfmt.Registry) error {

	if err := validate.Required("marketplaceId", "body", m.MarketplaceID); err != nil {
		return err
	}

	return nil
}

func (m *ItemRelationshipsByMarketplace) validateRelationships(formats strfmt.Registry) error {

	if err := validate.Required("relationships", "body", m.Relationships); err != nil {
		return err
	}

	for i := 0; i < len(m.Relationships); i++ {
		if swag.IsZero(m.Relationships[i]) { // not required
			continue
		}

		if m.Relationships[i] != nil {
			if err := m.Relationships[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("relationships" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("relationships" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this item relationships by marketplace based on the context it is used
func (m *ItemRelationshipsByMarketplace) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateRelationships(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ItemRelationshipsByMarketplace) contextValidateRelationships(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Relationships); i++ {

		if m.Relationships[i] != nil {
			if err := m.Relationships[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("relationships" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("relationships" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ItemRelationshipsByMarketplace) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ItemRelationshipsByMarketplace) UnmarshalBinary(b []byte) error {
	var res ItemRelationshipsByMarketplace
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
