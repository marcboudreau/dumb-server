# dumb-server
A simple dumb HTTP server useful for testing.  This server completely ignores
the incoming request parameters.  It always responds with the same response.

## Usage

`dumb-server [-port PORT] [-sc STATUS_CODE] [-resp RESPONSE_FILE]`

Where:
* PORT is the TCP port number that the server will listen on.  Defaults to 7979.
* STATUS_CODE is the HTTP status code of the response.  Defaults to 200.
* RESPONSE_FILE is the path to a response file that contains the data used to build the response.

The response file format consists of a header section and the body section.
The sections are separated by a single blank line.

```
header:header value

body section
continues
until the end of the file
```

If a response file is not specified, the server will respond with a single
HTTP header: `Content-Type: text/plain` and the following body: `dumb-server default response`. 

## Docker Usage

The server can be easily run inside a docker container as well.  This command
will run the server with its default configuration, as described above.

`docker run -d -p 7979:7979 marcboudreau/dumb-server`

The following command shows how to specify different values for all
configuration parameters.

`docker run -d -v /path/to/response_file:/response -p 5678:5678 marcboudreau/dumb-server -port 5678 -sc 404 -resp /response`

If specifying a response file (with the `-resp` option), don't forget to mount
that file into the container using the `-v` option of the Docker run command.
