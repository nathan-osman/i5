## ![](https://i.stack.imgur.com/08e0R.png)

[![Build Status](https://ci.quickmediasolutions.com/buildStatus/icon?job=i5)](https://ci.quickmediasolutions.com/job/i5/)
[![GoDoc](https://godoc.org/github.com/nathan-osman/i5?status.svg)](https://godoc.org/github.com/nathan-osman/i5)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

i5 is a reverse proxy for web services running in Docker.

### Features

- Monitors the Docker daemon for containers starting and stopping
- Routes traffic to running containers based on labels
- Serves static files directly
- Provides status pages for monitoring services
- Automatically obtains TLS certificates and redirects HTTP traffic
- Creates and initializes MySQL and PostgreSQL databases on-demand
- Runs within its own Docker container and requires very little configuration

### Building the App

i5 must be built in two steps.

#### Building the UI

The web interface uses [React](https://reactjs.org/) and must be built with npm. This can be done by running the following command in the `ui/` directory:

```shell
npm run build
```

The resulting files can be found in `ui/build/`.

#### Compiling the Application

The server application itself is written in [Go](https://golang.org/) and can be compiled with:

```shell
go build
```
