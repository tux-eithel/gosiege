# GOSIEGE, http/https stress tester made in GO
***

Inspired by [Siege](http://www.joedog.org/siege-home/) I have tried to do something similar. 
I have added some useful features like:
* filter response header
* send custom header in request
* send custom param in POST request

### Install

```
go get -u github.com/tux-eithel/gosiege/
```


### Command Line

```
-c x
	Where x is the number of concurrent connections. Default 1

-exp "HeaderField Value"
	Filters header using regular expression. You can define multiple regexp.
	HeaderField is a header field, Value must be a regular expression. 
	Examples are:
		`-exp "X-Cache HIT"`
		`-exp "X-Cache .*"`

-f fileName
	File which contains at least an url for every row 

-nasty=boolValue
	Uses all available CPUs, physic and logic. Default true

-per y	
	Test every url y times. When all url are been tested y time, gosiege quits.
	Default -1, so test will run until you press Ctrl+c	
	
-pexp
	Prints the header which matches HeaderField.
	If a header matches, it will print in green color, otherwise magenta is used
	
-rand=boolValue
	Takes random url from the list. Default is true. False may be useful 
	if you want to test the list sequentially

-s 
	Times to wait between each requests. 
	Values accepted according http://golang.org/pkg/time/#Duration. Default is 1s

-u
	List of urls (spaced separated) passed via CLI. It can be used along -f option
```


### Examples

Reads input file and pass some url via CLI. Starts 5 concurrent connection and hit
every url 2 times then quits
* `gosiege -f test -u "google.com https://duckduckgo.com/" -c 5 -per 2`

Reads urls from CLI and runs 40 concurrent connections. Between calls 30ms will be waited
* `gosiege -u "google.com https://duckduckgo.com/" -c 40 -s 30ms`

Reads urls from CLI, checks if in header is present "Server" header and prints the result
* `gosiege -u "google.com https://duckduckgo.com/about" -c 2 -exp "Server gws" -pexp`


### Output

Commnad: `gosiege -u "google.com https://duckduckgo.com/about" -c 2 -exp "Server gws" -pexp`
```
.
.
.
200 0.20s 46.09KB /about
	Server: 'nginx'		<-- this is printed because -pexp
200 0.29s 17.55KB /
	Server: 'gws'		<-- this is printed because -pexp
.
.
.
Received quit. Sending shutdown and waiting all goroutines...
.
.
.



Transactions: 10 hits
Availability: 100.00%
Elapsed time: 7.106084685s
Transaction rate: 1.41
Successful transactions: 10
Failed transactions: 0
Response time: 0.37s
Longest transaction: 0.75s
Shortest transaction: 0.12s
Throughput: 44.78KB
Average bytes for transaction:  31.82KB


Header Server: 'gws'
Transactions with this header:  10
Was present in 100.0% of total transactions
Match the regexp 50.0% transactions
Done.

```


### Statistics
* Transactions: number of server hits
* Availability: percentage connections successfully handled by the server. It's not included 40x and 50x errors
* Elapsed time: duration of entire test
* Transaction rate: transactions / elapsed time
* Successful transactions: hits with code < 400
* Failed transactions: hits with code >=400 (not included socket errors or timeouts)
* Response time: average time to respond to each requests
* Longest transaction: the slowest hit
* Shortest transaction: the quickest hit
* Throughput: average number of bytes transferred every second
* Average bytes for transaction: average number of bytes transferred for request

If some regular expression is defined using `-exp`, for every regular expression:
* name of HeaderField and Value of regular expression
* number of transaction where HeaderField is present
* percentage of connections with HeaderField present on total transactions
* percentage of connections which matches the regular expression over transactions where HeaderField was present


### Input File
Input file can contains different urls. Every url may have some parameters.
If an url doesn't match one of the below example, it will be skipped.
Parameters are spaced separated.

#### Valid urls are:
Simple GET request: automatically gosiege add http://
* `example.com`

GET request with some header. Header must be a valid json object (simple, only string)
* `http://example.com {"Cache-Control":"max-age=0"}`

Simple POST request: 
* `www.example.com POST`

POST request with some header. Header must be a valid json object (simple, only string)
* `www.example.com POST {"Cache-Control":"max-age=0"}`

POST request with header and parameters
* `www.example.com POST {"Content-Type":"application/x-www-form-urlencoded"} id=2&user=100`

POST request with only parameters
* `www.example.com POST {} id=2&user=100`


### Tips
* Ctrl+c for exit
* You can use -f and -u together


### Bugs
Yeah there are bugs... help me fix them :) !