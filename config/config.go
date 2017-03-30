package config

const defaultPrompt = "_ "

type Config struct {
	Prompt         string
	StartUpCommand []byte
	Alias          [][2]string
}

func (c *Config) Init() {
	if c.Prompt == "" {
		c.Prompt = defaultPrompt
	}
}
