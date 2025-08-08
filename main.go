package main

import (
	"fmt"
	"performance_testing/config"
	"performance_testing/loadtest"
)

func main() {
	cfg := config.LoadConfig()
	results := loadtest.RunTest(cfg)

	fmt.Println("\n--- Load Test Summary ---")
	loadtest.PrintSummary(results)
}
