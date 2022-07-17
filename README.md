# JUnit XML Viewer

A tiny cli application that can read a JUnit XML file and makes it readable.

## Usage

```text
Options:
  -f Path to the JUnit XML file.
  -m Method to render the JUnit XML file. Acceptable values: server, stdout, output.
  -o Path to the output file.
  -p Port to serve the dashboard on.
```

By default, the JUnit XML file is rendered to an HTML document and served on the default port 8080. 

```sh
jxv -f test.xml

Serving the dashboard at http://localhost:8080/
Press CTRL+C to stop the server.
```

You can write the HTML document to stdout by changing the render method to "stdout".

```sh
jvx -f test.xml -m stdout
```

You can write the HTLM document to a file by changing the render method to "output" and passing in the output path.

```sh
jvx -f test.xml -m output -o dashboard.html
```