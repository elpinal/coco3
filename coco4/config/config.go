package config

const defaultPrompt = "_ "

type Config struct {
	Prompt string
}

func (c *Config) Init() {
	if c.Prompt == "" {
		c.Prompt = defaultPrompt
	}
}
