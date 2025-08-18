package loadtest

import (
	"fmt"
	"net/http"
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

	// total requests counter
	totalCount := 0

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

				// log per request
				fmt.Printf("[%s] Request %d %s  %v ms  %s %d  CPU: %.2f%%  RAM: %.2f MB | Average Response: %.2f ms\n",
					time.Now().Format("2006-01-02 15:04:05"),
					id,
					u,
					latency.Milliseconds(),
					status,
					code,
					utils.GetCPUUsage(),
					utils.GetRAMUsage(),
					avgResp,
				)
				mu.Unlock()

			}(reqID, url)
		}
	}

	wg.Wait()
}
