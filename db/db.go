package db

import "errors"

var ErrNoTicketsRemaining = errors.New("no tickets remaining")
var ErrNoReservation = errors.New("no reservation for given name")
var ErrAlreadyCharged = errors.New("reservation already charged")
var ErrReservationExpired = errors.New("reservation has expired")

// TicketsDB defines database layer operations for ticket management.
type TicketsDB interface {
	// Remaining returns number of tickets still
	// available for sale.
	// Note that by the time this method returns, the number
	// may already be out of date.
	Remaining() (int, error)
	// Reserve reserves ticket for specified customer name.
	// Only one ticket per name can be reserved, if specified
	// customer already has a ticket reserved (but not charged),
	// he gets his reservation time renewed.
	//
	// Returns ErrNoTicketsRemaining if there are no free tickets anymore
	// Note that this does not need to corelate with any earlier result
	// of `.Remaining` calls.
	Reserve(fullname string) error
	// Charge flags the customer reservation as charged.
	//
	// Returns ErrNoReservation if there is no reservation for given customer,
	// ErrAlreadyCharged if reservation is already charged.
	Charge(fullname string) error
	// Guests returns a list of *actual* guests, i.e. ones that are
	// already both reserved and charged.
	Guests() ([]string, error)

	// Reset resets the database to initial state.
	//
	// Mostly for testing/debugging purposes.
	Reset() error
}
