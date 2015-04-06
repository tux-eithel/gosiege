package libgosiege

import (
	"math/rand"
	"net/url"
	"sync"
	"time"
)

// InputRequest represents a parsed url from cli or file
type InputRequest struct {
	Method string
	Url    string
	Header map[string]string
	Body   []byte
	Param  map[string]string
}

// NewInputRequest creates a new InputRequest from a url string.
// The new InputRequest is a basic GET request.
// The function also modify the input url adding Host and Scheme if not specified.
// For example google.com becames http://google.com
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

	return in, nil
}

// Requests keeps all the InputRequest in one array.
// Its structure also contains the current-1 object to extract in case of sequential reading,
// or an object for random number generation in case of random reading
type Requests struct {
	Reqs []*InputRequest
	Rand *rand.Rand
	Cont int
	sync.Mutex
}

// NewRequests inizialize a Requests object
func NewRequests() *Requests {
	appR := &Requests{
		Rand: rand.New(rand.NewSource(time.Now().UnixNano())),
		Cont: -1,
	}
	return appR
}

// AddRequest add to a Requests an InputRequest
// Don't use NewInputRequest directly as an argument. NewInputRequest may return errors!
func (r *Requests) AddRequest(ir *InputRequest) {
	r.Reqs = append(r.Reqs, ir)
}

// NextUri return the next uri to be processed. If isRandom is true, a random url from input is returned, else the next one in order.
// Pass isRandom as false value is useful if you want to process all the input urls
// Now NextUri is thread-safe when isRandom is false. Well is thread-safe the sum of index for next element, but not the result of the function!
// For that, you have to make thread-safe this function when you are going to call it. But if you really want this it's better run gosiege with -n 1 param
func (r *Requests) NextUri(isRandom bool) *InputRequest {

	if len(r.Reqs) == 0 {
		return nil
	}

	if len(r.Reqs) == 1 {
		return r.Reqs[0]
	}

	r.Lock()
	defer r.Unlock()

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

// init set the Seed for random number
func init() {
	rand.Seed(time.Now().UnixNano())
}
