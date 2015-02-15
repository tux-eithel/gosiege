package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"bitbucket.org/tux-eithel/gosiege/parseUrl"
)

// -n param
var numberConcurrent int

// -s param
var secToWait time.Duration

// -u param
var listUrls parseUrl.Urls

func init() {
	flag.IntVar(&numberConcurrent, "n", 1, "Number of concurrent request")
	flag.DurationVar(&secToWait, "s", time.Duration(1)*time.Second, "Time to wait until next request")
	flag.Var(&listUrls, "u", "Url(s) to test")
}

func main() {
	flag.Parse()
	//	fmt.Println(numberConcurrent)
	//	fmt.Println(secToWait)
	//	fmt.Println(listUrls.String())

	for i := 0; i < len(listUrls); i++ {
		ToRun(&listUrls)
	}
}

func ToRun(totest parseUrl.RandomUri) error {

	var uri *url.URL
	var srt *string
	var err error

	//	for {
	srt, err = totest.GetRandomUri()
	if err != nil {
		return err
	}
	uri, err = url.Parse(*srt)
	fmt.Printf("%#v\n", uri)
	// do Stuff
	time.Sleep(secToWait)
	//	}

	return nil
}
