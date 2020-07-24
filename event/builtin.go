package event

const (
	TypeData = Type(0)

	KeyCreate = Key(iota)
	KeyRead
	KeyUpdate
	KeyDelete
)

var (
	builtInEvents = map[Kind]*wrapper{
		NewKind(TypeData, KeyCreate): &wrapper{"data.create", NewDataEvent},
		NewKind(TypeData, KeyRead):   &wrapper{"data.read", NewDataEvent},
		NewKind(TypeData, KeyUpdate): &wrapper{"data.update", NewDataEvent},
		NewKind(TypeData, KeyDelete): &wrapper{"data.delete", NewDataEvent},
	}
)

type DataEvent struct {
	Old interface{}
	New interface{}
}

func NewDataEvent(old, newDat interface{}) Interface {
	return &DataEvent{Old: old, New: newDat}
}
