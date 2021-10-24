package messages

type Message struct {
	Key  MessageName
	Data map[string]string
}

type MessageName string

const (
	SetEditModeMsg = MessageName("setEditMode")
	UpdateValueMsg = MessageName("updateValue")

	SetExploreModeMsg = MessageName("setExploreMode")

	ShowHelpWindow  = MessageName("showHelpWindow")
	CloseHelpWindow = MessageName("closeHelpWindow")

	SetSaveAsModeMsg = MessageName("setSaveAsMode")
	SaveAsMsg        = MessageName("saveAs")

	ShowToastMsg = MessageName("showToast")
	HideToastMsg = MessageName("hideToast")

	SetAddColumnModeMsg = MessageName("setAddColumnMode")
	AddColumnMsg        = MessageName("addColumn")
)
