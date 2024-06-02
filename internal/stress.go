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
	wg := new(sync.WaitGroup)
	wg.Add(s.Concurrency)
	for range s.Concurrency {
		go s.ReqWorker.DoRequest(ctx, wg, requestsC)
	}
	fmt.Println("Stating requests ", s.TotalRequests)
	begin := time.Now()
	for range s.TotalRequests {
		requestsC <- struct{}{}
	}
	close(requestsC)
	wg.Wait()
	s.ReqWorker.ResultReport(begin)
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type requestWorker struct {
	httpClient HttpClient
	req        *http.Request
	respC      chan *http.Response
	responses  []*http.Response
}

func (r *requestWorker) responsesConsumer() {
	for response := range r.respC {
		r.responses = append(r.responses, response)
	}
}

func NewRequestWorker(client HttpClient, req *http.Request) *requestWorker {
	rw := &requestWorker{
		httpClient: client,
		req:        req,
		respC:      make(chan *http.Response),
	}
	go rw.responsesConsumer()
	return rw
}

func (r *requestWorker) ResultReport(begin time.Time) {
	close(r.respC)
	elaspsed := time.Since(begin)
	summary := make(map[int]int)

	for _, response := range r.responses {
		statusCode := response.StatusCode
		if _, ok := summary[statusCode]; !ok {
			summary[statusCode] = 0
		}
		summary[statusCode]++
	}
	fmt.Println(fmt.Sprintf("Quantidade de request status 200: %d", summary[200]))
	delete(summary, 200)
	fmt.Println("Outros status codes qtd", len(summary))
	for code, count := range summary {
		fmt.Println(fmt.Sprintf("Quantidade de request %d: %d", code, count))
	}
	fmt.Println("Elapsed time", elaspsed.Seconds(), "seconds")
}

func (r *requestWorker) DoRequest(ctx context.Context, wg *sync.WaitGroup, requestsC chan struct{}) {
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
		r.respC <- resp
	}
}
