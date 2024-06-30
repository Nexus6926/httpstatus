package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	numWorkers       = 20
	requestTimeout   = 10 * time.Second
	retryPause       = 2 * time.Second
	maxRetries       = 2
	outputFolderName = "output"
)

func main() {
	filename := flag.String("dL", "", "Specify the file containing website URLs")
	update := flag.Bool("up", false, "Update the script")
	flag.Parse()

	if *update {
		updateScript()
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

	if _, err := os.Stat(outputFolderName); os.IsNotExist(err) {
		err := os.Mkdir(outputFolderName, 0755)
		if err != nil {
			fmt.Println("Error creating output folder:", err)
			return
		}
	}

	urls := make(chan string, 1000)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(urls, &wg)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "http://" + url
		}
		urls <- url
	}

	close(urls)
	wg.Wait()

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func worker(urls chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{Timeout: requestTimeout}

	for url := range urls {
		for tries := 0; tries < maxRetries; tries++ {
			err := fetchAndSaveURL(client, url)
			if err == nil {
				break
			}
			if strings.Contains(err.Error(), "no such host") {
				fmt.Printf("Error fetching url %s: %v\n", url, err)
				break
			}
			fmt.Printf("Error fetching url %s: %v\n", url, err)
			time.Sleep(retryPause)
		}
	}
}

func fetchAndSaveURL(client *http.Client, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileName := fmt.Sprintf("%s/urls_%03d.txt", outputFolderName, resp.StatusCode)
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

func updateScript() {
	fmt.Println("Script is up to date.")
}
