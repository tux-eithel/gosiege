package libgosiege

import (
	"fmt"
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

		appR, err := NewInputRequest(value)
		if err != nil {
			fmt.Println("Error making Request for '", value, "' with error:", err)
			fmt.Println("Url ignored")
		} else {
			fu.Req.AddRequest(appR)
		}

	}
	return nil
}
