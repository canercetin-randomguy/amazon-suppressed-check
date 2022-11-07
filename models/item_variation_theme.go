// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ItemVariationTheme Variation theme indicating the combination of Amazon item catalog attributes that define the variation family.
//
// swagger:model ItemVariationTheme
type ItemVariationTheme struct {

	// Names of the Amazon catalog item attributes associated with the variation theme.
	Attributes []string `json:"attributes"`

	// Variation theme indicating the combination of Amazon item catalog attributes that define the variation family.
	// Example: COLOR_NAME/STYLE_NAME
	Theme string `json:"theme,omitempty"`
}

// Validate validates this item variation theme
func (m *ItemVariationTheme) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this item variation theme based on context it is used
func (m *ItemVariationTheme) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ItemVariationTheme) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ItemVariationTheme) UnmarshalBinary(b []byte) error {
	var res ItemVariationTheme
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
