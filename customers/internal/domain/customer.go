package domain

import (
	"eda-in-golang/internal/ddd"

	"github.com/stackus/errors"
)

type Customer struct {
	ddd.AggregateBase
	Name      string
	SmsNumber string
	Enabled   bool
}

var (
	ErrNameCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the customer name cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer ID cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the SMS number cannot be blank")
	ErrCustomerAlreadyEnabled  = errors.Wrap(errors.ErrBadRequest, "the customer is already enabled")
	ErrCustomerAlreadyDisabled = errors.Wrap(errors.ErrBadRequest, "the customer is already disabled")
	ErrCustomerNotAuthorized   = errors.Wrap(errors.ErrUnauthorized, "the customer is not authorized")
)

func NewCustomer(id, name, smsNumber string) (*Customer, error) {
	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if smsNumber == "" {
		return nil, ErrSmsNumberCannotBeBlank
	}

	customer := &Customer{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
		Name:      name,
		SmsNumber: smsNumber,
		Enabled:   true,
	}

	customer.AddEvent(&CustomerRegistered{Customer: customer})

	return customer, nil
}

func (c *Customer) Authorize() error {
	if !c.Enabled {
		return ErrCustomerNotAuthorized
	}

	c.AddEvent(&CustomerAuthorized{Customer: c})

	return nil
}

func (c *Customer) Enable() error {
	if c.Enabled {
		return ErrCustomerAlreadyEnabled
	}

	c.Enabled = true
	c.AddEvent(&CustomerEnabled{Customer: c})

	return nil
}

func (c *Customer) Disable() error {
	if !c.Enabled {
		return ErrCustomerAlreadyDisabled
	}

	c.Enabled = false
	c.AddEvent(&CustomerDisabled{Customer: c})

	return nil
}
