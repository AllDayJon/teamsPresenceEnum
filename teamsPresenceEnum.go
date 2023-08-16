package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type Presence struct {
	ID                string `json:"id"`
	Availability      string `json:"availability"`
	Activity          string `json:"activity"`
	OutOfOfficeStatus struct {
		IsOutOfOffice bool `json:"isOutOfOffice"`
	} `json:"outOfOfficeSettings"`
}

var client = &http.Client{}

func createRequest(objectID string) (*http.Request, error) {
	url := fmt.Sprintf("https://graph.office.net/en-us/graph/api/proxy?url=https%%3A%%2F%%2Fgraph.microsoft.com%%2Fbeta%%2Fusers%%2F%s%%2Fpresence", objectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.171 Safari/537.36")
	req.Header.Set("Authorization", "Bearer {token:https://graph.microsoft.com/}")
	req.Header.Set("Accept", "*/*")

	return req, nil
}

const (
	maxRetries  = 3
	baseDelayMs = 500 // 500 milliseconds
)

func processObjectID(objectID string, writer *csv.Writer) {
	var (
		resp *http.Response
		err  error
	)

	req, err := createRequest(objectID)
	if err != nil {
		log.Printf("Error creating request for object ID %s: %s\n", objectID, err)
		return
	}

	for i := 0; i <= maxRetries; i++ {
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if i < maxRetries {
			delay := time.Duration(baseDelayMs*math.Pow(2, float64(i))) * time.Millisecond
			time.Sleep(delay)
		}
	}

	if err != nil {
		log.Printf("Error making request for object ID %s after %d retries: %s\n", objectID, maxRetries, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for object ID %s: %s\n", objectID, err)
		return
	}

	var presence Presence
	err = json.Unmarshal(body, &presence)
	if err != nil {
		log.Printf("Error unmarshalling response for object ID %s: %s\n", objectID, err)
		return
	}

	fmt.Printf("%-40s %-15s %-15s %-5v\n", presence.ID, presence.Availability, presence.Activity, presence.OutOfOfficeStatus.IsOutOfOffice)

	if writer != nil {
		err = writer.Write([]string{presence.ID, presence.Availability, presence.Activity, fmt.Sprintf("%v", presence.OutOfOfficeStatus.IsOutOfOffice)})
		if err != nil {
			log.Printf("Failed to write CSV row for object ID %s: %s\n", objectID, err)
		}
	}
}

func main() {
	banner := `
  _______                         ______
 |__   __|                       |  ____|
    | | ___  __ _ _ __ ___  ___  | |__   _ __  _   _ _ __ ___
    | |/ _ \/ _` + "`" + ` | '_ ` + "`" + ` _ \/ __| |  __| | '_ \| | | | '_ ` + "`" + ` _ \
    | |  __/ (_| | | | | | \__ \ | |____| | | | |_| | | | | | |
    |_|\___|\__,_|_| |_| |_|___/ |______|_| |_|\__,_|_| |_| |_|`

	fmt.Println(banner)

	objectID := flag.String("o", "", "Single object ID")
	filePath := flag.String("f", "", "File path containing object IDs")
	exportPath := flag.String("path", "", "Path to export CSV file")

	flag.Parse()

	if *objectID == "" && *filePath == "" {
		fmt.Println("Please provide either a single object ID (-o) or a file path (-f)")
		return
	}

	var writer *csv.Writer
	if *exportPath != "" {
		file, err := os.Create(*exportPath)
		if err != nil {
			log.Fatalf("Failed to create file: %s", err)
		}
		defer file.Close()

		writer = csv.NewWriter(file)
		defer writer.Flush()

		err = writer.Write([]string{"Object ID", "Availability", "Activity", "Out of Office"})
		if err != nil {
			log.Fatalf("Failed to write CSV header: %s", err)
		}
	}

	fmt.Printf("%-40s %-15s %-15s %-5s\n", "Object ID", "Availability", "Activity", "Out of Office")

	if *objectID != "" {
		processObjectID(*objectID, writer)
	} else if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			id := strings.TrimSpace(scanner.Text())
			if id != "" {
				processObjectID(id, writer)
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}
