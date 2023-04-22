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

/*
Returns the target calculator function depending on domain and question id

and the number of arguments to be passed to it
*/
func selectCalculator(domainId, questionId int, m model) model {
	switch domainId {
	case 0: //
		switch questionId {
		case 0:
			result := test(m)
			return result
		}
	}
	return m
}

// TODO: remove | Helper function to print
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

	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type model struct {
	Choice       int
	ChoiceCalc   int
	ChosenDomain bool
	ChosenCalc   bool
	InsideCalc   bool
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
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}
	// When do we want to select calculator and enter data
	// When calc func is chosen and quitting/enter keys aren't entered
	if m.ChosenCalc && !m.InsideCalc {
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
		case "b", "back":
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
		case "ctrl+c", "esc":
			return m, tea.Quit

		// FIXME: to remove | Change cursor mode
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.Inputs))
			for i := range m.Inputs {
				cmds[i] = m.Inputs[i].SetCursorMode(m.cursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.FocusIndex == len(m.Inputs) {
				pr("======= OUTPUT ========")
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
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("b: back") + dot + subtle("q, esc: quit")

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
	tpl += subtle("j/k, up/down: select") + dot + subtle("enter: choose") + dot + subtle("b: back") + dot + subtle("q, esc: quit")

	return fmt.Sprintf(tpl, domainStore[c], colorFg("calculations", "79"))
}

// View for third section: Arguments
func argumentView(m model) string {
	var b strings.Builder

	for i := range m.Inputs {
		b.WriteString(m.Inputs[i].View())
		if i < len(m.Inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.FocusIndex == len(m.Inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

	return b.String()
}

func checkbox(label string, checked bool) string {
	if checked {
		return colorFg("[x] "+label, "212")
	}
	return fmt.Sprintf("[ ] %s", label)
}

func main() {
	initialModel := model{
		Choice:       0,
		ChoiceCalc:   0,
		ChosenDomain: false,
		ChosenCalc:   false,
		InsideCalc:   false,
		Loaded:       false,
		Quitting:     false,
		Inputs:       make([]textinput.Model, 0),
	}

	// var t textinput.Model
	// for i := range initialModel.Inputs {
	// 	t = textinput.New()
	// 	t.CursorStyle = cursorStyle
	// 	t.CharLimit = 32

	// 	switch i {
	// 	case 0:
	// 		t.Placeholder = "Nickname"
	// 		t.Focus()
	// 		t.PromptStyle = focusedStyle
	// 		t.TextStyle = focusedStyle
	// 	case 1:
	// 		t.Placeholder = "Email"
	// 		t.CharLimit = 64
	// 	case 2:
	// 		t.Placeholder = "Password"
	// 		t.EchoMode = textinput.EchoPassword
	// 		t.EchoCharacter = '•'
	// 	}

	// 	initialModel.Inputs[i] = t
	// }

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
