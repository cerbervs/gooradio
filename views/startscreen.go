package views

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type SearchMethod string

const (
	SearchTermAll                   SearchMethod = "All"
	SearchTermDetailedAll                        = "Detailed All"
	SearchTermTerm                               = "Term"
	SearchTermTermDetailed                       = "Term Detailed"
	SearchTermTermDetailedAll                    = "Term Detailed All"
	SearchTermCountry                            = "Country"
	SearchTermCountriesDetailed                  = "Countries Detailed"
	SearchTermCountriesCode                      = "Countries Code"
	SearchTermCountriesCodeDetailed              = "Countries Code Detailed"
	SearchTermEmpty                              = ""
)

type model struct {
	searchterm   string
	searchmethod SearchMethod
	textinput    textinput.Model
}

func NewStartScreen() (*model, error) {
	ti := textinput.New()
	ti.CharLimit = 30
	ti.Placeholder = "Search Term"

	return &model{textinput: ti}, nil
}

func (m model) Init() tea.Cmd {
	// TODO: Implement
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: Implement
	return m, nil
}

func (m model) View() string {
	b := &strings.Builder{}
	b.WriteString("Please enter your search term")
	b.WriteString(m.textinput.View())
	return b.String()
}
