package db

import (
	"time"
)

type customerState struct {
	reserved time.Time
	charged  bool
}

type MemTicketsDB struct {
	remaining int
	// This field is exposed for the sake of testing
	// In real scenario, this memdb whole will most
	// probably sit in the testing module and a "real"
	// DB implementation will sit here instead.
	Customers map[string]*customerState
}

func NewMemTicketsDB(number int) *MemTicketsDB {
	return &MemTicketsDB{
		remaining: number,
		Customers: map[string]*customerState{},
	}
}

func (m *MemTicketsDB) Remaining() (int, error) {
	now := time.Now().UTC()
	for _, cs := range m.Customers {
		if !cs.charged && now.Sub(cs.reserved).Minutes() > 5 {
			m.remaining++
		}
	}
	return m.remaining, nil
}

func (m *MemTicketsDB) Reserve(fullname string) error {
	customer, exists := m.Customers[fullname]
	if !exists {
		if m.remaining <= 0 {
			return ErrNoTicketsRemaining
		}
		m.remaining-- // TODO: Make this atomic
		m.Customers[fullname] = &customerState{
			reserved: time.Now().UTC(),
			charged:  false,
		}
		return nil
	}
	customer.reserved = time.Now().UTC()
	return nil
}

func (m *MemTicketsDB) Charge(fullname string) error {
	customer, exists := m.Customers[fullname]
	if !exists {
		return ErrNoReservation
	}
	if customer.charged {
		return ErrAlreadyCharged
	}
	customer.charged = true
	return nil
}

func (m *MemTicketsDB) Guests() ([]string, error) {
	names := []string{}
	for fullname, cs := range m.Customers {
		if cs.charged {
			names = append(names, fullname)
		}
	}
	return names, nil
}

func (m *MemTicketsDB) Reset() error {
	m.remaining += len(m.Customers)
	m.Customers = map[string]*customerState{}
	return nil
}
