package apiserver

import "github.com/polyk005/Notesync1.0/internal/app/blocknot"

// Config
type Config struct {
	BinAddr  string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Blocknot *blocknot.Config
}

func NewConfig() *Config {
	return &Config{
		BinAddr:  ":8080",
		LogLevel: "debug",
		Blocknot: blocknot.NewConfig(),
	}
}
