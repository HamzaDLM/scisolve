package main

import (
	"fmt"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/termenv"
)

// Define dict for domains and calculations available
var domainStore = map[int]string{
	0: "Math",
	1: "Biology",
	2: "Physics",
	3: "Chemistry",
}
var questionStore = map[int][]string{
	0: {"Polynomial Division", "Multiplying Binomials", "Inverse Variation", "Inverse Variation"},
	1: {"DNA Concentration", "Punnett Square", "Allele Frequency"},
	2: {"Polar Moment of Inertia", "Projectile Motion", "Projectile Motion"},
	3: {"Atomic Mass", "Effective Nuclear Charge", "Effective Nuclear Charge"},
}

// Helper function to print
func pr(arg any) {
	fmt.Println(arg)
}

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
	ChoiceCalc   int
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
		case "j", "down":
			m.Choice++
			if m.Choice > len(domainStore)-1 {
				m.Choice = len(domainStore) - 1
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
	c := m.Choice
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.ChoiceCalc++
			if m.ChoiceCalc > len(questionStore[c])-1 {
				m.ChoiceCalc = len(questionStore[c]) - 1
			}
		case "k", "up":
			m.ChoiceCalc--
			if m.ChoiceCalc < 0 {
				m.ChoiceCalc = 0
			}
		case "enter":
			m.ChosenCalc = true
			return m, nil
			// case "b"
		}
	}
	return m, nil
}

// Sub-views

// The first view, for choosing the domain
func choicesView(m model) string {
	c := m.Choice

	tpl := title("SciSolve\n\n")
	tpl += "Available domain:\n\n"

	keys := make([]int, 0)
	for k, _ := range domainStore {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		tpl += checkbox(domainStore[k], c == k) + "\n"
	}
	tpl += "\n"
	tpl += "Select to show the available %s for that domain.\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	return fmt.Sprintf(tpl, colorFg("calculations", "79"))
}

// The second view, after a domain has been chosen
func chosenView(m model) string {
	c := m.Choice
	n := m.ChoiceCalc

	tpl := title("SciSolve\n\n")
	tpl += "Available %s calculators:\n\n"

	keys := make([]int, 0)
	for k, _ := range questionStore[c] {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		tpl += checkbox(questionStore[c][k], n == k) + "\n"
	}

	tpl += "\n"
	tpl += "Select to show the available %s for that domain.\n\n"
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

	return fmt.Sprintf(tpl, domainStore[c], colorFg("calculations", "79"))
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func main() {
	initialModel := model{0, 0, false, false, false, false}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
