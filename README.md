# config-server
My attempt at builder a configuration management server in Golang

# Setup ⚙️
- Rename `example.config.env` to `config.env`
- Populate relevant configuration variables

# Run
- Run `go run main.go` to start server

# Usage :rocket:
Configs will be served at

| Name | Request Method | Endpoint  | Response |
| --- | --- | --- | --- |
| ServeConfig | GET | localhost:8080/{filepath}/{branch} | 200 OK |

Configs will be served as JSON



## Stil a WIP :anchor:


