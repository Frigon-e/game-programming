# SideProjectGames

An experimental Go project showcasing small game/prototyping modules. The project currently includes two games: Conway's Game of Life and Battleship, both rendered with Ebiten (a 2D game library for Go). It uses a shared domain layer for a generic 2D board, specialized game boards, and a simple configuration system driven by environment variables.

## Modules

This project contains two main game modules:

### 1. Conway's Game of Life

A simulation of Conway's Game of Life with interactive controls.

**Features:**
- Real-time simulation using Ebiten for rendering.
- Interactive controls:
  - **Left click:** Toggle the cell under the mouse.
  - **Arrow Down:** Speed up the simulation (decrease step time).
  - **Arrow Up:** Slow down the simulation (increase step time).
- Toroidal wrapping board (edges wrap around).

### 2. Battleship

A classic game of Battleship against a simple AI opponent.

**Features:**
- Player vs. AI gameplay.
- Separate boards for the player and the AI, displayed vertically.
- Turn-based attacking.
- Visual feedback for hits, misses, and sunk ships.
- Simple AI that takes turns automatically.

## Requirements
- Go 1.25 or newer (as declared in `go.mod`).
- A desktop environment supported by Ebiten (macOS, Windows, or Linux with X11/Wayland). For Linux, you may need additional system packages (OpenGL/GLFW support) depending on your distro.

## Getting Started

### Clone and Build
```
git clone <your-fork-or-repo-url> SideProjectGames
cd SideProjectGames
go build ./cmd
```
This creates an executable in the current directory. You can also run modules directly with `go run`.

### Configuration
Configuration is read from environment variables using `envconfig`, with optional support for `.env` files.

**Supported Variables:**
- `MODULE`: (Required) Specifies which game to run. Can be `GOL` or `BATTLESHIP`.
- `GOLWIDTH`, `GOLHEIGHT`: Board dimensions for Game of Life.
- `BATTLESHIPWIDTH`, `BATTLESHIPHEIGHT`: Board dimensions for Battleship.
- `ENVIRONMENT`: Set to `local` to load `.env.local` files.

Example `.env` file:
```
MODULE=BATTLESHIP
BATTLESHIPWIDTH=10
BATTLESHIPHEIGHT=10
```

### Running the Games

To run a specific module, set the `MODULE` environment variable.

**Run Game of Life:**
```
MODULE=GOL GOLWIDTH=80 GOLHEIGHT=60 go run ./cmd
```

**Run Battleship:**
```
MODULE=BATTLESHIP BATTLESHIPWIDTH=10 BATTLESHIPHEIGHT=10 go run ./cmd
```

## Project Structure
- `cmd/main.go`: Application entrypoint; reads the `MODULE` config and runs the selected game.
- `internal/config`: Configuration loading (env + .env support).
- `internal/ddd`: A generic 2D board implementation.
- `gameoflife/`: Contains the Game of Life module, including its specific board logic and Ebiten implementation.
- `battleship/`: Contains the Battleship module, including its board logic, AI, and Ebiten implementation.

## Development
- Run locally: `MODULE=<game> go run ./cmd`
- Build binary: `go build ./cmd`
- Formatting and linting: `go fmt`, `go vet`

## Roadmap / Ideas
- Add start/pause/reset hotkeys for Game of Life.
- Adjustable cell size and window scaling flags.
- Preset patterns (glider, pulsar, etc.) for Game of Life.
- More advanced AI for Battleship.
- Add a main menu to select games instead of using environment variables.

## License
Specify your license here (e.g., MIT).

## Acknowledgments
- Ebiten by Hajime Hoshi.
- `envconfig` and `stackus/dotenv` for configuration handling.
