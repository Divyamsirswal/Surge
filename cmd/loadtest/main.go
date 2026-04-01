package main

import (
	"bytes"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	totalRequests := 10000
	concurrency := 100

	url := "http://localhost:8080/api/v1/buy"
	payload := []byte(`{"user_id": 999, "product_id": 1}`)

	var successCount int32
	var soldOutCount int32
	var errorCount int32

	customTransport := &http.Transport{
		MaxIdleConns:        concurrency,
		MaxIdleConnsPerHost: concurrency,
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: customTransport,
	}

	fmt.Printf("Initiating Load Test: %d requests with %d concurrent workers...\n", totalRequests, concurrency)
	startTime := time.Now()

	jobs := make(chan int, totalRequests)
	var wg sync.WaitGroup

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for range jobs {
				req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
				req.Header.Set("Content-Type", "application/json")

				resp, err := client.Do(req)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					continue
				}

				if resp.StatusCode == http.StatusOK {
					atomic.AddInt32(&successCount, 1)
				} else if resp.StatusCode == http.StatusConflict {
					atomic.AddInt32(&soldOutCount, 1)
				} else {
					atomic.AddInt32(&errorCount, 1)
				}
				resp.Body.Close()
			}
		}()
	}

	for j := 1; j <= totalRequests; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
	duration := time.Since(startTime)

	fmt.Println("=======================================")
	fmt.Println("LOAD TEST RESULTS (Worker Pool Pattern)")
	fmt.Println("=======================================")
	fmt.Printf("Time Taken					   : %v\n", duration)
	fmt.Printf("Successful Orders (Got Ticket) : %d\n", successCount)
	fmt.Printf("Sold Out Rejections            : %d\n", soldOutCount)
	fmt.Printf("Server Errors                  : %d\n", errorCount)
	fmt.Println("=======================================")

}
