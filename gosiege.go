package main

import (
	"flag"
	"fmt"
	_ "log"
	_ "net/url"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"bitbucket.org/tux-eithel/gosiege/parseUrl"
)

// -n param
var numberConcurrent int

// -s param
var secToWait time.Duration

// -u param
var listUrls parseUrl.FlagUrl

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

	//	fmt.Println(numberConcurrent)
	//	fmt.Println(secToWait)
	//	fmt.Println(listUrls.String())

	quitChannel := make(chan os.Signal)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)

	shutdownChannel := make(chan bool)

	nextChannel := make(chan int, 1)

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)

	for i := 0; i < numberConcurrent; i++ {
		go ToRun(listUrls.Req, randomUrl, nextChannel, shutdownChannel, waitGroup)
	}

	if !randomUrl {
		nextChannel <- 1
	}

	<-quitChannel
	shutdownChannel <- true
	fmt.Println("Received quit. Sending shutdown and waiting on goroutines...")

	/*
	 * Block until wait group counter gets to zero
	 */
	waitGroup.Wait()
	fmt.Println("Done.")
}

func ToRun(totest *parseUrl.Requests, randomUrl bool, nextChannel chan int, shutdownChannel chan bool, waitGroup *sync.WaitGroup) error {

	defer waitGroup.Done()
	for {

		select {
		case _ = <-shutdownChannel:
			return nil

		default:
		}

		if !randomUrl {
			<-nextChannel
		}
		req := totest.NextUri(randomUrl)
		if !randomUrl {
			nextChannel <- 1
		}

		if req == nil {
			fmt.Println("Qualcosa di sbagliato")
		}

		fmt.Printf("%#v\n", req.Url)
		// do Stuff
		time.Sleep(secToWait)

	}

	return nil
}
