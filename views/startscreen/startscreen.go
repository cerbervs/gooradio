package startscreen

import (
	"gooradio/views/filterscreen"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	height     int
	width      int
	searchterm string
	terminput  textinput.Model
	styles     *Styles
}

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func NewStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("#0099cc")
	s.InputField = lipgloss.NewStyle().
		BorderForeground(s.BorderColor).
		BorderStyle(lipgloss.NormalBorder()).
		Padding(1).
		Width(80)

	return s
}

func NewStartScreen() (*model, error) {
	styles := NewStyles()
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Search Term"
	ti.Focus()

	return &model{
			terminput: ti,
			styles:    styles,
		},
		nil
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.searchterm = m.terminput.Value()
			m.terminput.SetValue("")
			if filter, err := filterscreen.NewFilter(m.searchterm); err != nil {
				return m, nil
			} else {
				filterScreen := filterscreen.NewModel(filter, m.width, m.height)
				return filterScreen, nil
			}
		}
	}

	m.terminput, cmd = m.terminput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
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
			"Please enter your search term",
			m.styles.InputField.Render(m.terminput.View()),
		),
	)
}
