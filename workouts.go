package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func mustListWorkouts(ctx context.Context, userID int, limit, offset int, cookies []*http.Cookie) {
	log.Printf("Dumping workouts for user %d, limit %d, offset %d\n", userID, limit, offset)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://www.endomondo.com/rest/v1/users/%d/workouts/history?limit=%d&offset=%d", userID, limit, offset), nil)
	if err != nil {
		log.Fatalf("Error create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error get session: %v", err)
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	var listWorkoutsResp struct {
		Data   []workout `json:"data"`
		Paging struct {
			Next string `json:"next"`
		} `json:"paging"`
	}
	if err := json.Unmarshal(bb, &listWorkoutsResp); err != nil {
		log.Println(string(bb))
		log.Fatalf("Error decode workouts list: %v", err)
	}
	if len(listWorkoutsResp.Data) == 0 {
		return
	}

	filename := fmt.Sprintf("dump/%d/raw/%d-%d-%d.json", userID, userID, limit, offset)
	if err := ioutil.WriteFile(filename, bb, os.ModePerm); err != nil {
		log.Fatalf("Failed to write contents: %v", err)
	}
}

func dumpWorkout(ctx context.Context, userID, workoutID int, cookies []*http.Cookie) error {
	filename := fmt.Sprintf("dump/%d/gpx/%d.gpx", userID, workoutID)
	if fileExists(filename) {
		return nil
	}

	time.Sleep(time.Second)
	log.Printf("Dumping workout %d for user %d\n", workoutID, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://www.endomondo.com/rest/v1/users/%d/workouts/%d/export?format=GPX", userID, workoutID), nil)
	if err != nil {
		return fmt.Errorf("error create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error get session: %w", err)
	}
	defer resp.Body.Close()

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	if err := ioutil.WriteFile(filename, bb, os.ModePerm); err != nil {
		return fmt.Errorf("failed to write contents: %v", err)
	}
	return nil
}

func mustDumpGPX(ctx context.Context, userID int, dir string, sports []int, cookies []*http.Cookie) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %q: %w", path, err)
		}

		var dumpContents struct {
			Data []workout `json:"data"`
		}
		if err := json.Unmarshal(file, &dumpContents); err != nil {
			return fmt.Errorf("failed to parse file %q: %w", path, err)
		}

		for _, workout := range dumpContents.Data {

			if !contains(gpxSports, workout.Sport) {
				continue
			}

			if err := dumpWorkout(ctx, userID, workout.ID, cookies); err != nil {
				return fmt.Errorf("failed to dump workout %d: %w", workout.ID, err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Failed to walk through workout files: %v", err)
	}
}

type workout struct {
	ID             int    `json:"id"`
	Sport          int    `json:"sport"`
	LocalStartTime string `json:"local_start_time"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
