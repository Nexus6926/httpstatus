package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Define command line flags
	filename := flag.String("dL", "", "Specify the file containing website URLs")
	flag.Parse()

	// Check if filename flag is provided
	if *filename == "" {
		fmt.Println("Please provide a filename using the -dL flag")
		return
	}

	// Open the file
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create output folder if it doesn't exist
	outputFolder := "output"
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.Mkdir(outputFolder, 0755)
		if err != nil {
			fmt.Println("Error creating output folder:", err)
			return
		}
	}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line in the file
	for scanner.Scan() {
		// Read the URL from the current line
		url := scanner.Text()

		// Check if the URL contains a valid scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}

		// Fetch URL and get status code
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			saveURLToFile(outputFolder+"/000.txt", url)
			continue
		}
		defer resp.Body.Close()

		// Create file for status code and write URL
		fileName := fmt.Sprintf("%s/urls_%03d.txt", outputFolder, resp.StatusCode)
		saveURLToFile(fileName, url)
	}

	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

// Function to save URL to a file
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
