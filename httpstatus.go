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
	filename := flag.String("dL", "", "Specify the file containing website URLs")
	flag.Parse()

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

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error fetching %s: %v\n", url, err)
			saveURLToFile(outputFolder+"/000.txt", url)
			continue
		}
		defer resp.Body.Close()

		fileName := fmt.Sprintf("%s/urls_%03d.txt", outputFolder, resp.StatusCode)
		saveURLToFile(fileName, url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
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
