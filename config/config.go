package config

import "text/template"

const defaultPrompt = "_ "

type Config struct {
	Prompt         string
	PromptTmpl     *template.Template
	StartUpCommand []byte
	Alias          [][2]string
}

func (c *Config) Init() {
	if c.Prompt == "" {
		c.Prompt = defaultPrompt
	}
}

type Info struct {
	WD string
}
