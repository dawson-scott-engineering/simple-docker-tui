# Simple Docker TUI

A basic Docker management TUI written in Go for quickly managing containers used in networking labs.

## Features

* List Docker containers and view status
* Start, stop, restart, and delete containers
* Deploy Debian lab containers with automatic naming and port assignment

## Requirements

* Go
* Docker
* `lab-debian:latest` image for container deployment

## Usage

```bash
go run .
```

## Notes

New containers use the `lab-debian:latest` image and map container SSH port `22` to an available host port starting at `10000`.

Built as a utility for managing Docker infrastructure in networking labs.
