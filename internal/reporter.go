package internal

import (
	"fmt"
	"time"
)

type resultsReporter struct {
	responses     []RequestResult
	consumerDoneC chan struct{}
}

func NewResultsReporter() *resultsReporter {
	return &resultsReporter{
		responses:     make([]RequestResult, 0),
		consumerDoneC: make(chan struct{}),
	}
}

func (r *resultsReporter) ConsumeResponses(responsesC chan RequestResult) {
	for requestResult := range responsesC {
		r.responses = append(r.responses, requestResult)
	}
	r.consumerDoneC <- struct{}{}
}

func (r *resultsReporter) ResultReport(begin time.Time) {
	<-r.consumerDoneC
	elapsed := time.Since(begin)
	fmt.Println("Elapsed time", elapsed.Seconds(), "seconds")

	summary := make(map[string]int)
	for _, requestResult := range r.responses {
		if requestResult.Error != nil {
			summary["error"]++
			continue
		}
		statusCode := fmt.Sprintf("%d", requestResult.Response.StatusCode)
		if _, ok := summary[statusCode]; !ok {
			summary[statusCode] = 0
		}
		summary[statusCode]++
		requestResult.Response.Body.Close()
	}

	fmt.Println(fmt.Sprintf("=> status 200: %d", summary["200"]))
	delete(summary, "200")
	for code, count := range summary {
		fmt.Println(fmt.Sprintf("=> status %s: %d", code, count))
	}
}
