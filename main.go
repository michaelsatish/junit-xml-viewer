package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
)

// xmlPath is the path to the JUnit XML file.
// renderMethod is the method to render the JUnit XML file.
// outputFile is the path to the output file.
// serverPort is the port to serve the dashboard on.
var (
	xmlPath      string
	renderMethod string
	outputFile   string
	serverPort   string
)

type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Errors    string     `xml:"errors,attr"`
	Failures  string     `xml:"failures,attr"`
	Skipped   string     `xml:"skipped,attr"`
	Tests     string     `xml:"tests,attr"`
	Time      string     `xml:"time,attr"`
	Timestamp string     `xml:"timestamp,attr"`
	Hostname  string     `xml:"hostname,attr"`
	TestCases []TestCase `xml:"testcase"`
}

// GetSuccessCount returns the number of successful tests.
func (ts *TestSuite) GetSuccessCount() int {
	tests, err := strconv.Atoi(ts.Tests)
	checkError(err)

	failures, err := strconv.Atoi(ts.Failures)
	checkError(err)

	errors, err := strconv.Atoi(ts.Errors)
	checkError(err)

	skipped, err := strconv.Atoi(ts.Skipped)
	checkError(err)

	return tests - failures - errors - skipped
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Classname string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	Failure   Failure  `xml:"failure,omitempty"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
}

// checkError checks if an error occurred and if so, it logs it and exits the program.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Render the dashboard. The method is determined by the renderMethod flag.
// Possible values: server, stdout, output.
// Server: Render the dashboard on the browser.
// Stdout: Render the dashboard to the stdout.
// Output: Render the dashboard to the output file.
func render(method string, tmpl *template.Template, testSuites *TestSuites) error {
	switch method {
	case "server":
		// Create a new temporary html file.
		tmpFile, err := ioutil.TempFile("", "dashboard.*.html")
		if err != nil {
			return err
		}

		defer tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		// Write the dashboard template to the tmp file.
		err = tmpl.Execute(tmpFile, testSuites)
		if err != nil {
			return err
		}

		// Open the temporary html file in the web browser.
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, tmpFile.Name())
		})
		log.Printf("Serving the dashboard at http://localhost:%s/", serverPort)
		log.Println("Press CTRL+C to stop the server.")

		err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
		if err != nil {
			return err
		}
	case "stdout":
		// Write the dashboard template to the stdout.
		err := tmpl.Execute(os.Stdout, testSuites)
		if err != nil {
			return err
		}
	case "output":
		// Write the dashboard template to the output file.
		f, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer f.Close()

		err = tmpl.Execute(f, testSuites)
		if err != nil {
			return err
		}
	default:
		return errors.New("Please provide a valid render method.")
	}

	return nil
}

func main() {
	// Parse the command line flags.
	flag.StringVar(&xmlPath, "f", "", "Path to the JUnit XML file.")
	flag.StringVar(&renderMethod, "m", "server", "Method to render the JUnit XML file. Possible values: server, stdout, output.")
	flag.StringVar(&outputFile, "o", "", "Path to the output file.")
	flag.StringVar(&serverPort, "p", "8080", "Port to serve the dashboard on.")
	flag.Parse()

	// Check if the xml file exists and read the contents.
	if xmlPath == "" {
		log.Fatal("Please provide a path to the JUnit XML file.")
	}

	if _, err := os.Stat(xmlPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("The JUnit XML file does not exist in the given path.")
	}

	// If the render method is output, check if the output file exists.
	if renderMethod == "output" && outputFile == "" {
		log.Fatal("Please provide a valid output file path. Please use the -o flag to provide output file path.")
	}

	data, err := ioutil.ReadFile(xmlPath)
	checkError(err)

	// Unmarshal the xml file contents into a TestSuites struct.
	var testSuites TestSuites
	err = xml.Unmarshal(data, &testSuites)
	checkError(err)

	// Create a new dashboard template.
	tmpl, err := template.ParseFiles("dashboard.html")
	checkError(err)

	// Render the dashboard template.
	err = render(renderMethod, tmpl, &testSuites)
	checkError(err)
}
