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

type habitsMsg []habit

func fetchHabits() tea.Msg {
	req, _ := http.NewRequest(http.MethodGet, "https://api.habitify.me/journal", nil)

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

	var responseObject response
	json.Unmarshal(responseData, &responseObject)

	return habitsMsg(responseObject.Habits)
}
