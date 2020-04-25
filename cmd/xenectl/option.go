package main

type xeneCtlConfig struct {
	APIServerAddr string `yaml:"apiServerAddr"`

	AuthToken string `yaml:"authToken"`
}

var config = &xeneCtlConfig{}
