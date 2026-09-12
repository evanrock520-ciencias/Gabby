package panels

import (
	"client/internal/ui/messages"

	tea "github.com/charmbracelet/bubbletea"
)

type WizardStore struct {
	choices map[string][]string
}

func newWizardStore() WizardStore {
	return WizardStore{choices: make(map[string][]string)}
}

func (store *WizardStore) Set(key string, values []string) {
	store.choices[key] = values
}

func (store *WizardStore) Get(key string) ([]string, bool) {
	value, ok := store.choices[key]
	return value, ok
}

type WizardModal struct {
	Layout
	store       WizardStore
	steps       []Modal
	currentStep int
	onDone      func(WizardStore) messages.ModalResultMsg
	done        bool
}

func NewWizardModal(factory func(store *WizardStore) []Modal, onDone func(WizardStore) messages.ModalResultMsg) *WizardModal {
	w := &WizardModal{onDone: onDone}
	w.store = newWizardStore()
	w.steps = factory(&w.store)
	return w
}

func (w *WizardModal) Init() tea.Cmd {
	return nil
}

func (w *WizardModal) IsCapturingInput() bool {
	return !w.done
}

func (w *WizardModal) Focus() tea.Cmd {
	if len(w.steps) == 0 {
		return nil
	}
	return w.steps[0].Focus()
}

func (w *WizardModal) SetSize(width, height int) {
	w.Layout.SetSize(width, height)
	for _, step := range w.steps {
		step.SetSize(width, height)
	}
}

func (w *WizardModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	if w.done || w.currentStep >= len(w.steps) {
		return w, nil
	}

	updated, cmd := w.steps[w.currentStep].Update(msg)
	w.steps[w.currentStep] = updated

	if !updated.IsCapturingInput() {

		if cmd == nil {
			w.done = true
			return w, nil
		}

		w.currentStep++
		if w.currentStep >= len(w.steps) {
			w.done = true
			return w, tea.Sequence(cmd, func() tea.Msg {
				return w.onDone(w.store)
			})
		}

		focusCmd := w.steps[w.currentStep].Focus()
		return w, tea.Batch(cmd, focusCmd)
	}

	return w, cmd
}

func (w *WizardModal) View() string {
	if w.currentStep >= len(w.steps) {
		return ""
	}
	return w.steps[w.currentStep].View()
}
