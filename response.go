package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type response struct {
	Errors  []any   `json:"errors"`
	Message string  `json:"message"`
	Habits  []habit `json:"data"`
	Version string  `json:"version"`
	Status  bool    `json:"status"`
}

type habit struct {
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
}

var (
	// https://atlassian.design/foundations/color
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5630"))
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("#36B37E"))
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFC400"))
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("#5E6C84"))
)

func (i habit) Title() string { return i.Name }
func (i habit) Description() string {
	s := ""
	switch i.Status {
	case "completed":
		s += green.Render("Completed")
	case "in_progress":
		s += yellow.Render("In progress")
	case "failed":
		s += red.Render("Failed")
	}

	s += gray.Render(fmt.Sprintf(" (%d/%d)", i.Progress.CurrentValue, i.Progress.TargetValue))

	if i.Area.Name != "" {
		s += " • " + i.Area.Name
	}

	return s
}
func (i habit) FilterValue() string { return i.Name }
