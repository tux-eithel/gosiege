package main

import (
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

// -n param
var numberConcurrent int

// -s param
var secToWait time.Duration

// -u param
var listUrls libgosiege.FlagUrl

// -f param
var inputFile string

// -nasty param
var isNasty bool

// r param
var randomUrl bool

func init() {

	listUrls.Init()
	flag.IntVar(&numberConcurrent, "n", 1, "Number of concurrent request")
	flag.DurationVar(&secToWait, "s", time.Duration(1)*time.Second, "Time to wait until next request")
	flag.Var(&listUrls, "u", "Url(s) to test")
	flag.StringVar(&inputFile, "f", "", "Input file with urls")
	flag.BoolVar(&isNasty, "nasty", false, "Use all available CPU cores")
	flag.BoolVar(&randomUrl, "r", true, "Use random urls from list")
}

func main() {
	flag.Parse()

	if isNasty {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	waitData := &sync.WaitGroup{}
	waitData.Add(1)
	dataChannel := make(chan *libgosiege.SimpleCounter, numberConcurrent*2)

	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdownChannel := make(chan bool, numberConcurrent)
	shutdownProcessData := make(chan bool)

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(numberConcurrent)

	go libgosiege.ProcessData(dataChannel, shutdownProcessData, waitData)

	for i := 0; i < numberConcurrent; i++ {
		go ToRun(listUrls.Req, dataChannel, randomUrl, shutdownChannel, waitGroup)
	}

	<-quitChannel
	fmt.Println("Received quit. Sending shutdown and waiting on goroutines...")
	for i := 0; i < numberConcurrent; i++ {
		shutdownChannel <- true
	}

	// Block until wait group counter gets to zero
	waitGroup.Wait()

	shutdownProcessData <- true
	waitData.Wait()
	fmt.Println("Done.")
}

func ToRun(totest *libgosiege.Requests, dataChannel chan *libgosiege.SimpleCounter, randomUrl bool, shutdownChannel chan bool, waitGroup *sync.WaitGroup) error {

	var t0 time.Time
	var diff time.Duration
	defer waitGroup.Done()
	for {

		select {
		case _ = <-shutdownChannel:
			return nil

		default:
		}

		req := totest.NextUri(randomUrl)

		if req == nil {

			fmt.Println("Seems strange, no Url recover")

		} else {

			t0 = time.Now()
			r, err := http.DefaultClient.Do(req.ReadyUrl)
			diff = time.Since(t0)

			if err != nil {

				fmt.Printf("Response Error: %v | Response Object:  %+v\n", err, r)
				// return body, statusCode, response_headers, err

			} else {

				defer r.Body.Close()

				body, err := ioutil.ReadAll(r.Body)
				qtaBody := -1
				if err == nil {
					qtaBody = len(body)
				}

				// TODO: here we'll put goroutine to manage result data

				dataChannel <- libgosiege.NewSimpleCounter(float64(qtaBody), diff.Seconds(), r.StatusCode, req.ReadyUrl.URL.Path)

			}

		}

		// do Stuff
		time.Sleep(secToWait)

	}

	return nil
}
