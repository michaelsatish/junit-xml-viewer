# JUnit XML Viewer

A tiny cli application that can read a JUnit XML file and makes it readable.

## Usage

```text
Options:
  -f Path to the JUnit XML file.
  -e To render the dashboard to stdout.
  -p Port to serve the dashboard on.
```

By default, the JUnit XML file is rendered to an HTML document and served on the default port 8080. 

```sh
jxv -f test.xml

Serving the dashboard at http://localhost:8080/
Press CTRL+C to stop the server.
```

You can write the HTML document to stdout by passing the -e flag.

```sh
jvx -f test.xml -e
```