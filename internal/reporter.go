package internal

import (
	"fmt"
	"net/http"
	"time"
)

type resultsReporter struct {
	responses     []*http.Response
	consumerDoneC chan struct{}
}

func NewResultsReporter() *resultsReporter {
	return &resultsReporter{
		responses:     make([]*http.Response, 0),
		consumerDoneC: make(chan struct{}),
	}
}

func (r *resultsReporter) ConsumeResponses(responsesC chan *http.Response) {
	for response := range responsesC {
		r.responses = append(r.responses, response)
	}
	r.consumerDoneC <- struct{}{}
}

func (r *resultsReporter) ResultReport(begin time.Time) {
	<-r.consumerDoneC
	elapsed := time.Since(begin)
	fmt.Println("Elapsed time", elapsed.Seconds(), "seconds")

	summary := make(map[int]int)
	for _, response := range r.responses {
		statusCode := response.StatusCode
		if _, ok := summary[statusCode]; !ok {
			summary[statusCode] = 0
		}
		summary[statusCode]++
	}

	fmt.Println(fmt.Sprintf("=> status 200: %d", summary[200]))
	delete(summary, 200)
	for code, count := range summary {
		fmt.Println(fmt.Sprintf("=> status %d: %d", code, count))
	}
}
