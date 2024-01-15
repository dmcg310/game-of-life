# game-of-life

### Running

Build
```sh
go build -o bin/game-of-life ./src
```

Run
```sh
./bin/game-of-life
```

or

```sh
go build -o bin/game-of-life ./src && ./bin/game-of-life
```

```sh
go install
```

this will install into $GOPATH/bin then run

```sh
game-of-life
```

```sh
game-of-life -h
```

### Patterns

case sensitive
- blinker
- toad
- beacon
- lwss
- gosper-glider-gun
- glider
- block
- random

### Colors

Following colors are supported:
- black
- maroon
- green
- olive
- navy
- purple
- teal
- silver
- gray
- red
- lime
- yellow
- blue
- fuchsia
- aqua
- white

### References

- [tcell](https://github.com/gdamore/tcell)
- [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life)
- [cli](https://github.com/urfave/cli)
