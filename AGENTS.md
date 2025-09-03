# Agent Guide: SideProjectGames

Welcome, agent! This guide will help you contribute to the SideProjectGames repository.

## Project Overview

This project is a collection of small game modules written in Go, using the Ebiten 2D game engine for rendering. The goal is to create a modular system where different games can be developed and run from a common entry point.

The currently implemented modules are:
1.  **Conway's Game of Life**: An interactive cellular automaton.
2.  **Battleship**: A classic board game against a simple AI.

### Directory Structure

- `cmd/main.go`: The main entry point for the application. It reads the `MODULE` environment variable to decide which game to run.
- `gameoflife/`: Contains the primary module for the Game of Life simulation.
- `battleship/`: Contains the primary module for the Battleship game.
- `internal/`: Holds shared code, including configuration and a generic board implementation.
- `go.mod`, `go.sum`: Manage the project's dependencies.

## Development

### Running the Application

To run a specific game module, you must set the `MODULE` environment variable.

**To run the Game of Life simulation:**
```bash
MODULE=GOL go run ./cmd/main.go
```

**To run the Battleship game:**
```bash
MODULE=BATTLESHIP go run ./cmd/main.go
```

### Running Tests

To run all unit tests for the project, use the following command from the root directory:

```bash
go test ./...
```
