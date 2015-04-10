package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	_ "log"
	"net/http"
	_ "net/url"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"bitbucket.org/tux-eithel/gosiege/libgosiege"
)

// -c param
var numberConcurrent int

// -s param
var secToWait time.Duration

// -u param
var listUrls libgosiege.FlagUrl

// -f param
var inputFile string

// -nasty param
var isNasty bool

// rand param
var randomUrl bool

// exp param
var listRegexp libgosiege.FlagRegexp

// pexp param
var printRegexp bool

// per param
var maxRequestForURL int

func init() {

	listUrls.Init()
	listRegexp.Init()

	flag.IntVar(&numberConcurrent, "c", 1, "Number of concurrent request")
	flag.DurationVar(&secToWait, "s", time.Duration(1)*time.Second, "Time to wait until next request")
	flag.Var(&listUrls, "u", "Url(s) to test")
	flag.StringVar(&inputFile, "f", "", "Input file with urls")
	flag.BoolVar(&isNasty, "nasty", true, "Use all available CPU cores")
	flag.BoolVar(&randomUrl, "rand", true, "Use random urls from list")
	flag.Var(&listRegexp, "exp", "Regular expression for filter response header. Ex. \"X-Cache HIT\"")
	flag.BoolVar(&printRegexp, "pexp", false, "Print result flag during execution")
	flag.IntVar(&maxRequestForURL, "per", -1, "Number of time url will be hit")
}

func main() {
	flag.Parse()

	if isNasty {
		fmt.Println("We are going to use", runtime.NumCPU(), "CPU")
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if inputFile != "" {
		libgosiege.ParseAllInputFile(inputFile, listUrls.Req)
	}

	// print regexp during execution
	listRegexp.Rexp.PrintRegexp = printRegexp

	// set number of totat hit for request
	listUrls.Req.MaxRequest = maxRequestForURL

	waitData := &sync.WaitGroup{}
	waitData.Add(1)
	dataChannel := make(chan *libgosiege.SimpleCounter, numberConcurrent*2)

	quitChannel := make(chan os.Signal, numberConcurrent)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdownChannel := make(chan bool)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(numberConcurrent)

	go libgosiege.ProcessData(dataChannel, listRegexp.Rexp, waitData)

	fmt.Println("Prepare ", numberConcurrent, " goroutines for the battle")
	for i := 0; i < numberConcurrent; i++ {
		go ToRun(listUrls.Req, dataChannel, randomUrl, secToWait, shutdownChannel, quitChannel, waitGroup)
	}

	<-quitChannel
	fmt.Println("Received quit. Sending shutdown and waiting all goroutines...")
	close(shutdownChannel)

	waitGroup.Wait()
	close(dataChannel)

	waitData.Wait()
	fmt.Println("Done.")
}

// ToRun is the function which do the dirty work :)
// It creates a new http.NewRequest and makes the request
func ToRun(
	totest *libgosiege.Requests,
	dataChannel chan *libgosiege.SimpleCounter,
	randomUrl bool,
	secToWait time.Duration,
	shutdownChannel chan bool,
	quitChannel chan os.Signal,
	waitGroup *sync.WaitGroup) {

	var t0 time.Time
	var diff time.Duration
	var r *http.Response
	var rq *http.Request
	var err error
	var body []byte

	defer waitGroup.Done()
	for {

		select {
		case _ = <-shutdownChannel:
			// fmt.Println("Routine closed")
			return

		default:

			req := totest.NextUri(randomUrl)

			if req == nil {
				quitChannel <- syscall.SIGQUIT
				return

			}

			// fmt.Println("URL: ", req.Url)

			rq, err = http.NewRequest(req.Method, req.Url, bytes.NewBufferString(req.Body))
			if err != nil {

				dataChannel <- &libgosiege.SimpleCounter{
					Error: errors.New("Error preparing url '" + req.Url + "': " + err.Error()),
				}

			} else {
				for key, value := range req.Header {
					rq.Header.Set(key, value)
				}
				fmt.Println(rq.Body)

				t0 = time.Now()
				r, err = http.DefaultClient.Do(rq)
				diff = time.Since(t0)

				if err != nil {

					dataChannel <- &libgosiege.SimpleCounter{
						Error: errors.New("Response Error: " + err.Error()),
					}

				} else {

					body, err = ioutil.ReadAll(r.Body)
					qtaBody := -1
					if err == nil {
						qtaBody = len(body)
					}
					r.Body.Close()

					dataChannel <- libgosiege.NewSimpleCounter(float64(qtaBody), diff.Seconds(), r.StatusCode, rq.URL.Path, r.Header)

				}
			}

			// wait the next call
			time.Sleep(secToWait)
		}

	}

	return
}
