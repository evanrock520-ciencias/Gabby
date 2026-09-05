package panels

// InputCapturer define el comportamiento de paneles que pueden capturar
// entrada de texto (por ejemplo, textinput) y requieren suprimir atajos globales.
type InputCapturer interface {
	IsCapturingInput() bool
}
