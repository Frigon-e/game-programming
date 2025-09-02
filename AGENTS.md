# Game of Life Agent Guide

Welcome, agent! This guide will help you contribute to the Game of Life project.

## Project Overview

This project is a simulation of Conway's Game of Life, built in Go using the Ebiten 2D game engine. The core game logic is separated from the rendering and application startup.

### Directory Structure

- `cmd/main.go`: The main entry point for the application.
- `gameoflife/`: Contains the primary module for the Game of Life simulation.
- `gameoflife/internal/ddd/board.go`: Defines the game board and its core logic, such as counting live neighbors.
- `go.mod`, `go.sum`: Manage the project's dependencies.

## Development

### Running the Application

To run the Game of Life simulation, execute the following command from the root directory:

```bash
go run cmd/main.go
```

### Running Tests

To run the unit tests, use the following command:

```bash
go test ./...
```
