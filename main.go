package main

import (
	"performance_testing/config"
	"performance_testing/loadtest"
)

func main() {
	cfg := config.LoadConfig()
	loadtest.RunTest(cfg) // This now prints all per-request logs
}
