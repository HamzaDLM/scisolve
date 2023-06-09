package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	0: {
		"Polynomial Division",
		"Multiplying Binomials",
		"Inverse Variation",
	},
	1: {
		"DNA Concentration",
		"Punnett Square",
		"Allele Frequency",
	},
	2: {
		"Polar Moment of Inertia",
		"Projectile Motion",
	},
	3: {
		"Atomic Mass",
		"Effective Nuclear Charge",
	},
}

func selectCalculator(domainId, questionId int, m model) model {
	switch domainId {
	case 0: // Math
		switch questionId {
		case 0:

		}
	case 1:
		switch questionId {
		case 0:
			result := wrapperDNAConcentration(m)
			return result
		}
	}
	return m
}

// General stuff for styling the view
var (
	term    = termenv.EnvColorProfile()
	keyword = makeFgStyle("211")
	subtle  = makeFgStyle("241")
	dot     = colorFg(" • ", "236")
	title   = makeFgBgStyle("#ffffff", "#341fd1")

	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

type model struct {
	Choice       int
	ChoiceCalc   int
	ChosenDomain bool
	ChosenCalc   bool
	InsideCalc   bool
	Description  string
	Result       string
	Loaded       bool
	Quitting     bool
	FocusIndex   int
	Inputs       []textinput.Model
	cursorMode   textinput.CursorMode
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Main update function.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if (k == "q" || k == "ctrl+c") && !m.ChosenCalc {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	if m.ChosenCalc { // && !m.InsideCalc
		m = selectCalculator(m.Choice, m.ChoiceCalc, m)
	}
	if !m.ChosenDomain {
		return updateChoices(msg, m)
	} else if m.ChosenDomain && !m.ChosenCalc {
		return updateChosen(msg, m)
	}
	return updateArguments(msg, m)
}

// The main view, which just calls the appropriate sub-view
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	// When do we want to select calculator and enter data
	// When calc func is chosen and quitting/enter keys aren't entered
	if m.ChosenCalc { //&& !m.InsideCalc
		m = selectCalculator(m.Choice, m.ChoiceCalc, m)
	}
	if !m.ChosenDomain {
		s = choicesView(m)
	} else if m.ChosenDomain && !m.ChosenCalc {
		s = chosenView(m)
	} else {
		s = argumentView(m)
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
		case "esc":
			return m, tea.Quit
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
		case "esc", "back":
			m.ChosenDomain = false
			return m, nil
		}
	}
	return m, nil
}

func updateArguments(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// If we go back, reset all necessary values
			m.ChosenCalc = !m.ChosenCalc
			m.Inputs = make([]textinput.Model, 0)
			m.Description = ""
			m.Result = ""
		case "ctrl+c":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.FocusIndex == len(m.Inputs) {
				pr("======= Redirect to output page ========")
				pr(m.Inputs)
				// return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > len(m.Inputs) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Inputs)
			}

			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := 0; i <= len(m.Inputs)-1; i++ {
				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Inputs[i].Focus()
					m.Inputs[i].PromptStyle = focusedStyle
					m.Inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Inputs[i].Blur()
				m.Inputs[i].PromptStyle = noStyle
				m.Inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.Inputs {
		m.Inputs[i], cmds[i] = m.Inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
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
	tpl += subtle("j/k/up/down: select") + dot + subtle("enter: choose") + dot + subtle("q, esc: quit")

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
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("esc: back") + dot + subtle("q, ctrl+c: quit")

	return fmt.Sprintf(tpl, domainStore[c], colorFg("calculations", "79"))
}

// View for third section: Arguments
func argumentView(m model) string {
	var b strings.Builder

	if m.Description != "" {
		b.WriteString(m.Description)
		fmt.Fprintf(&b, "\n\n")
	}

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if m.Result != "" {
		fmt.Fprintf(&b, "\n\n")
		b.WriteString(m.Result)
	}

	fmt.Fprintf(&b, "\n\n\n")

	b.WriteString(helpStyle.Render("Press esc to quit"))

	return b.String()
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg(dot, "212") + colorFg(label, "212")
	}
	return fmt.Sprintf(label)
}

func main() {
	initialModel := model{
		Choice:       0,
		ChoiceCalc:   0,
		ChosenDomain: false,
		ChosenCalc:   false,
		InsideCalc:   false,
		Description:  "",
		Loaded:       false,
		Quitting:     false,
		Inputs:       make([]textinput.Model, 0),
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
