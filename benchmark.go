package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hgtpcastro/go-expert-lab-stress-test/benchmark/utils"
)

// StressParameters stress params for worker
type StressParameters struct {
	Requests    int    `json:"n"`   // N is the total number of requests to make.
	Concurrency int    `json:"c"`   // C is the concurrency level, the number of concurrent workers to run.
	Url         string `json:"url"` // Request url.
}

func (p *StressParameters) String() string {
	body, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err.Error()
	}
	return string(body)
}

type RequestResult struct {
	Err        error
	StatusCode int
}

var (
	urlstr       = flag.String("url", "", "")
	n            = flag.Int("requests", 0, "")     // Number of requests to run
	c            = flag.Int("concurrency", 50, "") // Number of requests to run concurrently
	printExample = flag.Bool("example", false, "")
)

const (
	usage = `Usage: go-http-bench [options...] <url>
Options:
	-n		Number of requests to run.
	-c		Number of requests to run concurrently.	
	-url		Request single url.
	-example	Print some stress test examples (default false).`

	examples = `
1.Example stress test:
	./bin/go-http-bench -url "http://google.com" -requests 10 -concurrency 2 
	./bin/go-http-bench "http://google.com" -requests 10 -concurrency 2`
)

func main() {
	flag.Usage = func() {
		fmt.Println(usage)
	}

	var params StressParameters

	flag.Parse()

	for flag.NArg() > 0 {
		if len(*urlstr) == 0 {
			*urlstr = flag.Args()[0]
		}
		os.Args = flag.Args()[0:]
		flag.Parse()
	}

	if *printExample {
		println(examples)
		return
	}

	params.Requests = *n
	params.Concurrency = *c
	params.Url = *urlstr

	if params.Requests <= 0 {
		utils.UsageAndExitt("n and c cannot be smaller than 1.")
	}

	if params.Requests < params.Concurrency {
		utils.UsageAndExitt("n cannot be less than c.")
	}

	startTime := time.Now()
	requestResult := make(chan *RequestResult)

	go executeStress(params, requestResult)

	stressResult := NewStressResult()
	var totalRequests int

	for i := 0; i < params.Requests; i++ {
		result := <-requestResult
		totalRequests++
		stressResult.append(result)
	}
	stressResult.Duration = time.Since(startTime)
	stressResult.TotalRequests = totalRequests

	stressResult.print()

	// fmt.Println(params.String())

}

func executeStress(params StressParameters, requestResult chan *RequestResult) {
	requestsPerWorker := params.Requests / params.Concurrency
	httpClient := &http.Client{}

	for i := 0; i < params.Concurrency; i++ {
		go func() {
			for j := 0; j < requestsPerWorker; j++ {
				resp, err := httpClient.Get(params.Url)
				if err != nil {
					requestResult <- &RequestResult{StatusCode: -1, Err: err}
				} else {
					requestResult <- &RequestResult{StatusCode: resp.StatusCode, Err: err}
				}
			}
		}()
	}

	remainder := params.Requests % params.Concurrency
	for j := 0; j < remainder; j++ {
		resp, err := httpClient.Get(params.Url)
		requestResult <- &RequestResult{StatusCode: resp.StatusCode, Err: err}
	}
}

func NewStressResult() *StressResult {
	return &StressResult{
		Duration:       0,
		TotalRequests:  0,
		ErrorDist:      make(map[string]int, 0),
		StatusCodeDist: make(map[int]int, 0),
	}
}

var resultRdMutex sync.RWMutex

type StressResult struct {
	Duration      time.Duration
	TotalRequests int
	//StatusCode1xx, StatusCode2xx, StatusCode3xx, StatusCode4xx, StatusCode5xx []int
	ErrorDist      map[string]int
	StatusCodeDist map[int]int
}

func (result *StressResult) print() {
	println("Summary:")
	println(fmt.Sprintf("Total time: %s", result.Duration))
	println(fmt.Sprintf("Total requests: %d", result.TotalRequests))
	result.printStatusCodes()
	if len(result.ErrorDist) > 0 {
		result.printErrors()
	}
}

func (result *StressResult) printStatusCodes() {
	println("\nStatus code distribution:")
	for code, num := range result.StatusCodeDist {
		println(fmt.Sprintf("  [%d]\t%d responses", code, num))
	}
}

func (result *StressResult) printErrors() {
	println("\nError distribution:")
	for err, num := range result.ErrorDist {
		println(fmt.Printf("  [%d]\t%s", num, err))
	}
}

func (result *StressResult) append(res *RequestResult) {
	resultRdMutex.Lock()
	defer resultRdMutex.Unlock()

	if res.Err != nil {
		result.ErrorDist[res.Err.Error()]++
	} else {
		result.StatusCodeDist[res.StatusCode]++
	}
}
