package options

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Version string = "v0.0.1"

type LogOptions struct {
	Level string `yaml:"level"`
	Path  string `yaml:"path"`
}

type WebOptions struct {
	Listen string `yaml:"listen"`
}

type ScrapeOptions struct {
	Interval string `yaml:"interval"`
}

type CoreOptions struct {
	Log LogOptions `yaml:"log"`
}
type Options struct {
	Core   CoreOptions   `yaml:"core"`
	Web    WebOptions    `yaml:"web"`
	Scrape ScrapeOptions `yaml:"scrape"`
}

func NewOptions() (opts Options) {
	optsSource := viper.AllSettings()
	err := createOptions(optsSource, &opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "create options failed:", err)
		os.Exit(1)
	}
	return
}
