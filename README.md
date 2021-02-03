# Simple Web Server

A simple web server written in Go. The server exposes the following endpoints:
- `/`                     responds with `echo`
- `/headers`              responds with the request headers in JSON
- `/env`                  responds with the environment variables in JSON
- `/encrypt`              responds with an encrypted message
- `/decrypt/<ciphertext>` responds with the plaintext of `<ciphertext>`

## How to run
1. Install Go https://golang.org/doc/install.
2. Clone this repository and `cd` into the directory.
3. Build the server with `go build -o bin/server main.go`.
4. Run `bin/server`.
5. Test the connection by requesting `http://localhost/`, e.g. `curl http://localhost`.

## Configuration

### Command-line arguments
- `-host`   the host to listen on (defaults to `127.0.0.1`)
- `-port`   the port to listen on (defaults to `80`)
- `-config` the server config file in JSON format (defaults to `config.json`)
