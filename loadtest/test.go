package loadtest

import (
	"net/http"
	"sync"
	"time"
	"performance_testing/utils"
	"performance_testing/config"
)

type Result struct {
	Success       int
	Fail          int
	Duration      time.Duration
	AvgResp       float64
	CPUPercent    float64
	RAMUsedMB     float64
	GoroutinesNum int
}

func RunTest(cfg config.Config) Result {
	start := time.Now()
	var wg sync.WaitGroup
	var mu sync.Mutex

	success, fail := 0, 0
	times := []time.Duration{}

	sem := make(chan struct{}, cfg.Concurrency)

	for i := 0; i < cfg.TotalReq; i++ {
		for _, url := range cfg.URLs {
			wg.Add(1)
			sem <- struct{}{}

			go func(u string) {
				defer wg.Done()
				defer func() { <-sem }()

				reqStart := time.Now()
				resp, err := http.Get(u)
				duration := time.Since(reqStart)

				mu.Lock()
				times = append(times, duration)
				if err != nil || resp.StatusCode >= 400 {
					fail++
				} else {
					success++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}(url)
		}
	}

	wg.Wait()
	totalDuration := time.Since(start)

	var total float64
	for _, d := range times {
		total += float64(d.Milliseconds())
	}
	avg := total / float64(len(times))

	return Result{
		Success:       success,
		Fail:          fail,
		Duration:      totalDuration,
		AvgResp:       avg,
		CPUPercent:    utils.GetCPUUsage(),
		RAMUsedMB:     utils.GetRAMUsage(),
		GoroutinesNum: utils.GetGoroutineCount(),
	}
}
