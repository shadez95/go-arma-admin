# go-arma-admin

[![Go Report Card](https://goreportcard.com/badge/github.com/shadez95/go-arma-admin)](https://goreportcard.com/report/github.com/shadez95/go-arma-admin) [![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

A web service for controlling and managing ARMA servers through a web API. The web API consists of RESTful API endpoints and websockets to manage users and ARMA servers.

## Status

Currently this is under heavy development and API's are subject to change without notice. I will work on a front-end in the future but getting the backend working and stable is priority first.

## Contribute

If you wish to contribute, feel free to contribute, but please contact me first through Twitter PM, email, submit a ticket. This project will be under constant development and don't want people to be wasting time and effort over something I may be working on already, or is out of the scope of this project.

This project uses [dep](https://github.com/golang/dep) for dependency management.

If you want a better debugging experience, especially if you are using VS Code, it is highly recommended to install [delve](https://github.com/derekparker/delve).

Tested in Go 1.12.1.

To install locally for development run `go get -u github.com/shadez95/go-arma-admin`
