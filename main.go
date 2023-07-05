package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Response struct {
	Errors  []any  `json:"errors"`
	Message string `json:"message"`
	Data    []struct {
		ID         string    `json:"id"`
		Name       string    `json:"name"`
		IsArchived bool      `json:"is_archived"`
		StartDate  time.Time `json:"start_date"`
		TimeOfDay  []string  `json:"time_of_day"`
		Goal       struct {
			UnitType    string `json:"unit_type"`
			Value       int    `json:"value"`
			Periodicity string `json:"periodicity"`
		} `json:"goal"`
		GoalHistoryItems []struct {
			UnitType    string `json:"unit_type"`
			Value       int    `json:"value"`
			Periodicity string `json:"periodicity"`
		} `json:"goal_history_items"`
		LogMethod  string `json:"log_method"`
		Recurrence string `json:"recurrence"`
		Remind     []any  `json:"remind"`
		Area       struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Priority string `json:"priority"`
		} `json:"area"`
		CreatedDate time.Time `json:"created_date"`
		Priority    int64     `json:"priority"`
		Status      string    `json:"status"`
		Progress    struct {
			CurrentValue  int       `json:"current_value"`
			TargetValue   int       `json:"target_value"`
			UnitType      string    `json:"unit_type"`
			Periodicity   string    `json:"periodicity"`
			ReferenceDate time.Time `json:"reference_date"`
		} `json:"progress"`
		HabitType int `json:"habit_type"`
	} `json:"data"`
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

type model struct {
	choices  []string         // items on the to-do list
	cursor   int              // which to-do list item our cursor is pointing at
	selected map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	req, err := http.NewRequest(http.MethodGet, "https://api.habitify.me/journal", nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("target_date", "2023-07-05T00:00:00+00:00")

	req.URL.RawQuery = q.Encode()

	apiKey := os.Getenv("HABITIFY_API_KEY")
	if apiKey == "" {
		log.Fatal("Absent Habitify api key. Have you forgotten to declare a HABITIFY_API_KEY environment variable?")
		os.Exit(1)
	}

	req.Header.Add("Authorization", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to fetch data...")
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	choices := []string{}
	selected := make(map[int]struct{})
	for index, element := range responseObject.Data {
		choices = append(choices, element.Name)
		if element.Status == "completed" {
			selected[index] = struct{}{}
		}
	}

	return model{
		choices:  choices,
		selected: selected,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Is it a key press?
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Habits:\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
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
