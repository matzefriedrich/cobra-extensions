package reflection

import (
	"sync"
)

// UidSequence represents an interface for generating unique sequential identifiers.
type UidSequence interface {

	// Next generates the next unique sequential identifier in the sequence and returns it as an unsigned 64-bit integer.
	Next() uint64
}

type uidSequence struct {
	value uint64
	m     sync.RWMutex
}

// globalUidSequence is a singleton instance of UidSequence used to generate globally unique sequential identifiers.
var globalUidSequence UidSequence = &uidSequence{}

// UidSequenceOption represents a functional option used to customize the behavior of an uidSequence instance.
type UidSequenceOption func(sequence *uidSequence)

// WithInitialValue returns an option function that sets the initial value for the given uidSequence to the provided unsigned 64-bit integer.
func WithInitialValue(value uint64) UidSequenceOption {
	return func(sequence *uidSequence) {
		sequence.value = value
	}
}

// NewUidSequence creates and returns a new UidSequence instance with optional customizable behavior through UidSequenceOptions.
func NewUidSequence(option ...UidSequenceOption) UidSequence {
	instance := &uidSequence{
		value: 0,
		m:     sync.RWMutex{},
	}
	instance.m.Lock()
	defer instance.m.Unlock()
	for _, option := range option {
		option(instance)
	}
	return instance
}

// Next increments the sequence value in a thread-safe manner and returns the updated value.
func (s *uidSequence) Next() uint64 {
	s.m.Lock()
	defer s.m.Unlock()
	s.value++
	return s.value
}
