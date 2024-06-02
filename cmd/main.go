package main

import (
	"fmt"
	"github.com/airtongit/fc-stress-test/internal"
	"net/http"

	"github.com/spf13/cobra"
)

func main() {
	var (
		url                   string
		requests, concurrency int
	)

	rootCmd := &cobra.Command{
		Use:   "fc-stress-test",
		Short: "A simple stress test for Firecracker",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s, %d requests, %d concurrency!\n", url, requests, concurrency)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			reqWorker := internal.NewRequestWorker(http.DefaultClient, req)

			stressTester := internal.NewStressTester(reqWorker, requests, concurrency)
			stressTester.Run(cmd.Context())
			reqWorker.ResultReport()
		},
	}

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "Target service to stress test")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Total number of requests to send")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests to send")
	rootCmd.MarkFlagRequired("concurrency")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
