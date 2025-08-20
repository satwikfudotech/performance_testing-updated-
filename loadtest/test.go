package loadtest

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"performance_testing/config"
	"performance_testing/utils"
)

func RunTest(cfg config.Config) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	times := []time.Duration{}
	reqID := 0
	sem := make(chan struct{}, cfg.Concurrency)

	totalCount := 0

	// 🔹 create CSV file
	file, err := os.Create("results.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 🔹 write header row
	writer.Write([]string{
		"Timestamp", "RequestID", "URL", "Latency(ms)", "Status", "HTTPCode",
		"CPU(%)", "RAM(MB)", "AverageResponse(ms)",
	})

	for totalCount < cfg.TotalReq {
		for _, url := range cfg.URLs {
			if totalCount >= cfg.TotalReq {
				break
			}

			reqID++
			totalCount++
			wg.Add(1)
			sem <- struct{}{}

			go func(id int, u string) {
				defer wg.Done()
				defer func() { <-sem }()

				start := time.Now()
				resp, err := http.Get(u)
				latency := time.Since(start)

				status := "✅"
				code := 0
				if resp != nil {
					code = resp.StatusCode
					resp.Body.Close()
				}
				if err != nil || code >= 400 {
					status = "❌"
				}

				mu.Lock()
				times = append(times, latency)

				// running average
				var total float64
				for _, t := range times {
					total += float64(t.Milliseconds())
				}
				avgResp := total / float64(len(times))

				// system stats
				cpu := utils.GetCPUUsage()
				ram := utils.GetRAMUsage()

				// timestamp
				timestamp := time.Now().Format("2006-01-02 15:04:05")

				// 🔹 log to console
				fmt.Printf("[%s] Request %d %s  %v ms  %s %d  CPU: %.2f%%  RAM: %.2f MB | Average Response: %.2f ms\n",
					timestamp, id, u, latency.Milliseconds(), status, code, cpu, ram, avgResp,
				)

				// 🔹 write to CSV
				writer.Write([]string{
					timestamp,
					fmt.Sprintf("%d", id),
					u,
					fmt.Sprintf("%d", latency.Milliseconds()),
					status,
					fmt.Sprintf("%d", code),
					fmt.Sprintf("%.2f", cpu),
					fmt.Sprintf("%.2f", ram),
					fmt.Sprintf("%.2f", avgResp),
				})
				writer.Flush()

				mu.Unlock()

			}(reqID, url)
		}
	}

	wg.Wait()
}
