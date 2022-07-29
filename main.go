package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

const (
	envEstuaryKeyKey     = "ESTUARY_KEY"
	envEstuaryShuttleKey = "ESTUARY_SHUTTLE"
	userAgent            = "github.com/Jorropo/estuary-push"
)

func main() {
	err := mainRet()
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func mainRet() error {
	key := os.Getenv(envEstuaryKeyKey)
	shuttle := os.Getenv(envEstuaryShuttleKey)

	if key == "" {
		return fmt.Errorf("error empty " + envEstuaryKeyKey + " envKey")
	}
	if shuttle == "" {
		return fmt.Errorf("error empty " + envEstuaryShuttleKey + " envKey")
	}

	if len(os.Args) != 2 {
		return fmt.Errorf("expected 1 argument (car path) got %d", len(os.Args))
	}

	car, err := os.Open(os.Args[1])
	if err != nil {
		return fmt.Errorf("opening the car file: %w", err)
	}

	stats, err := car.Stat()
	if err != nil {
		return fmt.Errorf("stating the car file: %w", err)
	}

	req, err := http.NewRequest("POST", "https://"+shuttle+"/content/add-car", car)
	if err != nil {
		return fmt.Errorf("creating the request failed: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/vnd.ipld.car")
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Length", strconv.FormatUint(uint64(stats.Size()), 10))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("posting failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non 200 result code: %d, body:\n%s", resp.StatusCode, string(b))
	}

	return nil
}
