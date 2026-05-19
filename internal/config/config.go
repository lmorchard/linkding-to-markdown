package config

// Config holds application configuration. Most settings are read straight
// from viper at the call site; only the cross-cutting logging flags live on
// this struct.
type Config struct {
	Verbose bool
	Debug   bool
	LogJSON bool
}
