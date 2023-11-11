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
	Offset     uint
	Limit      uint
	Reverse    bool
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

type FilterOption struct {
	Title     string
	Options   []string
	OptionIdx int
	Response  int
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

type model struct {
	input      textinput.Model
	filter     *Filter
	styles     Styles
	stations   []goradios.Station
	filterForm FilterForm
	height     int
	width      int
	inputIdx   int
}

func NewModel(width int, height int) *model {
	ti := textinput.New()
	ti.CharLimit = 80
	ti.Placeholder = "Type here"

	return &model{
		filter:     new(Filter),
		width:      width,
		height:     height,
		styles:     *NewStyles(),
		input:      ti,
		filterForm: NewFilterForm(),
	}
}

// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *model) Init() tea.Cmd {
	m.filterForm = NewFilterForm()
	return nil
}

func (m *model) MoveSelection(next bool, prev bool) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	if m.inputIdx < len(m.filterForm.Options)-1 {
		m.input.Blur()
		if next {
			if m.filterForm.Options[m.inputIdx].OptionIdx < len(
				m.filterForm.Options[m.inputIdx].Options,
			)-1 {
				m.filterForm.Options[m.inputIdx].OptionIdx++
				return tea.Batch(cmds...)
			}

			if m.inputIdx < len(m.filterForm.Options)-1 {
				m.inputIdx++
			}
		} else {
			if m.inputIdx == 0 {
				m.inputIdx = len(m.filterForm.Options)
				cmd = m.input.Focus()
				cmds = append(cmds, cmd)
				return tea.Batch(cmds...)
			}

			if m.filterForm.Options[m.inputIdx].OptionIdx > 0 {
				m.filterForm.Options[m.inputIdx].OptionIdx--
				return tea.Batch(cmds...)
			}

			if m.inputIdx > 0 {
				m.inputIdx--
			}
		}
	} else {
		if next {
			cmd = m.input.Focus()
			m.inputIdx = 0
			for i := range m.filterForm.Options {
				m.filterForm.Options[i].OptionIdx = 0
			}
		} else {
			m.input.Blur()
			m.inputIdx--
		}
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
			if m.input.Focused() {
				m.filter.Term = m.input.Value()
				m.input.SetValue("")
			} else {
				if len(m.filterForm.Options[m.inputIdx].Options) > 0 {
					m.filterForm.Options[m.inputIdx].Response = m.filterForm.Options[m.inputIdx].OptionIdx
				} else {
					m.filterForm.Options[m.inputIdx].Response = 1
				}
			}
		case "tab":
			m.MoveSelection(true, false)
		case "shift+tab":
			m.MoveSelection(false, true)
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
		for i, option := range m.filterForm.Options {
			options = lipgloss.JoinVertical(
				lipgloss.Top,
				options,
				func() string {
					var (
						response string
						selected string
					)
					if len(option.Options) > 0 {
						var subopts string
						subopts = option.Title
						for j, subopt := range option.Options {
							if j == option.OptionIdx {
								selected = " >"
							} else {
								selected = " "
							}

							if j == option.Response {
								response = "x"
							} else {
								response = " "
							}

							subopts = lipgloss.JoinHorizontal(
								lipgloss.Right,
								subopts,
								fmt.Sprintf("%s [%s] %v", selected, response, subopt),
							)
						}
						return subopts
					}

					if i == m.inputIdx {
						selected = ">"
					} else {
						selected = " "
					}

					if option.Response == 1 {
						response = "x"
					} else {
						response = " "
					}

					return fmt.Sprintf("%s [%s] %s", selected, response, option.Title)
				}(),
			)
		}
		return options
	}()

	view := lipgloss.JoinVertical(
		lipgloss.Center,
		options,
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
