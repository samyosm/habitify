package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/arkan/dotconfig"
	tea "github.com/charmbracelet/bubbletea"
)

type habitsMsg []habit

type settings struct {
	ApiKey string `yaml:"api-key"`
}

func getApiKey() string {
	apiKey := os.Getenv("HABITIFY_API_KEY")
	if apiKey != "" {
		return apiKey
	}

	ss := settings{}

	if err := dotconfig.Load("habitify", &ss); err != nil {
		if err == dotconfig.ErrConfigNotFound {
			log.Fatal("Absent Habitify api key. Please declare a HABITIFY_API_KEY environment variable or add your api key in config.")
		}
	} else if err != nil {
		log.Fatal("Absent Habitify api key. Please declare a HABITIFY_API_KEY environment variable or add your api key in config.")
	}

	return ss.ApiKey
}

func fetchHabits() tea.Msg {
	req, _ := http.NewRequest(http.MethodGet, "https://api.habitify.me/journal", nil)

	q := req.URL.Query()
	q.Add("target_date", time.Now().Format(time.RFC3339))

	req.URL.RawQuery = q.Encode()

	req.Header.Add("Authorization", getApiKey())

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to fetch data...")
	}

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var responseObject response
	json.Unmarshal(responseData, &responseObject)

	return habitsMsg(responseObject.Habits)
}
