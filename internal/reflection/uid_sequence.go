package reflection

import (
	"sync"
)

type UidSequence interface {
	Next() uint64
}

type uidSequence struct {
	value uint64
	m     sync.RWMutex
}

var globalUidSequence UidSequence = &uidSequence{}

type UidSequenceOption func(sequence *uidSequence)

func WithInitialValue(value uint64) UidSequenceOption {
	return func(sequence *uidSequence) {
		sequence.value = value
	}
}

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

func (s *uidSequence) Next() uint64 {
	s.m.Lock()
	defer s.m.Unlock()
	s.value++
	return s.value
}
