package libgosiege

import (
	"fmt"
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
