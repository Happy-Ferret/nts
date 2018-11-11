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
	customers map[string]*customerState
}

func NewMemTicketsDB(number int) *MemTicketsDB {
	return &MemTicketsDB{
		remaining: number,
		customers: map[string]*customerState{},
	}
}

func (m *MemTicketsDB) Remaining() (int, error) {
	// Cleanup expired reservations
	now := time.Now().UTC()
	for fullname, cs := range m.customers {
		if !cs.charged && now.Sub(cs.reserved).Minutes() > 5 {
			delete(m.customers, fullname)
			m.remaining++
		}
	}
	return m.remaining, nil
}

func (m *MemTicketsDB) Reserve(fullname string) error {
	customer, exists := m.customers[fullname]
	if !exists {
		if m.remaining <= 0 {
			return ErrNoTicketsRemaining
		}
		m.remaining-- // TODO: Make this atomic
		m.customers[fullname] = &customerState{
			reserved: time.Now().UTC(),
			charged:  false,
		}
		return nil
	}
	customer.reserved = time.Now().UTC()
	return nil
}

func (m *MemTicketsDB) Charge(fullname string) error {
	customer, exists := m.customers[fullname]
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
	for fullname, cs := range m.customers {
		if cs.charged {
			names = append(names, fullname)
		}
	}
	return names, nil
}

func (m *MemTicketsDB) Reset() error {
	m.remaining += len(m.customers)
	m.customers = map[string]*customerState{}
	return nil
}
