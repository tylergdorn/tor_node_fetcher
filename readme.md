# Tor Nodes Fetcher

## Usage

`go run cmd/main.go --help` to see options.
Otherwise, `go run cmd/main.go` will start a functional server on port `8080`.
This server will by default use "./tor.db" as the sqlite database for the application

## APIs Exposed

- `GET /nodes`
    - returns list of ips, the time they were seen, and the list they've been seen on
        - `torproject` is `https://check.torproject.org/torbulkexitlist`, the tor project's list of tor exitnodes
        - `danmeuk` is `https://www.dan.me.uk/torlist/?exit`, a list of tor exit nodes
    - supports limit/offset pagination with the `limit` and `offset` query parameters
        - e.g. `GET /nodes?limit=10&offset=0`
- `GET /allowlist` returns the ips allowlist of excluded ips from the tor exit list node above
- `POST /allowlist/<ip address here>` adds a new ip to the allow list
    - using the `note` query parameter lets you specify a note to the ip
- `DELETE /allowlist/<ip address here>` removes an ip from the allowlist


## Building

This project does require CGO, and you will need `CGO_ENABLED=1` set to build it. (as well as all the requirements that apply to CGO)