package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	spinner spinner.Model
	table   table.Model
	cursor  int // which to-do list item our cursor is pointing at
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	columns := []table.Column{
		{Title: "Status", Width: 12},
		{Title: "Name", Width: 50},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	st := table.DefaultStyles()
	st.Header = st.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	st.Selected = st.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(st)

	return model{spinner: s, table: t}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHabits)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case habitsMsg:
		var rows []table.Row
		for _, habit := range msg {
			rows = append(rows, table.Row{habit.Status, habit.Name})
		}
		m.table.SetRows(rows)
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if len(m.table.Rows()) == 0 {
		s := fmt.Sprintf("\n\n   %s Loading habits...press q to quit\n\n", m.spinner.View())
		return s
	}

	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
