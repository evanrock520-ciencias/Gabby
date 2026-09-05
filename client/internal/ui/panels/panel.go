package panels

// InputCapturer define el comportamiento de paneles que pueden capturar
// entrada de texto (por ejemplo, textinput) y requieren suprimir atajos globales.
type InputCapturer interface {
	IsCapturingInput() bool
}

type Layout struct {
	width  int
	height int
}

func (layout *Layout) SetSize(width int, height int) {
	layout.width = width
	layout.height = height
}

type Panel struct {
	Layout
	focus  bool
	cursor int
}

func (panel *Panel) SetFocus(value bool) {
	panel.focus = value
}

func (panel *Panel) IsFocused() bool {
	return panel.focus
}
