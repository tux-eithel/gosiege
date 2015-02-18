package parseUrl

import (
	"fmt"
	"net/url"
	"strings"
)

type FlagUrl struct {
	Req *Requests
}

func (fu *FlagUrl) Init() {
	fu.Req = NewRequests()
}

func (fu *FlagUrl) String() string {
	var srt []string
	for i := 0; i < len(fu.Req.Reqs); i++ {
		srt = append(srt, fu.Req.Reqs[i].Url)
	}
	return strings.Join(srt, ",")
}

func (fu *FlagUrl) Set(srt string) error {
	fu.Req = NewRequests()
	app := strings.Split(srt, " ")
	for _, value := range app {
		_, err := url.Parse(value)
		if err != nil {
			fmt.Println("Url: ", value, "not correct, skipped")
		}
		fu.Req.AddRequest(NewInputRequest(value))
	}
	return nil
}
