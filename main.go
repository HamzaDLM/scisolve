package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

// Define dict for domains and calculations available
var questionsStore = map[string][]string{
	"Math":    []string{"Question1", "Question2"},
	"Biology": []string{"Question1", "Question2"},
	"Physics": []string{"foiezfjeio"},
}

// Secodn approach to defining the domains and questions
// 1: "Math"
// 2: "Biology"

// func (domain) {
// 	switch {

// 	case "Math"
// 		return ["dfeo", "ifjefi"]
// 	}
// }

// General stuff for styling the view
var (
	term    = termenv.EnvColorProfile()
	keyword = makeFgStyle("211")
	subtle  = makeFgStyle("241")
	dot     = colorFg(" • ", "236")
	title   = makeFgBgStyle("#ffffff", "#341fd1")
)

type model struct {
	Choice       int
	ChosenDomain bool
	ChosenCalc   bool
	Loaded       bool
	Quitting     bool
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if !m.ChosenDomain {
		return updateChoices(msg, m)
	}
	return updateChosen(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.ChosenDomain {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return indent.String("\n"+s+"\n\n", 2)
}

// Sub-update functions

// Update loop for the first view where you're choosing a task.
func updateChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// TODO: make the length check dynamic
		case "j", "down":
			m.Choice++
			if m.Choice > 3 {
				m.Choice = 3
			}
		case "k", "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.ChosenDomain = true
			return m, nil
		}
	}
	return m, nil
}

// Update loop for the second view after a choice has been made
func updateChosen(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	}
	return m, nil
}

// Sub-views

// The first view, for choosing the domain
func choicesView(m model) string {
	c := m.Choice

	tpl := title("Welcome to SciSolve\n\n")
	tpl += "Available domain:\n\n"
	tpl += "%s\n\n"
	tpl += "Select to show the available %s for that domain.\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	ch := "%s\n%s\n%s\n"
	i := 0
	var checkboxes []any
	for key, _ := range questionsStore {
		checkboxes = append(checkboxes, checkbox(key, c == i))
		i++
	}
	choices := fmt.Sprintf(ch, checkboxes...)

	return fmt.Sprintf(tpl, choices, colorFg("calculations", "79"))
}

// The second view, to show available calculations
// func calculationsView(m model) string {

// 	tpl := title("Welcome to SciSolve\n\n")
// 	tpl += "Available calculations:\n\n"
// 	tpl += "%s\n\n"
// 	tpl += "Select to perform the %s.\n\n"
// 	// helper text
// 	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

// 	choices := fmt.Sprintf(
// 		"%s\n%s\n%s\n%s",
// 		"math",
// 		"eofzo",
// 		"fzioej",
// 		"fe",
// 		// checkbox("Math", c == 0),
// 		// checkbox("Biology", c == 1),
// 		// checkbox("Physics", c == 2),
// 		// checkbox("Chemistry", c == 3),
// 	)

// 	return fmt.Sprintf(tpl, choices, colorFg("calculations", "79"))
// }

// The third view, after a task has been chosen
func chosenView(m model) string {
	var msg string

	switch m.Choice {
	case 0:
		msg = fmt.Sprintf("Carrot planting?\n\nCool, we'll need %s and %s...", keyword("libgarden"), keyword("vegeutils"))
	case 1:
		msg = fmt.Sprintf("A trip to the market?\n\nOkay, then we should install %s and %s...", keyword("marketkit"), keyword("libshopping"))
	case 2:
		msg = fmt.Sprintf("Reading time?\n\nOkay, cool, then we'll need a library. Yes, an %s.", keyword("actual library"))
	default:
		msg = fmt.Sprintf("This domain is not implemented yet.\n\nPlease Check later!")
	}

	return msg + "\n\n"
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func main() {
	initialModel := model{0, false, false, false, false}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
