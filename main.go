package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"test-stress-new/config"
	"time"
)

func main() {
	var (
		baseURL     = config.Config.BaseURL
		numRequests = config.Config.NumRequests
		concurrency = config.Config.Concurrency
	)

	// Read URL from txt file
	urlData, lengthUrlData := readURLFromTXT()

	// Read token from txt file
	tokenData, lengthData := readTokenFromTXT()
	if lengthData != 0 {
		concurrency = lengthData
	}

	// Create a pool of goroutines to send requests concurrently
	for n := 0; n < numRequests; n++ {
		for u := 0; u < lengthUrlData; u++ {
			// === Timestamp
			timestamp := time.Now().Format("2006-01-02 15:04:05")

			// === Print batch request
			fmt.Printf("\n\n\n")
			fmt.Println("==================")
			fmt.Println("code ", " token ", "time")
			fmt.Println("==================")

			// === Create a wait group to monitor goroutines
			var wg *sync.WaitGroup = new(sync.WaitGroup)
			targetURL := baseURL + urlData[u]
			results := make(chan int, concurrency)

			for c := 0; c < concurrency; c++ {
				wg.Add(1)
				go sendRequestNew(targetURL, tokenData[c], wg, results)
			}

			// === Wait for all goroutines to finish
			wg.Wait()
			close(results)

			// === Collect and analyze results
			totalTime := 0
			successfulRequests := 0
			for elapsed := range results {
				if elapsed > 0 {
					totalTime += elapsed
					successfulRequests++
				}
			}
			averageResponseTime := float64(totalTime) / float64(successfulRequests)
			failures := concurrency - successfulRequests

			// === Print results
			fmt.Println("==================")
			fmt.Printf("Endpoint: %v\n", targetURL)
			fmt.Printf("Successful Requests: %d\n", successfulRequests)
			fmt.Printf("Failed Requests: %d\n", failures)
			fmt.Printf("Total Requests: %d\n", concurrency)
			fmt.Printf("Batch Request: %d\n", n+1)
			fmt.Printf("Total Batch Requests: %d\n", numRequests)
			fmt.Printf("Average Response Time (ms): %.2f\n", averageResponseTime)
			fmt.Printf("\n\n\n")

			// export to csv
			exportToCSV(timestamp, targetURL, successfulRequests, failures, concurrency, n+1, numRequests, averageResponseTime)

			// === Sleep for 2 second while cleaning up
			time.Sleep(1 * time.Second)
			results = nil
			time.Sleep(1 * time.Second)
			wg = nil
		}
	}
}

func exportToCSV(timestamp string, targetURL string, successfulRequests int, failures int, concurrency int, batchRequest int, totalBatchRequests int, averageResponseTime float64) {
	// check if file exists
	if _, err := os.Stat("result.csv"); os.IsNotExist(err) {
		// create file
		f, err := os.Create("result.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// write header
		if _, err := f.WriteString("Timestamp;Endpoint;Successful Requests;Failed Requests;Total Requests;Batch Request;Total Batch Requests;Average Response Time (ms)\n"); err != nil {
			log.Fatal(err)
		}
	}

	// append to file
	f, err := os.OpenFile("result.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%.2f\n", timestamp, targetURL, successfulRequests, failures, concurrency, batchRequest, totalBatchRequests, averageResponseTime)); err != nil {
		log.Fatal(err)
	}
}

func sendRequestNew(url string, token string, wg *sync.WaitGroup, results chan<- int) {
	// === Close goroutine when function is finished
	defer wg.Done()

	// === Set new request
	startTime := time.Now()
	client := http.Client{}
	client.Timeout = time.Second * time.Duration(config.Config.RequestTimeout)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		results <- 0
		return
	}

	// === Set request header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-CF-Secret", "DCVoys8d9xxjYoLsvmhFLKi6k8l2QmIA")
	req.Header.Set("'Platform", "android")

	// === Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		results <- 0
		return
	}
	defer resp.Body.Close()

	// === Collect results
	elapsed := time.Since(startTime).Milliseconds()
	fmt.Println(resp.StatusCode, token[280:290], elapsed)
	if resp.StatusCode != 200 {
		results <- 0
	} else {
		results <- int(elapsed)
	}
}

func readTokenFromTXT() ([]string, int) {
	var data []string
	file, err := os.Open(config.Config.TokenListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := len(data)
	if count == 0 {
		log.Fatal("Err: Token list is empty")
	}

	return data, count
}

func readURLFromTXT() ([]string, int) {
	var data []string
	file, err := os.Open(config.Config.UrlListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	count := len(data)
	if count == 0 {
		log.Fatal("Err: Url list is empty")
	}

	return data, count
}
