package libgosiege

import (
	"fmt"
	"sync"
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
	NumRequest int
	NumSuccess int
	LongTrans  float64
	ShortTrans float64
	TotalTime  float64
}

type SimpleCounter struct {
	QtaBytes   float64
	Elapsed    float64
	StatusCode int
	Path       string
}

func NewSimpleCounter(qtaBytes float64, elapsedTime float64, code int, path string) *SimpleCounter {

	app_path := "/"
	if path != "" {
		app_path = path
	}

	return &SimpleCounter{
		qtaBytes,
		elapsedTime,
		code,
		app_path,
	}
}

func ProcessData(dataChannel chan *SimpleCounter, waitGroup *sync.WaitGroup) {

	var ok bool
	var data *SimpleCounter

	sumData := &GeneralCounter{}

	defer waitGroup.Done()

	for {

		select {
		case data, ok = <-dataChannel:

			if !ok {
				fmt.Printf("%+v\n", sumData)
				return
			}

			fmt.Println(data.StatusCode, fmt.Sprintf("%.2fs", data.Elapsed), ByteSize(data.QtaBytes), data.Path)
			// sum request
			sumData.NumRequest++

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

		default:
		}
	}
}
