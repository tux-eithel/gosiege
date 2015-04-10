package libgosiege

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ErrorFormat = "Formatting error"
)

func ParseAllInputFile(fileName string, request *Requests) {

	var err error
	var fd *os.File

	fd, err = os.Open(fileName) // For read access.
	if err != nil {
		log.Fatal("Error open Input file:", err)
	}
	defer fd.Close()

	scannerInput := bufio.NewScanner(fd)

	var row []string
	var rawRow string
	var appR *InputRequest
	var url string
	var jsonHeader map[string]string
	contRow := 1

	for scannerInput.Scan() {

		rawRow = scannerInput.Text()
		row = strings.Fields(rawRow)

		switch len(row) {

		case 1:

			url = row[0]
			appR, err = NewInputRequest(url)

		case 2:

			url = row[0]

			switch row[1] {

			case "POST":
				appR, err = NewInputRequestComplex(url, row[1], nil, nil)

			default:

				url = row[0]
				err = json.Unmarshal([]byte(row[1]), &jsonHeader)
				if err != nil {
					fmt.Println("Row - ", contRow, "Error parsing header:", err, " - ignored")
				}
				appR, err = NewInputRequestComplex(url, "GET", nil, jsonHeader)

			}

		case 3:
			url = row[0]

			switch row[1] {

			case "POST":
				err = json.Unmarshal([]byte(row[2]), &jsonHeader)
				if err != nil {
					fmt.Println("Row - ", contRow, "Error parsing header:", err, " - ignored")
				}

				appR, err = NewInputRequestComplex(url, row[1], nil, jsonHeader)

			default:
				err = errors.New(ErrorFormat)

			}

		case 4:

			url = row[0]

			switch row[1] {

			case "POST":
				err = json.Unmarshal([]byte(row[2]), &jsonHeader)
				if err != nil {
					log.Fatal("Error parsing header:", err)
				}

				appR, err = NewInputRequestComplex(url, row[1], []byte(row[3]), jsonHeader)

			default:
				err = errors.New(ErrorFormat)

			}

		default:
			err = errors.New(ErrorFormat)

		}

		if err != nil {
			fmt.Println("Row - ", contRow, "Url ignored '", url, "' with error:", err)
		} else {
			request.AddRequest(appR)
		}

		contRow++
		jsonHeader = nil

	}

	if err := scannerInput.Err(); err != nil {
		log.Fatal("Reading standard input:", err)
	}

}
