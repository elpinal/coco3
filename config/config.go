package config

import (
	"os"
	"path/filepath"
	"text/template"
)

const defaultPrompt = "_ "

var defaultHistFile = filepath.Join(os.Getenv("HOME"), ".coco3_history")

type Config struct {
	Prompt         string
	PromptTmpl     *template.Template
	StartUpCommand []byte
	Alias          [][2]string
	HistFile       string
}

func (c *Config) Init() {
	if c.Prompt == "" {
		c.Prompt = defaultPrompt
	}
	if c.HistFile == "" {
		c.HistFile = defaultHistFile
	}
}

type Info struct {
	WD string
}
