package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	spinner spinner.Model
	habits  []habit // items on the to-do list
	cursor  int     // which to-do list item our cursor is pointing at
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHabits)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.habits)-1 {
				m.cursor++
			}
		}
	case habitsMsg:
		m.habits = msg
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if len(m.habits) == 0 {
		s := fmt.Sprintf("\n\n   %s Loading habits...press q to quit\n\n", m.spinner.View())
		return s
	}

	// The header
	s := "Habits:\n\n"

	// Iterate over our choices
	for i, habit := range m.habits {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		switch habit.Status {
		case "completed":
			checked = "x"
		case "failed":
			checked = "f"
		case "in_progress":
			checked = "-"
		case "skipped":
			checked = "/"
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s (%d/%d %s)\n", cursor, checked, habit.Name, habit.Progress.CurrentValue, habit.Progress.TargetValue, habit.Progress.Periodicity)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
