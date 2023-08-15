package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Presence struct {
	ID                string `json:"id"`
	Availability      string `json:"availability"`
	Activity          string `json:"activity"`
	OutOfOfficeStatus struct {
		IsOutOfOffice bool `json:"isOutOfOffice"`
	} `json:"outOfOfficeSettings"`
}

func processObjectID(objectID string, writer *csv.Writer) {
	url := fmt.Sprintf("https://graph.office.net/en-us/graph/api/proxy?url=https%%3A%%2F%%2Fgraph.microsoft.com%%2Fbeta%%2Fusers%%2F%s%%2Fpresence", objectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.171 Safari/537.36")
	req.Header.Set("Authorization", "Bearer {token:https://graph.microsoft.com/}") // Replace with actual token
	req.Header.Set("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var presence Presence
	err = json.Unmarshal(body, &presence)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-40s %-15s %-15s %-5v\n", presence.ID, presence.Availability, presence.Activity, presence.OutOfOfficeStatus.IsOutOfOffice)

	if writer != nil {
		err = writer.Write([]string{presence.ID, presence.Availability, presence.Activity, fmt.Sprintf("%v", presence.OutOfOfficeStatus.IsOutOfOffice)})
		if err != nil {
			log.Fatalf("failed to write CSV row: %s", err)
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

	var writer *csv.Writer
	if *exportPath != "" {
		file, err := os.Create(*exportPath)
		if err != nil {
			log.Fatalf("failed to create file: %s", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		writer = csv.NewWriter(file)
		defer writer.Flush()

		err = writer.Write([]string{"Object ID", "Availability", "Activity", "Out of Office"})
		if err != nil {
			log.Fatalf("failed to write CSV header: %s", err)
		}
	}

	// Print the headers here
	fmt.Printf("%-40s %-15s %-15s %-5s\n", "Object ID", "Availability", "Activity", "Out of Office")

	if *objectID != "" {
		processObjectID(*objectID, writer)
	} else if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

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
	} else {
		fmt.Println("Please provide either a single object ID (-o) or a file path (-f)")
	}
}
