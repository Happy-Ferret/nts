package db

import (
	"sync"
	"time"

	"go.uber.org/atomic"
)

type customerState struct {
	reserved time.Time
	charged  *atomic.Bool
}

type MemTicketsDB struct {
	remaining int
	// This field is exposed for the sake of testing
	// In real scenario, this memdb whole will most
	// probably sit in the testing module and a "real"
	// DB implementation will sit here instead.
	Customers map[string]*customerState
	lock      sync.RWMutex
}

func NewMemTicketsDB(number int) *MemTicketsDB {
	return &MemTicketsDB{
		remaining: number,
		Customers: map[string]*customerState{},
	}
}

func (m *MemTicketsDB) Remaining() (int, error) {
	now := time.Now().UTC()
	m.lock.Lock()
	for _, cs := range m.Customers {
		if !cs.charged.Load() && now.Sub(cs.reserved).Minutes() > 5 {
			m.remaining++
		}
	}
	m.lock.Unlock()
	return m.remaining, nil
}

func (m *MemTicketsDB) Reserve(fullname string) error {
	m.lock.RLock()
	customer, exists := m.Customers[fullname]
	m.lock.RUnlock()
	if !exists {
		m.lock.Lock()

		if m.remaining <= 0 {
			m.lock.Unlock()
			return ErrNoTicketsRemaining
		}
		m.remaining--
		m.Customers[fullname] = &customerState{
			reserved: time.Now().UTC(),
			charged:  atomic.NewBool(false),
		}

		m.lock.Unlock()
		return nil
	}
	if customer.charged.Load() {
		return ErrAlreadyCharged
	}
	customer.reserved = time.Now().UTC()
	return nil
}

func (m *MemTicketsDB) Charge(fullname string) error {
	m.lock.RLock()
	customer, exists := m.Customers[fullname]
	m.lock.RUnlock()
	if !exists {
		return ErrNoReservation
	}
	if time.Now().UTC().Sub(customer.reserved).Minutes() > 5 {
		m.lock.Lock()
		m.remaining++
		m.lock.Unlock()
		return ErrReservationExpired
	}
	if customer.charged.CAS(false, true) {
		return nil
	}
	return ErrAlreadyCharged
}

func (m *MemTicketsDB) Guests() ([]string, error) {
	names := []string{}
	m.lock.RLock()
	for fullname, cs := range m.Customers {
		if cs.charged.Load() {
			names = append(names, fullname)
		}
	}
	m.lock.RUnlock()
	return names, nil
}

func (m *MemTicketsDB) Reset() error {
	m.lock.Lock()
	m.remaining += len(m.Customers)
	m.Customers = map[string]*customerState{}
	m.lock.Unlock()
	return nil
}
