# game-of-life

A terminal implementation of Conway's Game of Life, a cellular automaton created by mathematician John Conway. This program transforms your terminal into a dynamic playground for exploring various patterns and behaviors within the game.

![Game-of-Life](game-of-life.gif)

## Features
- **Customizable Patterns**: Choose from a variety of predefined patterns like blinker, toad, beacon, and more.
- **Dynamic Color Schemes**: Personalize your experience with a selection of colors for cells and backgrounds.
- **Flexible Configuration**: Tailor the game's behavior with custom FPS settings and more through a configuration file or command-line arguments.
- **Responsive Terminal Display**: Watch the Game of Life unfold in real-time within your terminal.

## Installation
1. **Prerequisites**: Ensure you have Go installed on your system.
2. **Getting the Program**:
   - Clone the repository or download the source code.
   - Alternatively, you can directly install the program using Go:
     ```sh
     go install
     ```
     This command installs the program into `$GOPATH/bin`. After installation, you can run the program using `game-of-life` command in your terminal.

## Usage
1. **Running the Game**:
   - Directly execute the binary after [building](#build-from-source):
     ```sh
     ./bin/game-of-life
     ```
   - Or, if installed via `go install`, simply run:
     ```sh
     game-of-life
     ```
2. **Command-Line Arguments**:
   - Optionally, specify a pattern and FPS:
     ```sh
     game-of-life [pattern] [fps]
     ```
   - Use `-h` for help:
     ```sh
     game-of-life -h
     ```
3. **During the Game**:
   - Press 'p' to pause/resume the game at the current turn. Once paused you can use 'space' to step through one turn at a time.

## Configuration
- The program looks for `gol-config.json` in the current directory or the default configuration location based on your system (e.g., `/Users/<username>/Library/Application Support/gol/gol-config.json` for MacOS, and `/home/&lt;username&gt;/.config/gol/` for Linux file systems).
- A sample configuration file:
  ```js
  {
      "preset": "random",
      "cell-color": "gray",
      "background-color": "white",
      "scale-factor": 1,
      "fps": 23
  }
  ```
- To find out the configuration directory path, use:
```sh
game-of-life cl
```

## Available Patterns
The game supports various patterns which are case-sensitive:
- blinker, toad, beacon, lwss, gosper-glider-gun, glider, block, random

## Color Options
You can choose from a range of supported colors for cells and background:
- black, maroon, green, olive, navy, purple, teal, silver, gray, red, lime, yellow, blue, fuchsia, aqua, white

<a name="build-from-source"></a>
## Building from Source
- To compile the program:
```sh
go build -o bin/game-of-life
```
This command generates an executable in the bin directory.
- You can also build and run in one step:
```sh
go build -o bin/game-of-life && ./bin/game-of-life
```

## References and Acknowledgements

- **tcell Library**: This project leverages the [tcell library](https://github.com/gdamore/tcell) for managing terminal graphics and events.
- **Conway's Game of Life**: Learn more about the original concept on [Wikipedia](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life).
- **cli Package**: Command-line interactions are powered by the [cli package](https://github.com/urfave/cli).
