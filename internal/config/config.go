package config

import "time"

// Config holds application configuration
type Config struct {
	// Core settings
	Database string
	Verbose  bool
	Debug    bool
	LogJSON  bool

	// Linkding settings
	Linkding struct {
		URL     string
		Token   string
		Timeout time.Duration
	}

	// Fetch command settings
	Fetch struct {
		Days          int
		Since         string
		Until         string
		Query         string
		Output        string
		Title         string
		NoNotes       bool
		NoTags        bool
		NoGroupByDate bool
		DateFormat    string
	}
}
