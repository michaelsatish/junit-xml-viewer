package main

import (
	"embed"
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
// serverPort is the port to serve the dashboard on.
// export is the flag to render the dashboard to stdout.
var (
	xmlPath    string
	serverPort string
	export     bool
	version    string
	vFlag      bool
)

//go:embed dashboard.html
var f embed.FS

type TestSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Errors    string     `xml:"errors,attr"`
	Failures  string     `xml:"failures,attr"`
	Skipped   string     `xml:"skipped,attr,omitempty"`
	Tests     string     `xml:"tests,attr"`
	Time      string     `xml:"time,attr"`
	Timestamp string     `xml:"timestamp,attr"`
	Hostname  string     `xml:"hostname,attr"`
	TestCases []TestCase `xml:"testcase"`
}

// GetSuccessCount returns the number of successful tests.
func (ts *TestSuite) GetSuccessCount() int {
	intCov := func(s string) int {
		i, err := strconv.Atoi(s)
		checkError(err)
		return i
	}

	tests := intCov(ts.Tests)
	failures := intCov(ts.Failures)
	errors := intCov(ts.Errors)

	// Not all test suites have skipped tests.
	if ts.Skipped == "" {
		return tests - failures - errors
	}

	skipped := intCov(ts.Skipped)
	return tests - failures - errors - skipped
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      string   `xml:"time,attr"`
	Failure   Failure  `xml:"failure,omitempty"`
	Error     Error    `xml:"error,omitempty"`
	Skipped   Skipped  `xml:"skipped,omitempty"`
}

type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Value   string   `xml:",chardata"`
}

type Error struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:"message,attr"`
	Value   string   `xml:",chardata"`
}

type Skipped struct {
	XMLName xml.Name `xml:"skipped"`
	Type    string   `xml:"type,attr"`
	Message string   `xml:"message,attr"`
}

// checkError checks if an error occurred and if so, it logs it and exits the program.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// serve serves the dashboard.
func serve(tmpl *template.Template, ts *TestSuite) error {
	// Create a new temporary html file.
	tmpFile, err := ioutil.TempFile("", "dashboard.*.html")
	if err != nil {
		return err
	}

	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// Write the dashboard template to the tmp file.
	err = tmpl.Execute(tmpFile, ts)
	if err != nil {
		return err
	}

	// Serve the dashboard.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, tmpFile.Name())
	})
	log.Printf("Serving the dashboard at http://localhost:%s/", serverPort)
	log.Println("Press CTRL+C to stop the server.")

	err = http.ListenAndServe(fmt.Sprintf(":%s", serverPort), nil)
	if err != nil {
		return err
	}

	return nil
}

// expStdout renders the dashboard to stdout.
func expStdout(tmpl *template.Template, ts *TestSuite) error {
	err := tmpl.Execute(os.Stdout, ts)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// Parse the command line flags.
	flag.StringVar(&xmlPath, "f", "", "Path to the JUnit XML file.")
	flag.StringVar(&serverPort, "p", "8080", "Port to serve the dashboard on.")
	flag.BoolVar(&export, "e", false, "Render to stdout.")
	flag.BoolVar(&vFlag, "v", false, "Print the version.")
	flag.Parse()

	// Print the version.
	if vFlag {
		fmt.Println("Version:", version)
		os.Exit(0)
	}

	// Check if the xml file exists and read the contents.
	if xmlPath == "" {
		log.Fatal("Please provide a path to the JUnit XML file.")
	}

	if _, err := os.Stat(xmlPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal("The JUnit XML file does not exist in the given path.")
	}

	data, err := os.ReadFile(xmlPath)
	checkError(err)

	// Unmarshal the xml file contents into a TestSuites struct.
	var testSuites TestSuites
	err = xml.Unmarshal(data, &testSuites)
	checkError(err)

	// Create a new dashboard template.
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
	}

	b, err := f.ReadFile("dashboard.html")
	checkError(err)

	tmpl, err := template.New("dashboard").Funcs(funcMap).Parse(string(b))
	checkError(err)

	// Render the dashboard template.
	ts := testSuites.TestSuites[0]
	if export {
		err := expStdout(tmpl, &ts)
		checkError(err)
	} else {
		err := serve(tmpl, &ts)
		checkError(err)
	}
}
