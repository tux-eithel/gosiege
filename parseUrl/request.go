package parseUrl

import (
	"math/rand"
	"time"
)

type InputRequest struct {
	Method string
	Url    string
	Header map[string]string
	Body   string
	Param  map[string]string
}

func NewInputRequest(url string) *InputRequest {
	return &InputRequest{
		Method: "GET",
		Url:    url,
	}
}

type Requests struct {
	Reqs []*InputRequest
	Rand *rand.Rand
	Cont int
}

func NewRequests() *Requests {
	return &Requests{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		Cont: 0,
	}
}

func (r *Requests) AddRequest(ir *InputRequest) {
	r.Reqs = append(r.Reqs, ir)
}

func (r *Requests) NextUri(isRandom bool) *InputRequest {
	if isRandom {
		return r.Reqs[r.Rand.Intn(len(r.Reqs))]
	}
	next := r.Cont + 1
	if next == len(r.Reqs) {
		next = 0
	} else {
		r.Cont++
	}
	return r.Reqs[next]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
