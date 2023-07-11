package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func addHabitLog(habitId, unit_type string, value int) int {
	data := []byte(fmt.Sprintf(`{
		"unit_type": "%s",
		"value": %d,
		"target_date": "%s"
	}`, unit_type, value, time.Now().Format(time.RFC3339)))

	req, _ := http.NewRequest(http.MethodPost, "https://api.habitify.me/logs/"+habitId, bytes.NewBuffer(data))

	req.Header.Set("Authorization", getApiKey())
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to add log...")
	}

	if resp.StatusCode != 200 {
		log.Print(habitId)
		log.Print(unit_type)
		log.Print(time.Now().Format(time.RFC3339))

		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Fatal(string(responseData))
	}

	return resp.StatusCode
}

func putHabitStatus(habitId, status string) int {
	data := []byte(fmt.Sprintf(`{
		"status": "%s",
		"target_date": "%s"
	}`, status, time.Now().Format(time.RFC3339)))

	req, _ := http.NewRequest(http.MethodPut, "https://api.habitify.me/status/"+habitId, bytes.NewBuffer(data))

	req.Header.Set("Authorization", getApiKey())
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to fetch data...")
	}

	return resp.StatusCode
}
