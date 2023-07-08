package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arkan/dotconfig"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	spinner spinner.Model
	list    list.Model
	cursor  int // which to-do list item our cursor is pointing at
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Habits"

	return model{spinner: s, list: l}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHabits)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case habitsMsg:
		var items []list.Item
		for _, habit := range msg {
			items = append(items, habit)
		}
		m.list.SetItems(items)

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
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if len(m.list.Items()) == 0 {
		s := fmt.Sprintf("\n\n   %s Loading habits...press q to quit\n\n", m.spinner.View())
		return s
	}

	return m.list.View()
}

func main() {
	args := os.Args[1:]
	// TODO: A proper CLI with a help menu
	if len(args) == 2 && args[0] == "init" {
		ss := settings{
			ApiKey: args[1],
		}
		if err := dotconfig.Save("habitify", ss); err != nil {
			log.Fatal("An error occured while trying to set api key")
		} else {
			fmt.Println("Successfully set api key")
			return
		}
	} else if len(args) != 0 {
		log.Fatal("Unexpected arguments")
	}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
