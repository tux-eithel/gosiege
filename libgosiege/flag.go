package libgosiege

import (
	"errors"
	"fmt"
	"strings"
)

// Type FlagUrl re-define Requests struct to implement flag.Value interface
// so you can use a FlagUrl struct to parse directly cli param
type FlagUrl struct {
	Req *Requests
}

// Init initialize the struct
func (fu *FlagUrl) Init() {
	fu.Req = NewRequests()
}

// String return all url parsed in one single string
func (fu *FlagUrl) String() string {

	var srt []string
	for i := 0; i < len(fu.Req.Reqs); i++ {
		srt = append(srt, fu.Req.Reqs[i].Url)
	}
	return strings.Join(srt, ",")

}

// When it parse an url, it will be added to a Requests struct.
// If url gives error, url will be skipped
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

type FlagRegexp struct {
	Rexp *CompareHeader
}

func (fr *FlagRegexp) Init() {
	fr.Rexp = NewCompareHeader()
}

func (fr *FlagRegexp) String() string {
	return fr.Rexp.String()
}

func (fr *FlagRegexp) Set(srt string) error {

	app := strings.Split(srt, " ")
	if len(app) != 2 {
		return errors.New("Indicate header key and regex separated from space")
	}

	fr.Rexp.Add(app[0], app[1])

	return nil

}
