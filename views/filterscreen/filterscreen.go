package filterscreen

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/AgentNemo/goradios"
)

type SearchMethod string

type Filter struct {
	By         goradios.StationsBy
	Term       string
	Order      goradios.Order
	Reverse    bool
	Offset     uint
	Limit      uint
	HideBroken bool
}

func NewFilter(term string) (*Filter, error) {
	return &Filter{
			By:         goradios.StationsByName,
			Term:       term,
			Order:      goradios.OrderName,
			Reverse:    false,
			Offset:     0,
			Limit:      0,
			HideBroken: true,
		},
		nil
}

type model struct {
	height   int
	width    int
	filter   *Filter
	stations []goradios.Station
}

func NewModel(filter *Filter, width int, height int) *model {
	return &model{
		filter: filter,
		width:  width,
		height: height,
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, nil
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			lipgloss.JoinVertical(
				lipgloss.Center,
				"Loading...",
			),
		)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		"Term: "+m.filter.Term,
		"Reverse: "+lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render(strconv.FormatBool(m.filter.Reverse)),
		"Offset: "+lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render(strconv.Itoa(int(m.filter.Offset))),
		"Limit: "+lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render(strconv.Itoa(int(m.filter.Limit))),
		"HideBroken: "+lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Render(strconv.FormatBool(m.filter.HideBroken)),
	)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}
