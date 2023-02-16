package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type inputFile struct {
	filepath  string
	separator string
	pretty    bool
}

func getFileInfo() (inputFile, error) {
	// Validating the number of given arguments
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required")
	}

	// Parsing the command option
	separator := flag.String("separator", ",", "Column separator")
	pretty := flag.Bool("pretty", false, "Generate pretty JSON")
	flag.Parse()

	// Getting the file location path
	fileLocation := flag.Arg(0)

	// Validating the given separator
	if !(*separator == "," || *separator == ";") {
		return inputFile{}, errors.New("only comma ',' or semicolon ';' separators are allowed")
	}

	return inputFile{fileLocation, *separator, *pretty}, nil
}

func isValidCSVFile(filename string) (bool, error) {
	if fileExtension := filepath.Ext(filename); fileExtension != ".csv" {
		return false, fmt.Errorf("file %s is not a CSV file", filename)
	}

	if _, err := os.Stat(filename); err != nil && os.IsNotExist(err) {
		return false, fmt.Errorf("file %s does not exist", filename)
	}

	return true, nil
}

func properlyExit(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	os.Exit(1)
}

func check(err error) {
	if err != nil {
		properlyExit(err)
	}
}

func processLine(headers []string, columns []string) (map[string]string, error) {
	if len(columns) != len(headers) {
		return nil, errors.New("line doesn't match headers format. Skipping")
	}

	recordMap := make(map[string]string)
	for i, name := range headers {
		recordMap[name] = columns[i]
	}

	return recordMap, nil
}

func processCSVFile(fileInfo inputFile, writerChannel chan<- map[string]string) {
	// Opening the file for reading
	file, err := os.Open(fileInfo.filepath)
	// checking for errors, but not necessary
	check(err)
	// We have to close the file once everything is done
	defer file.Close()

	reader := csv.NewReader(file)
	if fileInfo.separator == ";" {
		reader.Comma = ';'
	}

	// Read the header of each column
	headers, err := reader.Read()
	check(err)

	for {
		line, err := reader.Read()
		if err == io.EOF {
			close(writerChannel)
			break
		}
		check(err)

		// Processing a CSV line
		record, err := processLine(headers, line)

		if err != nil {
			log.Printf("Line: %s; Error: %s\n", line, err)
			continue
		}

		// If there is no error, we send the record to the channel
		writerChannel <- record
	}
}

func createStringWriter(path string) func(string, bool) {
	jsonDir := filepath.Dir(path)
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(path), ".csv"))
	finalPathLocation := filepath.Join(jsonDir, jsonName)
	// Opening the JSON file
	file, err := os.Create(finalPathLocation)
	check(err)

	return func(data string, isClose bool) {
		_, err := file.WriteString(data)
		check(err)
		// Close if it's the end
		if isClose {
			_ = file.Close()
		}
	}
}

func getJSONFunc(pretty bool) (func(map[string]string) string, string) {
	if pretty {
		jsonFunc := func(record map[string]string) string {
			jsonData, _ := json.MarshalIndent(record, "  ", "  ")
			return "  " + string(jsonData)
		}
		return jsonFunc, "\n"
	}

	return func(record map[string]string) string {
		jsonData, _ := json.Marshal(record)
		return string(jsonData)
	}, ""
}

func writeJSONFile(csvFilePath string, writerChanel <-chan map[string]string, done chan<- bool, pretty bool) {
	writeString := createStringWriter(csvFilePath)
	jsonFunc, breakline := getJSONFunc(pretty)
	// Log for user information
	fmt.Println("Writing JSON file ...")
	writeString("["+breakline, false)
	first := true

	for {
		// wait until there is records in the channel
		record, more := <-writerChanel
		if more {
			if !first {
				writeString(","+breakline, false)
			} else {
				first = false
			}

			jsonData := jsonFunc(record)
			writeString(jsonData, false)
		} else {
			writeString(breakline+"]", true)
			fmt.Println("Completed!")
			done <- true
			break
		}
	}
}

func main() {
	//Customizing the help display
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <csvFile>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
	// Getting file information
	fileInfo, err := getFileInfo()
	check(err)

	// Validating the given file
	_, err = isValidCSVFile(fileInfo.filepath)
	check(err)

	// Creating the needed chanel
	writerChanel := make(chan map[string]string)
	done := make(chan bool)

	// Running both function in go-routine
	go processCSVFile(fileInfo, writerChanel)
	go writeJSONFile(fileInfo.filepath, writerChanel, done, fileInfo.pretty)

	// waiting until done
	<-done
}
