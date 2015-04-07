package libgosiege

import (
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"
	_ "time"
)

const (
	_          = iota // ignore first value by assigning to blank identifier
	KB float64 = 1 << (10 * iota)
	MB
	GB
	TB
)

func ByteSize(b float64) string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

type GeneralCounter struct {
	NumRequest  int
	NumSuccess  int
	NumBadError int
	LongTrans   float64
	ShortTrans  float64
	TotalTime   float64
	TotalByte   float64
	TransTime   []float64
	TotalRun    time.Duration
}

func (gc *GeneralCounter) AddTrans(time float64) {

	gc.TransTime = append(gc.TransTime, time)

}

func (gc *GeneralCounter) Results(parseHeader *CompareHeader) {

	fmt.Printf("\n\n")
	fmt.Println("Transactions:", gc.NumRequest, "hits")
	fmt.Printf("Availability: %.2f%%\n", 100-(float64(gc.NumBadError)*100/float64(gc.NumRequest)))
	fmt.Println("Elapsed time:", gc.TotalRun.String())
	fmt.Printf("Transaction rate: %.2f\n", float64(gc.NumSuccess)/gc.TotalRun.Seconds())
	fmt.Println("Successful transactions:", gc.NumSuccess)
	fmt.Println("Failed transactions:", gc.NumRequest-gc.NumSuccess)
	fmt.Printf("Response time: %.2fs\n", gc.TotalTime/float64(gc.NumRequest-gc.NumBadError))
	fmt.Printf("Longest transaction: %.2fs\n", gc.LongTrans)
	fmt.Printf("Shortest transaction: %.2fs\n", gc.ShortTrans)
	fmt.Println("Throughput:", ByteSize(gc.TotalByte/gc.TotalRun.Seconds()))
	fmt.Println("Average bytes for transaction: ", ByteSize(gc.TotalByte/float64(gc.NumRequest-gc.NumBadError)))

	for _, value := range parseHeader.list {
		fmt.Printf("\n\n")
		fmt.Printf("Header %s: '%s'\n", value.Key, value.Value)
		fmt.Println("Transactions with this header: ", value.ContTot)
		fmt.Printf("Was present in %.1f%% of total transactions\n", float64(value.ContTot)*100/float64(gc.NumSuccess))
		fmt.Printf("Match the regexp %.1f%% transactions\n", float64(value.ContHit)*100/float64(value.ContTot))

	}

}

type CompareHeader struct {
	list        map[string]*FilterHeader
	PrintRegexp bool
}

func NewCompareHeader() *CompareHeader {
	return &CompareHeader{
		make(map[string]*FilterHeader),
		false,
	}
}

func (ch *CompareHeader) Add(key, value string) {

	if r, err := regexp.Compile(value); err != nil {
		fmt.Printf("Unable to compile '%s' in a regular expression, skipped\n", value)
	} else {
		ch.list[key] = &FilterHeader{
			Key:     key,
			Value:   value,
			ContTot: 0,
			ContHit: 0,
			Rexp:    r,
		}
	}

}

func (ch *CompareHeader) CompareAll(header http.Header) {

	for _, val := range ch.list {
		val.Compare(header, ch.PrintRegexp)
	}

}

func (ch *CompareHeader) String() string {

	app := ""
	for key, value := range ch.list {
		app += fmt.Sprintf("%s: %s\n", key, value)
	}
	return app

}

type FilterHeader struct {
	Key     string
	Value   string
	ContTot int
	ContHit int
	Rexp    *regexp.Regexp
	sync.Mutex
}

func (fh *FilterHeader) Compare(header http.Header, printRegexp bool) {

	value := header.Get(fh.Key)

	if value == "" {
		return
	}

	fh.Lock()
	defer fh.Unlock()

	if printRegexp {
		fmt.Printf("\t%s: '%s'\n", fh.Key, value)
	}

	fh.ContTot++

	if fh.Rexp.Match([]byte(value)) {
		fh.ContHit++
	}

}

type SimpleCounter struct {
	QtaBytes   float64
	Elapsed    float64
	StatusCode int
	Path       string
	Header     http.Header
	Error      error
}

func NewSimpleCounter(qtaBytes float64, elapsedTime float64, code int, path string, header http.Header) *SimpleCounter {

	app_path := "/"
	if path != "" {
		app_path = path
	}

	return &SimpleCounter{
		qtaBytes,
		elapsedTime,
		code,
		app_path,
		header,
		nil,
	}
}

func ProcessData(dataChannel chan *SimpleCounter, HC *CompareHeader, waitGroup *sync.WaitGroup) {

	var ok bool
	var data *SimpleCounter

	sumData := &GeneralCounter{}
	start := time.Now()

	defer waitGroup.Done()

	for {

		select {
		case data, ok = <-dataChannel:

			if !ok {

				sumData.TotalRun = time.Since(start)
				sumData.Results(HC)
				return

			}

			// sum request
			sumData.NumRequest++

			if data.Error != nil {

				// sum bad error, socket, system limit
				sumData.NumBadError++
				fmt.Println(data.Error)

			} else {

				fmt.Println(data.StatusCode, fmt.Sprintf("%.2fs", data.Elapsed), ByteSize(data.QtaBytes), data.Path)
				HC.CompareAll(data.Header)

				// qta bytes
				sumData.TotalByte += data.QtaBytes

				sumData.AddTrans(data.Elapsed)

				// if status code <400 it's a success request
				if data.StatusCode < 400 {
					sumData.NumSuccess++
				}

				// save the shortest request
				if sumData.ShortTrans == 0 || sumData.ShortTrans > data.Elapsed {
					sumData.ShortTrans = data.Elapsed
				}

				// save the longest request
				if sumData.LongTrans == 0 || sumData.LongTrans < data.Elapsed {
					sumData.LongTrans = data.Elapsed
				}
				// sum the total time
				sumData.TotalTime += data.Elapsed

			}

		default:
		}
	}
}
