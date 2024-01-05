package main

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"os"
	"text/template"
)

type InfoConfig struct {
	Template string `yaml:"template"`
}

type LogoConfig struct {
	Template string `yaml:"template"`
}

type Config struct {
	Info InfoConfig `yaml:"info"`
	Logo LogoConfig `yaml:"logo"`
}

func NewConfig(configPath string) *Config {
	configYaml, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	config := &Config{}
	err = yaml.Unmarshal(configYaml, config)
	return config
}

func (config *Config) getInfo(info Info) string {
	tmpl := template.Must(template.New("info").Parse(config.Info.Template))
	var b bytes.Buffer
	tmpl.Execute(&b, &info)

	return b.String()
}

func (config *Config) getLogo() string {
	return config.Logo.Template
}
