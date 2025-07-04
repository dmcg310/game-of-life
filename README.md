# Game of Life

A terminal implementation of Conway's Game of Life written in Go. This cellular automaton transforms your terminal into a dynamic playground for exploring various patterns and behaviors.

![Game-of-Life](game-of-life.gif)

## Features

- **Multiple Patterns**: Choose from predefined patterns like blinker, toad, beacon, glider, and more
- **Customizable Colors**: Configure cell and background colors through JSON configuration
- **Flexible Settings**: Adjust FPS and other behaviors via config file or command-line arguments
- **Interactive Controls**: Pause, resume, and step through generations manually

## Installation

### Prerequisites
- Go 1.21 or later

### Install from Source
```bash
git clone https://github.com/dmcg310/game-of-life.git
cd game-of-life
make install
```

### Quick Build
```bash
make build
```

## Usage

### Basic Usage
```bash
# Run with default settings
make run

# Run with specific pattern and FPS
make run-pattern PATTERN=glider FPS=30

# Or use the binary directly
./bin/game-of-life [pattern] [fps]
```

### Available Commands
```bash
# Show available patterns
make patterns

# Show configuration directory
make config-location

# Show help
make help
```

### Controls
- **p**: Pause/resume the simulation
- **Space**: Step through one generation when paused
- **q/Esc/Ctrl+C**: Quit the program

## Configuration

The program looks for `gol-config.json` in the current directory or your system's config directory:

```json
{
    "preset": "random",
    "cell-color": "white",
    "background-color": "black",
    "fps": 23
}
```

### Available Patterns
- `blinker`, `toad`, `beacon`, `lwss`, `gosper-glider-gun`, `glider`, `block`, `random`

### Supported Colors
- `black`, `maroon`, `green`, `olive`, `navy`, `purple`, `teal`, `silver`
- `gray`, `red`, `lime`, `yellow`, `blue`, `fuchsia`, `aqua`, `white`

## Development

### Build Commands
```bash
make all          # Clean, format, vet, and build
make build        # Build the application
make dev-build    # Build with race detection
make clean        # Clean build artifacts
```

### Code Quality
```bash
make fmt          # Format code
make vet          # Vet code
make tidy         # Tidy dependencies
```

### Cross-Platform Builds
```bash
make build-all    # Build for Linux, Windows, and macOS
make release      # Create release archives
```

## Project Structure

```
game-of-life/
├── cmd/                    # Application entry point
├── internal/
│   ├── cli/               # Command-line interface
│   ├── config/            # Configuration management
│   ├── display/           # Terminal display handling
│   ├── errors/            # Error handling
│   └── game/              # Core game logic and grid
├── Makefile               # Build automation
└── gol-config.json        # Default configuration
```

## References

- [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) - Original concept
- [tcell](https://github.com/gdamore/tcell) - Terminal interface library
- [cli](https://github.com/urfave/cli) - Command-line interface framework
