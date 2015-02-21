package parseUrl

import (
	"math/rand"
	"time"
)

type InputRequest struct {
	Method   string
	Url      string
	Header   map[string]string
	Body     []byte
	Param    map[string]string
	ReadyUrl *http.Request
}

func NewInputRequest(inputUrl string) (*InputRequest, error) {

	u, err := url.Parse(inputUrl)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		inputUrl = "//" + inputUrl
	}
	if u.Scheme == "" {
		inputUrl = "http:" + inputUrl
	}

	in := &InputRequest{
		Method: "GET",
		Url:    inputUrl,
	}

	req, err := http.NewRequest(in.Method, in.Url, bytes.NewBuffer(in.Body))
	if err != nil {
		return nil, err
	}

	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}

	for key, value := range in.Header {
		req.Header.Set(key, value)
	}

	in.ReadyUrl = req

	return in, nil
}

type Requests struct {
	Reqs []*InputRequest
	Rand *rand.Rand
	Cont int
}

func NewRequests() *Requests {
	return &Requests{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		Cont: -1,
	}
}

func (r *Requests) AddRequest(ir *InputRequest) {
	r.Reqs = append(r.Reqs, ir)
}

func (r *Requests) NextUri(isRandom bool) *InputRequest {
	if isRandom {
		return r.Reqs[r.Rand.Intn(len(r.Reqs))]
	}

	if r.Cont+1 == len(r.Reqs) {
		r.Cont = 0
	} else {
		r.Cont++
	}
	return r.Reqs[r.Cont]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
