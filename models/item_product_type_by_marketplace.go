// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ItemProductTypeByMarketplace Product type associated with the Amazon catalog item for the indicated Amazon marketplace.
//
// swagger:model ItemProductTypeByMarketplace
type ItemProductTypeByMarketplace struct {

	// Amazon marketplace identifier.
	MarketplaceID string `json:"marketplaceId,omitempty"`

	// Name of the product type associated with the Amazon catalog item.
	// Example: LUGGAGE
	ProductType string `json:"productType,omitempty"`
}

// Validate validates this item product type by marketplace
func (m *ItemProductTypeByMarketplace) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this item product type by marketplace based on context it is used
func (m *ItemProductTypeByMarketplace) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ItemProductTypeByMarketplace) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ItemProductTypeByMarketplace) UnmarshalBinary(b []byte) error {
	var res ItemProductTypeByMarketplace
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}