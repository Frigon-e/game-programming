# SideProjectGames

An experimental Go project showcasing small game/prototyping modules. The current implemented module is Conway's Game of Life rendered with Ebiten (2D game library). It uses a small domain layer for a generic 2D board and a specialized Game of Life board, plus a simple configuration system driven by environment variables.

## Features
- Conway's Game of Life simulation using Ebiten for real‑time rendering.
- Interactive controls:
  - Left click: toggle the cell under the mouse.
  - Arrow Down: speed up simulation (decrease step time).
  - Arrow Up: slow down simulation (increase step time).
- Toroidal wrapping board (edges wrap around).
- Simple, environment‑variable driven configuration for board size.

## Requirements
- Go 1.25 or newer (as declared in go.mod).
- A desktop environment supported by Ebiten (macOS, Windows, or Linux with X11/Wayland). For Linux, you may need additional system packages (OpenGL/GLFW support) depending on your distro.

## Getting Started

### Clone and build
```
git clone <your-fork-or-repo-url> SideProjectGames
cd SideProjectGames
go build ./cmd
```
This will create an executable in the current directory (on some systems named `cmd` or `cmd.exe`). You can also run it directly with `go run`:
```
go run ./cmd
```

### Configuration
Configuration is read from environment variables using `envconfig`, optionally loading a `.env` file depending on the `ENVIRONMENT` variable via `stackus/dotenv`.

Supported variables:
- GOLWIDTH: integer, width of the Game of Life grid in cells.
- GOLHEIGHT: integer, height of the Game of Life grid in cells.
- ENVIRONMENT: optional, used by the dotenv loader to decide which .env files to read. See below.

How .env loading works (internal/config/config.go):
- On startup, the program calls `dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT")))`.
- If `ENVIRONMENT` is empty, dotenv typically attempts to load from default `.env` (and variants) in the project root if present.
- If you set `ENVIRONMENT=local` (for example), dotenv will prefer files like `.env.local` (depending on the library’s resolution rules).

Example `.env` file:
```
GOLWIDTH=80
GOLHEIGHT=60
```
You can also export variables directly in your shell before running.

### Run
```
GOLWIDTH=80 GOLHEIGHT=60 go run ./cmd
```
This launches a window sized to width*cellSize by height*cellSize pixels. The default cell size is 10 pixels and can be changed in code (gameoflife/module.go).

## Controls
- Left mouse button: Toggle the clicked cell alive/dead.
- Up Arrow: Increase step time (slows the simulation).
- Down Arrow: Decrease step time (speeds up the simulation).

## Project Structure
- cmd/main.go — application entrypoint; initializes config and runs the Game of Life module.
- internal/config — configuration loading (env + .env support).
- internal/ddd — a generic 2D board implementation using a flat slice with wrapping coordinates.
- gameoflife/internal/ddd — Game of Life specific board (bool grid) with helpers to seed and count neighbors.
- gameoflife/module.go — Ebiten game implementation: update loop, drawing, input handling, and rules.

## Development
- Run locally:
  - `go run ./cmd`
- Build binary:
  - `go build ./cmd`
- Formatting and linting: follow standard Go tooling (`go fmt`, `go vet`).

### Notes on the Board implementation
- Coordinates wrap around (toroidal world). Access outside the range wraps back into bounds.
- Board uses a flat slice; index math is handled by helper methods.

## Roadmap / Ideas
- Add start/pause/reset hotkeys.
- Adjustable cell size and window scaling flags.
- Preset patterns (glider, pulsar, etc.) and randomized seeding levels.
- Add a Battleship game module (grid-based placement, turn-based targeting, basic AI/opponent). 
- Additional mini‑modules beyond Game of Life.

## License
Specify your license here (e.g., MIT). If you add a LICENSE file, reference it.

## Acknowledgments
- Ebiten by Hajime Hoshi.
- envconfig and stackus/dotenv for configuration handling.
