// Package libgosiege package define all the structs and functions necessary for Gosiege
// Inside this package there is some new flags for better parse CLI params
// There is a struct for save Urls to parse and retrive in different way
// Inside results.go different structs works together for handle all the statistics
package libgosiege

import (
	"errors"
	"fmt"
	"strings"
)

// FlagUrl re-define Requests struct to implement flag.Value interface
// so you can use a FlagUrl struct to parse directly cli param
type FlagUrl struct {
	Req *Requests
}

// Init initializes the FlagUrl struct
func (fu *FlagUrl) Init() {
	fu.Req = NewRequests()
}

// String returns all url parsed in one single string
func (fu *FlagUrl) String() string {

	var srt []string
	for i := 0; i < len(fu.Req.Reqs); i++ {
		srt = append(srt, fu.Req.Reqs[i].Url)
	}
	return strings.Join(srt, ",")

}

// Set parse urls and puts into Requests struct.
// If url gives error, url will be skipped
// Multiple urls may be defined
func (fu *FlagUrl) Set(srt string) error {

	app := strings.Split(srt, " ")
	for _, value := range app {

		appR, err := NewInputRequest(value)
		if err != nil {
			fmt.Println("Url ignored '", value, "' with error:", err)
		} else {
			fu.Req.AddRequest(appR)
		}

	}
	return nil

}

// FlagRegexp allows to use struct CompareHeader to create new "bucket"
// Every "bucket" is made from a value and a regexp
type FlagRegexp struct {
	Rexp *CompareHeader
}

// Init initializes the FlagRegexp struct
func (fr *FlagRegexp) Init() {
	fr.Rexp = NewCompareHeader()
}

// String returns all value-regexp in one single string
func (fr *FlagRegexp) String() string {
	return fr.Rexp.String()
}

// Set parse FlagRegexp flag.
// A flag is defined as "value regexp" and may be used multiple times
func (fr *FlagRegexp) Set(srt string) error {

	app := strings.Split(srt, " ")
	if len(app) != 2 {
		return errors.New("indicate header key and regex separated from space")
	}

	fr.Rexp.Add(app[0], app[1])

	return nil

}
