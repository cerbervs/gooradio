package filterscreen

import (
	tea "github.com/charmbracelet/bubbletea"
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

func NewFilter(term string) (Filter, error) {
	return Filter{
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
	filter   Filter
	stations []goradios.Station
}

func NewModel(filter Filter) *model {
	return &model{
		filter: filter,
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
		case "ctrl+c", "q":
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
	panic("not implemented") // TODO: Implement
}
