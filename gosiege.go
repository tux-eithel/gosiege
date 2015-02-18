package main

import (
	"flag"
	"fmt"
	_ "net/url"
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

func init() {

	listUrls.Init()
	flag.IntVar(&numberConcurrent, "n", 1, "Number of concurrent request")
	flag.DurationVar(&secToWait, "s", time.Duration(1)*time.Second, "Time to wait until next request")
	flag.Var(&listUrls, "u", "Url(s) to test")
	flag.StringVar(&inputFile, "f", "", "Input file with urls")
}

func main() {
	flag.Parse()
	//	fmt.Println(numberConcurrent)
	//	fmt.Println(secToWait)
	//	fmt.Println(listUrls.String())

	for i := 0; i < len(listUrls.Req.Reqs); i++ {
		ToRun(listUrls.Req)
	}
}

func ToRun(totest *parseUrl.Requests) error {

	//	for {
	req := totest.NextUri(false)
	if req == nil {
		fmt.Println("Qualcosa di sbagliato")
	}

	fmt.Printf("%#v\n", req)
	// do Stuff
	time.Sleep(secToWait)
	//	}

	return nil
}
