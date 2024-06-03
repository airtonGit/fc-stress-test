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
	responsesC := make(chan *http.Response)
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

//func (r *requestWorker) ResponsesConsumer(doneC chan struct{}, responsesC chan *http.Response) {
//	for response := range responsesC {
//		r.responses = append(r.responses, response)
//	}
//	doneC <- struct{}{}
//}

//func (r *requestWorker) ResultReport(begin time.Time) {
//	elapsed := time.Since(begin)
//	fmt.Println("Elapsed time", elapsed.Seconds(), "seconds")
//
//	summary := make(map[int]int)
//	for _, response := range r.responses {
//		statusCode := response.StatusCode
//		if _, ok := summary[statusCode]; !ok {
//			summary[statusCode] = 0
//		}
//		summary[statusCode]++
//	}
//
//	fmt.Println(fmt.Sprintf("=> Quantidade de request status 200: %d", summary[200]))
//	delete(summary, 200)
//	fmt.Println("=> Outros status codes qtd", len(summary))
//
//	for code, count := range summary {
//		fmt.Println(fmt.Sprintf("=> Quantidade de request %d: %d", code, count))
//	}
//}

func (r *requestWorker) DoRequest(ctx context.Context, wg *sync.WaitGroup, requestsC chan struct{}, responsesC chan *http.Response) {
	defer wg.Done()
	for range requestsC {
		if ctx.Err() != nil {
			fmt.Print("context error", ctx.Err())
			return
		}
		resp, err := r.httpClient.Do(r.req)
		if err != nil {
			fmt.Print("error", err)
		}
		responsesC <- resp
	}
}
