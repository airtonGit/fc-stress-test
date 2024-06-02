package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/airtongit/fc-stress-test/internal"
	"github.com/spf13/cobra"
)

func main() {
	var (
		targetURL             string
		requests, concurrency int
	)

	rootCmd := &cobra.Command{
		Use:   "fc-stress-test",
		Short: "A simple stress test for Firecracker",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s, %d requests, %d concurrency!\n", targetURL, requests, concurrency)

			if !strings.HasPrefix(targetURL, "http://") {
				targetURL = "http://" + targetURL
			}

			requestURL, err := url.Parse(targetURL)
			if err != nil {
				fmt.Println(err)
				return
			}

			req, err := http.NewRequest("GET", requestURL.String(), nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			req.Header.Set("Content-Type", "application/json")

			reqWorker := internal.NewRequestWorker(http.DefaultClient, req)

			stressTester := internal.NewStressTester(reqWorker, requests, concurrency)
			stressTester.Run(cmd.Context())
		},
	}

	rootCmd.Flags().StringVarP(&targetURL, "targetURL", "u", "", "Target service to stress test")
	rootCmd.MarkFlagRequired("targetURL")
	rootCmd.Flags().IntVarP(&requests, "requests", "r", 1, "Total number of requests to send")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests to send")
	rootCmd.MarkFlagRequired("concurrency")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
