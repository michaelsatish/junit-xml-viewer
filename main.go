package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
)

var xmlPath string

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
	Timestamp string     `xml:"timestamp,attr"`
	Hostname  string     `xml:"hostname,attr"`
	TestCases []TestCase `xml:"testcase"`
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

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	flag.StringVar(&xmlPath, "f", "", "Path to the JUnit XML file.")
	flag.Parse()

	// Check if the xml file exists and read the contents.
	if xmlPath == "" {
		log.Fatal("Please provide a path to the JUnit XML file.")
	}

	if _, err := os.Stat(xmlPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("The JUnit XML file does not exist in the given path.")
	}

	data, err := ioutil.ReadFile(xmlPath)
	checkError(err)

	// unmarshal the xml file contents into a TestSuites struct.
	var testSuites TestSuites
	err = xml.Unmarshal(data, &testSuites)
	checkError(err)

	// Create a new dashboard template.
	tmpl, err := template.ParseFiles("dashboard.html")
	checkError(err)

	// Create a new temporary html file.
	tmpFile, err := ioutil.TempFile("", "dashboard.*.html")
	checkError(err)
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Write the dashboard template to the tmp file.
	err = tmpl.Execute(tmpFile, testSuites)
	checkError(err)

	// Open the temporary html file in the web browser.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, tmpFile.Name())
	})
	log.Println("Serving the dashboard at http://localhost:8080/")
	log.Println("Press CTRL+C to stop the server.")

	err = http.ListenAndServe(":8080", nil)
	checkError(err)
}
