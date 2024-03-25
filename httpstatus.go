package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	filename := flag.String("dL", "", "Specify the file containing website URLs")
	updateFlag := flag.Bool("up", false, "Update the script")
	flag.Parse()

	if *updateFlag {
		fmt.Println("Updating script...")
		// Code to update the script
		return
	}

	if *filename == "" {
		fmt.Println("Please provide a filename using the -dL flag")
		return
	}

	file, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	outputFolder := "output"
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.Mkdir(outputFolder, 0755)
		if err != nil {
			fmt.Println("Error creating output folder:", err)
			return
		}
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()

		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		tries := 0
		for tries < 2 {
			err := fetchAndSaveURL(url, outputFolder)
			if err == nil {
				break // If fetched successfully, move to the next URL
			}

			if strings.Contains(err.Error(), "dial tcp: lookup") && strings.Contains(err.Error(), "no such host") {
				fmt.Println("Skipping to next URL after 2 attempts due to DNS lookup error.")
				break
			}

			fmt.Printf("Error fetching url: %v\n", err)
			fmt.Println("Pausing for 5 seconds before retrying...")
			time.Sleep(5 * time.Second) // Pause for 5 seconds before retrying
			tries++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func fetchAndSaveURL(url, outputFolder string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err // Return the error if unable to fetch URL
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("%s/urls_%03d.txt", outputFolder, resp.StatusCode)
	saveURLToFile(fileName, url)

	return nil
}

func saveURLToFile(fileName, url string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", fileName, err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", url)
	if err != nil {
		fmt.Printf("Error writing URL to file %s: %v\n", fileName, err)
	}
}
