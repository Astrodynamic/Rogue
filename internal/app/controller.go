package app

// Controller is the interface that UI adapters depend on.
// This allows UI to work with App without knowing its concrete type.
type Controller interface {
	Handle(in Input) (ViewModel, bool, error)
	SetWorldSize(width, height int)
}
