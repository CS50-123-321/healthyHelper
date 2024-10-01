package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	Intro           string = "I’m working with my eight nephews and family members to help them grow in different areas like confidence, English, and most recently, healthy habits—especially daily sports. Since they spend a lot of time on their phones, I’ve developed a Telegram mini app that tracks their habit streaks and helps them stay committed to their goals. The goal is to not just have fun with daily sports but to understand the long-term benefits of living a healthy life. Through daily messages, they’re encouraged to stay motivated, become consistent, and realize that by building good habits, they’re shaping their future selves in a positive way. Now"
	Topups          string = "Please create an original Arabic version of the motivational message that captures the same essence as the English version but is distinct and not a translation. Ensure it reflects an Iraqi spirit and is culturally relevant, use simple Arabic. Get creative with emojies. USE ONLY THE FOLLOWING SYMBOLS IN THE TEXT:  '!',',','?','#'"
	maxRequestTries int    = 3
)

func GetAiResponse(habit Habit) (r string, err error) {
	if maxRequestTries == 0 {
		fmt.Println("GetAiResponse, reaching the max tried")
		return "", err
	}
	r, err = gneratePersonlizeResponse(habit)
	if err != nil {
		fmt.Println("GetAiResponse, err", maxRequestTries, err)
		time.Sleep(5) // 5 sec wait to do thenext requesttt
		return GetAiResponse(habit)
	}
	return r, nil
}

type Candidate struct {
	Content struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"content"`
}

type APIResponse struct {
	Candidates []Candidate `json:"candidates"`
}

func gneratePersonlizeResponse(habit Habit) (string, error) {
	prompt := generateHabitPrompt(habit)
	Gem_Token := os.Getenv("GEMINI_API_TOKEN")
	fmt.Println("Preparing request...")
	data := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	postBody, _ := json.Marshal(data)
	responseBody := bytes.NewBuffer(postBody)
	fmt.Println("Sending request...")
	resp, err := http.Post(fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", Gem_Token), "application/json", responseBody)
	if err != nil {
		return "", fmt.Errorf("An Error Occurred: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Reading response...")

	// Read and parse the response body

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read the response body for error details
		return "", fmt.Errorf("Request failed with status: %s, response: %s", resp.Status, string(bodyBytes))
	}
	body, _ := io.ReadAll(resp.Body)
	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", err
	}

	// Extract and return the text from the first candidate
	if len(apiResponse.Candidates) > 0 && len(apiResponse.Candidates[0].Content.Parts) > 0 {
		return apiResponse.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("No valid response received")
}

func generateHabitPrompt(habit Habit) string {
	insight := fmt.Sprintf(
		" %s Here's how they've been doing: %s has been working on '%s' for %d days since %s. They've completed %d days out of their %d-day goal, with their highest streak being %d days! ",
		Intro, habit.Name, habit.HabitName, habit.TotalDays, habit.CreatedAt.Format("Jan 2, 2006"), habit.Streaked, habit.CommitmentPeriod, habit.TopHit)
	detailedLogs := generateDayLogs(habit.DaysLog)
	prompt := fmt.Sprintf(
		"%s\n%s\nPlease generate a fun readable motivational and fun teenagers friendly message in English to encourage %s to keep going! also include insights for their records. %s",
		insight, detailedLogs, habit.Name, Topups)
	return prompt
}

func generateDayLogs(daysLog map[string]bool) string {
	logSummary := "Progress over the last few days:\n"
	for day, completed := range daysLog {
		if completed {
			logSummary += fmt.Sprintf("Day %d: Completed\n", day)
		} else {
			logSummary += fmt.Sprintf("Day %d: Missed\n", day)
		}
	}
	return logSummary
}

func GenerateText(prompt string) (string, error) {
	Gem_Token := os.Getenv("GEMINI_API_TOKEN")
	fmt.Println("Preparing request...")
	data := map[string]interface{}{
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	postBody, _ := json.Marshal(data)
	responseBody := bytes.NewBuffer(postBody)
	fmt.Println("Sending request...")
	resp, err := http.Post(fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=%s", Gem_Token), "application/json", responseBody)
	if err != nil {
		return "", fmt.Errorf("An Error Occurred: %v", err)
	}
	defer resp.Body.Close()
	fmt.Println("Reading response...")

	// Read and parse the response body

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read the response body for error details
		return "", fmt.Errorf("Request failed with status: %s, response: %s", resp.Status, string(bodyBytes))
	}
	body, _ := io.ReadAll(resp.Body)
	var apiResponse APIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", err
	}

	// Extract and return the text from the first candidate
	if len(apiResponse.Candidates) > 0 && len(apiResponse.Candidates[0].Content.Parts) > 0 {
		return apiResponse.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("No valid response received")
}