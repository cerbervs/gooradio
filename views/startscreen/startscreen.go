package startscreen

import (
	"gooradio/views/filterscreen"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	height     int
	width      int
	searchterm string
}

func NewStartScreen() (*model, error) {
	return &model{}, nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			return filterscreen.NewModel(m.width, m.height), nil
		}
	}

	return m, nil
}

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Welcome to GooRadio!",
			"Press enter to continue or Ctrl-C to quit...",
		),
	)
}
