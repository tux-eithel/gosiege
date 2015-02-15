package parseUrl

import (
	"fmt"
	"net/url"
	"strings"
)

// A Urls rappresent an array of strings
type Urls []string

// String returns all urls in a single string
func (u *Urls) String() string {
	return fmt.Sprint(*u)
}

// Set splits input string by space and use them as url. If string is not a valid url, it'll be skipped
func (u *Urls) Set(srt string) error {
	app := strings.Split(srt, " ")

	for _, value := range app {
		_, err := url.Parse(value)
		if err != nil {
			fmt.Println("Url: ", value, "not correct, skipped")
		}
		*u = append(*u, value)
	}
	return nil
}

// Return a random url
func (u Urls) GetRandomUri() (*string, error) {
	return &(u)[0], nil
}
