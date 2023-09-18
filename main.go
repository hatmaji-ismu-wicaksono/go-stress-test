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

		wg sync.WaitGroup
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
			targetURL := baseURL + urlData[u]
			results := make(chan int, numRequests)
			for c := 0; c < concurrency; c++ {
				// fmt.Println("\n", "Batch Request: ", n, "; Virtual User: ", c, "; Target URL: ", targetURL, "; Token:", tokenData[c])
				wg.Add(1)
				go sendRequestNew(targetURL, tokenData[c], &wg, results)
			}
			// Wait for all goroutines to finish
			wg.Wait()
			close(results)

			// Collect and analyze results
			totalTime := 0
			successfulRequests := 0

			for elapsed := range results {
				if elapsed > 0 {
					totalTime += elapsed
					successfulRequests++
				}
			}

			averageResponseTime := float64(totalTime) / float64(successfulRequests)
			fmt.Printf("\n")
			fmt.Printf("\n==========> Batch Request: %d\n", n+1)
			fmt.Printf("Endpoint: %v\n", targetURL)
			fmt.Printf("Total Batch Requests: %d\n", numRequests)
			fmt.Printf("Total Requests: %d\n", concurrency)
			fmt.Printf("Successful Requests: %d\n", successfulRequests)
			fmt.Printf("Average Response Time (ms): %.2f\n", averageResponseTime)
		}
	}
}

func sendRequestNew(url string, token string, wg *sync.WaitGroup, results chan<- int) {
	startTime := time.Now()
	client := &http.Client{}
	client.Timeout = time.Second * 1
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		results <- 0
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", "k1DzR0yYWIyZgvTvixiDHnb4Nl08wSU0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		results <- 0
		return
	}
	resp.Body.Close()
	elapsed := time.Since(startTime).Milliseconds()
	results <- int(elapsed)
	wg.Done()
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

	return data, count
}
