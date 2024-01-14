package main

func main() {
	c, colors := initialiseConfigAndColors()
	s := InitScreen()
	w, h := s.Size()
	g := NewGrid(w, h, c.Preset)

	NewGame(c).Run(s, g, colors)
}

func initialiseConfigAndColors() (*Config, *Colors) {
	c := NewConfig(ReadConfig())
	if c == nil {
		c = NewConfigWithDefaults()
		return c, DefaultColors()
	}

	return c, CustomColors(c)
}
