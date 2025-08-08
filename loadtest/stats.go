package loadtest

import "fmt"

func PrintSummary(r Result) {
	fmt.Printf("Total Duration     : %v\n", r.Duration)
	fmt.Printf("Average Response   : %.2f ms\n", r.AvgResp)
	fmt.Printf("Successful Requests: %d\n", r.Success)
	fmt.Printf("Failed Requests    : %d\n", r.Fail)
	fmt.Printf("CPU Usage          : %.2f%%\n", r.CPUPercent)
	fmt.Printf("RAM Usage          : %.2f MB\n", r.RAMUsedMB)
	fmt.Printf("Goroutines Used    : %d\n", r.GoroutinesNum)
}
