package usecase

// Controller is the inbound port for the presentation layer (TUI).
// The UI depends on this interface, not on the concrete Usecase struct.
type Controller interface {
	Handle(in Input) (ViewModel, bool, error)
	SetWorldSize(width, height int)
}
