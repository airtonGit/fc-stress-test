package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var (
		name, requests string
		concurrency    int
	)

	rootCmd := &cobra.Command{
		Use:   "fc-stress-test",
		Short: "A simple stress test for Firecracker",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hello, %s!\n", name)
		},
	}

	rootCmd.Flags().StringVarP(&name, "url", "u", "", "Target service to stress test")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().StringVarP(&requests, "requests", "r", "", "Total number of requests to send")
	rootCmd.MarkFlagRequired("requests")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 1, "Number of concurrent requests to send")
	rootCmd.MarkFlagRequired("concurrency")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
