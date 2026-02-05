package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Modal defines the interface for our unified modal components.
type Modal interface {
	Init() tea.Cmd
	Update(tea.Msg) (Modal, tea.Cmd)
	View(width, height int) string
}

// ModalResult is a generic message sent when a modal finishes.
type ModalResult struct {
	Canceled bool
	Value    any
}

// --- Modal Frame ---

func renderModal(width, height int, title string, content string) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlight).
		Padding(1).
		Width(60)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlight).
		MarginBottom(1)

	header := titleStyle.Render(title)
	modal := style.Render(header + "\n" + content)

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, modal)
}

// --- Confirmation Modal ---

type ConfirmModal struct {
	Title  string
	Prompt string
}

func (m ConfirmModal) Init() tea.Cmd { return nil }
func (m ConfirmModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			return m, func() tea.Msg { return ModalResult{Value: true} }
		case "n", "N", "esc":
			return m, func() tea.Msg { return ModalResult{Canceled: true} }
		}
	}
	return m, nil
}
func (m ConfirmModal) View(w, h int) string {
	content := fmt.Sprintf("%s\n\n(y/n)", m.Prompt)
	return renderModal(w, h, m.Title, content)
}

// --- Text Input Modal ---

type TextInputModal struct {
	Title string
	Input textinput.Model
}

func NewTextInputModal(title, initialValue, placeholder string) TextInputModal {
	ti := textinput.New()
	ti.SetValue(initialValue)
	ti.Placeholder = placeholder
	ti.Focus()
	return TextInputModal{
		Title: title,
		Input: ti,
	}
}

func (m TextInputModal) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, func() tea.Msg { return ModalResult{Value: m.Input.Value()} }
		case "esc":
			return m, func() tea.Msg { return ModalResult{Canceled: true} }
		}
	}
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m TextInputModal) View(w, h int) string {
	content := m.Input.View()
	return renderModal(w, h, m.Title, content)
}

// --- List Selector Modal ---

type ListOption struct {
	ID    string
	Label string
}

type ListSelectorModal struct {
	Title   string
	Options []ListOption
	Cursor  int
}

func (m ListSelectorModal) Init() tea.Cmd { return nil }
func (m ListSelectorModal) Update(msg tea.Msg) (Modal, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Options)-1 {
				m.Cursor++
			}
		case "enter":
			if len(m.Options) > 0 {
				return m, func() tea.Msg { return ModalResult{Value: m.Options[m.Cursor].ID} }
			}
		case "esc":
			return m, func() tea.Msg { return ModalResult{Canceled: true} }
		}
	}
	return m, nil
}

func (m ListSelectorModal) View(w, h int) string {
	var b strings.Builder
	if len(m.Options) == 0 {
		b.WriteString("No options available.")
	} else {
		for i, opt := range m.Options {
			cursor := "  "
			if m.Cursor == i {
				cursor = "> "
			}
			style := lipgloss.NewStyle()
			if m.Cursor == i {
				style = style.Foreground(special)
			}
			b.WriteString(fmt.Sprintf("%s%s\n", cursor, style.Render(opt.Label)))
		}
	}
	return renderModal(w, h, m.Title, b.String())
}