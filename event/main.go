package event

import (
	"errors"
	"time"
)

// Type is a higher level grouping
type Type uint32

// Key is a lower level grouping (multiple Keys for Type)
type Key uint32

// Kind is Type | Key
type Kind uint64

func NewKind(t Type, k Key) Kind {
	return Kind(uint64(t)<<32&0xffffffff00000000 | uint64(k))
}

func (k Kind) Key() Key {
	return Key(k & 0x00000000ffffffff)
}

func (k Kind) Type() Type {
	return Type(k >> 32)
}

// Interface is a generic type for the contextual data of the event.
type Interface interface {
}

// Object is a struct that holds all the context of an Event.
type Object struct {
	Timestamp time.Time
	Kind      Kind
	Context   Interface
}

func New(ctx Interface) *Object {
	return &Object{
		Timestamp: time.Now(),
		Context:   ctx,
	}
}

type NewEventFunc func() Interface

type wrapper struct {
	str     string
	newFunc NewEventFunc
}

var (
	registeredEvents = make(map[Kind]*wrapper)
)

func Register(kind Kind, str string, newFunc NewEventFunc) error {
	if _, exists := registeredEvents[kind]; exists {
		return errors.New("event kind already exists")
	}
	registeredEvents[kind] = &wrapper{str, newFunc}
	return nil
}
