package main

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_getFileInfo(t *testing.T) {
	// Define the structure of our test case
	tests := []struct {
		name    string
		want    inputFile
		wantErr bool
		osArgs  []string
	}{
		// Define our test cases
		{"Default parameters", inputFile{"test.csv", ",", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Semicolon enabled", inputFile{"test.csv", ";", false}, false, []string{"cmd", "--separator=;", "test.csv"}},
		{"Pretty enabled", inputFile{"test.csv", ",", true}, false, []string{"cmd", "--pretty", "test.csv"}},
		{"Pretty and Semicolon enabled", inputFile{"test.csv", ";", true}, false, []string{"cmd", "--pretty", "--separator=;", "test.csv"}},
		{"Separator not identified", inputFile{}, true, []string{"cmd", "--separator=|", "test.csv"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualOsArgs := os.Args
			// This function will run after the test is done
			defer func() {
				os.Args = actualOsArgs
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs
			got, err := getFileInfo()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFileInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFileInfo() = %v, want = %v", got, tt.want)
			}
		})
	}
}

func checkTest(err error) {
	if err != nil {
		panic(err)
	}
}

func Test_isValidCSVFile(t *testing.T) {
	// Creating an empty CSV file
	tmpfile, err := os.CreateTemp("", "test*.csv")
	checkTest(err)
	// Remove the file after the test
	defer os.Remove(tmpfile.Name())
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"File does exist", args{tmpfile.Name()}, true, false},
		{"File doesn't exist", args{"nowhere/test.csv"}, false, true},
		{"File is not csv", args{"toto.txt"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isValidCSVFile(tt.args.filename)
			// Checking the error
			if (err != nil) != tt.wantErr {
				t.Errorf("isValidCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Checking the returning value
			if got != tt.want {
				t.Errorf("isValidCSVFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processCSVFile(t *testing.T) {
	wantResult := []map[string]string{
		{"Name": "Stan", "Age": "23", "Height": "185"},
		{"Name": "Dan", "Age": "32", "Height": "175"},
	}

	tests := []struct {
		name           string
		csvFileContent string
		separator      string
	}{
		{"Comma separator", "Name,Age,Height\nStan,23,185\nDan,32,175\n", ","},
		{"Semicolon separator", "Name;Age;Height\nStan;23;185\nDan;32;175\n", ";"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating an empty CSV file
			tmpfile, err := os.CreateTemp("", "test*.csv")
			checkTest(err)
			// Remove the file after the test
			defer os.Remove(tmpfile.Name())
			tmpfile.Sync()

			_, err = tmpfile.WriteString(tt.csvFileContent)
			fileInfo := inputFile{
				filepath:  tmpfile.Name(),
				separator: tt.separator,
				pretty:    false,
			}
			writer := make(chan map[string]string)
			// Call the function as go routine
			go processCSVFile(fileInfo, writer)
			for _, item := range wantResult {
				record := <-writer
				if !reflect.DeepEqual(record, item) {
					t.Errorf("processCSVFile() = %v, want %v", record, item)
				}
			}
		})
	}
}

func Test_writeJSONFile(t *testing.T) {
	transformedData := []map[string]string{
		{"Name": "Stan", "Age": "23", "Height": "185"},
		{"Name": "Dan", "Age": "32", "Height": "175"},
	}

	tests := []struct {
		name         string
		csvFilePath  string
		jsonFilePath string
		pretty       bool
	}{
		{"Compact JSON Output", "compact.csv", "compact.json", false},
		{"Pretty JSON Output", "pretty.csv", "pretty.json", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Creating a chanel
			writerChanel := make(chan map[string]string)
			done := make(chan bool)
			// Run a go-routine to write in the channel
			go func() {
				defer close(writerChanel)
				for _, record := range transformedData {
					writerChanel <- record
				}
			}()
			// Running the function
			writeJSONFile(tt.csvFilePath, writerChanel, done, tt.pretty)
			// waiting for the end of running
			<-done
			// getting the test result from file
			testRes, err := os.ReadFile(tt.jsonFilePath)
			if err != nil {
				t.Errorf("writeJSONFile(), Output file got error: %v", err)
			}
			// remove the file afterward
			defer os.Remove(tt.jsonFilePath)
			wantedRes, err := os.ReadFile(filepath.Join("testFiles", tt.jsonFilePath))
			check(err)
			if string(testRes) != string(wantedRes) {
				t.Errorf("writeJSONFile() = %v, want %v", string(testRes), string(wantedRes))
			}
		})
	}
}
