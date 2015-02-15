// parseUrl let you to parse urls from command line or a file (passed via command line)
package parseUrl

// RandomUri is an interface for retrive a random string form a object
type RandomUri interface {
	GetRandomUri() (*string, error)
}
