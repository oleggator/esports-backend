// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// Tournament tournament
// swagger:model Tournament

type Tournament struct {

	// Идентификатор турнира
	ID int64 `json:"id,omitempty"`

	// Slug
	Slug string `json:"slug,omitempty"`

	// Название туринира
	// Required: true
	Title *string `json:"title"`
}

/* polymorph Tournament id false */

/* polymorph Tournament slug false */

/* polymorph Tournament title false */

// Validate validates this tournament
func (m *Tournament) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTitle(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Tournament) validateTitle(formats strfmt.Registry) error {

	if err := validate.Required("title", "body", m.Title); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Tournament) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Tournament) UnmarshalBinary(b []byte) error {
	var res Tournament
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
