package domain

import (
	"eda-in-golang/internal/ddd"
	"errors"
)

var (
	ErrStoreNameIsBlank               = errors.New("store name cannot be blank")
	ErrStoreLocationIsBlank           = errors.New("store location cannot be blank")
	ErrStoreIsAlreadyParticipating    = errors.New("store is already participating")
	ErrStoreIsAlreadyNotParticipating = errors.New("store is not participating")
)

type Store struct {
	ddd.AggregateBase
	Name          string
	Location      string
	Participating bool
}

func CreateStore(id, name, location string) (store *Store, err error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}
	if location == "" {
		return nil, ErrStoreLocationIsBlank
	}

	store = &Store{
		AggregateBase: ddd.AggregateBase{
			ID: id,
		},
		Name:          name,
		Location:      location,
		Participating: false,
	}

	store.AddEvent(&StoreCreated{
		Store: store,
	})

	return
}

func (s *Store) EnableParticipation() (err error) {
	if s.Participating {
		return ErrStoreIsAlreadyParticipating
	}

	s.Participating = true

	s.AddEvent(&StoreParticipationEnabled{
		Store: s,
	})

	return nil
}

func (s *Store) DisableParticipation() (err error) {
	if !s.Participating {
		return ErrStoreIsAlreadyNotParticipating
	}

	s.Participating = false

	s.AddEvent(&StoreParticipationDisabled{
		Store: s,
	})

	return nil
}
