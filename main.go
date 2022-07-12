package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
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

	// Check if the xml file exists.
	if xmlPath == "" {
		log.Fatal("Please provide a path to the JUnit XML file.")
	}

	if _, err := os.Stat(xmlPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("The JUnit XML file does not exist in the given path.")
	}

	data, err := ioutil.ReadFile(xmlPath)
	checkError(err)

	// unmarshal the xml file into a TestSuites struct.
	var testSuites TestSuites
	err = xml.Unmarshal(data, &testSuites)
	checkError(err)

	// Create a new dashboard template.
	tmpl, err := template.ParseFiles("dashboard.html")
	checkError(err)

	// Create a new tmp file.
	tmpFile, err := ioutil.TempFile("junit-xml-viewer", "dashboard.*.html")
	checkError(err)
	defer os.Remove(tmpFile.Name())

	// Write the dashboard template to the tmp file.
	err = tmpl.Execute(tmpFile, testSuites)
	checkError(err)

	// Open the tmp file in the web browser.
	fmt.Println(tmpFile.Name())
}
