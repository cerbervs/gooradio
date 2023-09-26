package filterscreen

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gitlab.com/AgentNemo/goradios"
)

type SearchMethod string

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
	height        int
	width         int
	filter        *Filter
	stations      []goradios.Station
	input         textinput.Model
	styles        Styles
	inputResponse string
	inputIdx      int
	filterForm    FilterForm
}

func NewModel(width int, height int) *model {
	ti := textinput.New()
	ti.CharLimit = 80
	ti.Placeholder = "Type here"
	ti.Focus()

	return &model{
		filter: new(Filter),
		width:  width,
		height: height,
		styles: *NewStyles(),
		input:  ti,
	}
}

type FilterOption struct {
	Title     string
	Options   []string
	OptionIdx int
	Response  interface{}
}

type FilterForm struct {
	Options []FilterOption
}

func NewFilterForm() FilterForm {
	return FilterForm{
		Options: []FilterOption{
			{
				Title: "Sort By",
				Options: []string{
					"Name",
					"Country",
					"Votes",
					"Clicks",
				},
			},
			{
				Title: "Reverse",
			},
			{
				Title: "Hide Broken",
			},
			{
				Title: "Exact",
			},
		},
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *model) Init() tea.Cmd {
	m.filterForm = NewFilterForm()
	return nil
}

func (m *model) Next() tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	if m.inputIdx < len(m.filterForm.Options)-1 {
		if len(m.filterForm.Options[m.inputIdx].Options) > 0 {
			m.filterForm.Options[m.inputIdx].OptionIdx++
		} else {
			m.inputIdx++
		}
	} else {
		cmd = m.input.Focus()
		cmds = append(cmds, cmd)
	}

	return tea.Batch(cmds...)
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
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			m.filter.Term = m.input.Value()
			m.inputResponse = m.input.Value()
			m.input.SetValue("")
		case "tab":
			m.Next()
		case "space":
			if len(m.filterForm.Options[m.inputIdx].Options) > 0 {
				m.filterForm.Options[m.inputIdx].Response = m.filterForm.Options[m.inputIdx].OptionIdx
			} else {
				m.filterForm.Options[m.inputIdx].Response = true
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
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

	options := func() string {
		var options string
		for _, option := range m.filterForm.Options {
			options = lipgloss.JoinHorizontal(
				lipgloss.Left,
				option.Title,
				lipgloss.JoinHorizontal(
					lipgloss.Left,
					option.Options[option.OptionIdx],
				),
			)
		}
		return options
	}()

	fmt.Println(options)

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.InputField.Render(m.input.View()),
	)
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}
