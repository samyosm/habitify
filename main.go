package main

import (
	"fmt"
	"log"
	"os"

	"github.com/arkan/dotconfig"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	spinner spinner.Model
	list    list.Model
	habits  []habit
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Habits"
	l.Paginator.KeyMap.NextPage.SetEnabled(false)
	l.Paginator.KeyMap.PrevPage.SetEnabled(false)
	l.KeyMap.NextPage.SetEnabled(false)
	l.KeyMap.PrevPage.SetEnabled(false)

	return model{spinner: s, list: l}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchHabits)
}

func (m model) AdditionalShortHelpKeys() func() []key.Binding {
	increment := key.NewBinding()
	increment.SetKeys("l", "right")
	increment.SetHelp("→/l", "increment")

	decrement := key.NewBinding()
	decrement.SetKeys("h", "left")
	decrement.SetHelp("←/h", "decrement")

	skip := key.NewBinding()
	skip.SetKeys("s")
	skip.SetHelp("s", "skip")

	currentHabit := m.habits[m.list.Cursor()]
	if currentHabit.Progress.CurrentValue-1 < 0 {
		return func() []key.Binding {
			return []key.Binding{
				increment,
				skip,
			}
		}
	}

	// TODO: Disable skip when habit status is set to completed
	// TODO: Disable increment when status is set to completed
	// TODO: Appropriate keymap for bad habit (habit_type 2)

	return func() []key.Binding {
		return []key.Binding{
			increment,
			decrement,
			skip,
		}
	}
}

func (m *model) updateListItems() {
	var items []list.Item
	for _, habit := range m.habits {
		items = append(items, habit)
	}

	m.list.SetItems(items)
}

func (m *model) updateSelectedHabit(status string) {
	currentHabit := m.habits[m.list.Cursor()]
	if currentHabit.Status == status {
		m.habits[m.list.Cursor()].Status = "in_progress"
	} else {
		m.habits[m.list.Cursor()].Status = status
	}

	m.fetchAndNotify(func() int {
		return putHabitStatus(currentHabit.ID, status)
	})

	m.updateListItems()
}

func (m *model) completeSelectedHabit(value int) tea.Cmd {
	currentHabit := m.habits[m.list.Cursor()]
	newValue := currentHabit.Progress.CurrentValue + value
	if newValue < 0 {
		return m.list.NewStatusMessage("Can't have negative results...")
	}

	m.fetchAndNotify(func() int {
		return addHabitLog(currentHabit.ID, currentHabit.Goal.UnitType, value)
	})

	m.habits[m.list.Cursor()].Progress.CurrentValue = newValue

	if newValue == currentHabit.Progress.TargetValue {
		m.habits[m.list.Cursor()].Status = "completed"
	} else {
		m.habits[m.list.Cursor()].Status = "in_progress"
	}

	m.updateListItems()
	return nil
}

func (m *model) fetchAndNotify(a func() int) {
	code := a()
	if code == 200 {
		m.list.NewStatusMessage("Operation successfull")
	} else {
		m.list.NewStatusMessage("Operation failed")
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case habitsMsg:
		m.habits = msg
		m.updateListItems()
		m.list.AdditionalShortHelpKeys = m.AdditionalShortHelpKeys()

	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "l", "right":
			return m, m.completeSelectedHabit(1)
		case "h", "left":
			return m, m.completeSelectedHabit(-1)

		case "s":
			m.updateSelectedHabit("skipped")
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)
	if len(m.habits) > 0 {
		m.list.AdditionalShortHelpKeys = m.AdditionalShortHelpKeys()
	}
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
