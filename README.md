# config-server
My attempt at builder a configuration management server in Golang

# Setup ⚙️
- Rename `example.config.env` to `config.env`
- Populate relevant configuration variables
- Setup authentication for github by setting up personal access token for your account and grant necessary permissions. Check the [GitHub Docs](https://docs.github.com/en/rest/overview/other-authentication-methods?apiVersion=2022-11-28) for guidance
- Setup a new github repo to store your configuration files in YAML.

| Variable Name | Description | 
| --- | --- |
| REPO_NAME | The name of your github repository created to store your env variables |
| OWNER | Your github username |
| AUTH_TOKEN | Personal access token used for GitHib Authentications |

# Run
- Run `go run main.go` to start server

# Usage :rocket:
Configs will be served at

| Name | Request Method | Endpoint  | Response |
| --- | --- | --- | --- |
| ServeConfig | GET | localhost:8080/{filepath}/{branch} | 200 OK |

Configs will be served as JSON



## Stil a WIP :anchor:


