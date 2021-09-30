package messages

type Message struct {
	Key  MessageName
	Data map[string]string
}

type MessageName string

const (
	SetEditModeMsg    = MessageName("setEditMode")
	UpdateValueMsg    = MessageName("updateValue")
	SetExploreModeMsg = MessageName("setExploreMode")
)
