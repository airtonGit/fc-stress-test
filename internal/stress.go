package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type stressTester struct {
	ReqWorker     *requestWorker
	TotalRequests int
	Concurrency   int
}

func NewStressTester(reqWorker *requestWorker, totalRequests int, concurrency int) *stressTester {
	return &stressTester{
		ReqWorker:     reqWorker,
		TotalRequests: totalRequests,
		Concurrency:   concurrency,
	}
}

func (s *stressTester) Run(ctx context.Context) {
	requestsC := make(chan struct{})
	responsesC := make(chan RequestResult)
	reporter := NewResultsReporter()
	go reporter.ConsumeResponses(responsesC)

	wg := new(sync.WaitGroup)
	wg.Add(s.Concurrency)
	for i := 1; i <= s.Concurrency; i++ {
		go s.ReqWorker.DoRequest(ctx, wg, requestsC, responsesC)
	}
	fmt.Println("Starting requests ", s.TotalRequests)
	requestsBeginTimestamp := time.Now()
	for i := 1; i <= s.TotalRequests; i++ {
		requestsC <- struct{}{}
	}
	close(requestsC)
	wg.Wait()
	close(responsesC)
	reporter.ResultReport(requestsBeginTimestamp)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type requestWorker struct {
	httpClient HttpClient
	req        *http.Request
	responses  []*http.Response
}

func NewRequestWorker(client HttpClient, req *http.Request) *requestWorker {
	rw := &requestWorker{
		httpClient: client,
		req:        req,
	}

	return rw
}

type RequestResult struct {
	Response *http.Response
	Error    error
}

func (r *requestWorker) DoRequest(ctx context.Context, wg *sync.WaitGroup, requestsC chan struct{}, responsesC chan RequestResult) {
	defer wg.Done()
	for range requestsC {
		if ctx.Err() != nil {
			fmt.Print("context error", ctx.Err())
			return
		}
		resp, err := r.httpClient.Do(r.req)
		if err != nil {
			fmt.Println("error ", err)
			responsesC <- RequestResult{
				Response: nil,
				Error:    err,
			}
			return
		}
		responsesC <- RequestResult{
			Response: resp,
			Error:    nil,
		}
	}
}
