package main

const version = "1.0.0.0"

// Category represents a category with an identifier, name, and optional parent category.
type Category struct {
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	ParentID     int    `json:"parent_id"`
}

// Stream represents a media stream with associated metadata.
type Stream struct {
	Num          int    `json:"num"`
	Name         string `json:"name"`
	StreamType   string `json:"stream_type"`
	StreamID     int    `json:"stream_id"`
	StreamIcon   string `json:"stream_icon"`
	EpgChannelID string `json:"epg_channel_id"`
	CategoryID   string `json:"category_id"`
}

// Config represents the application configuration.
type Config struct {
	Mode    string        `json:"mode"`
	Server  ServerConfig  `json:"server"`
	Xtream  XtreamConfig  `json:"xtream"`
	Output  OutputConfig  `json:"output"`
	Logging LoggingConfig `json:"logging"`
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Port            int `json:"port"`
	Host            string `json:"host"`
	ReadTimeoutSec  int `json:"read_timeout_seconds"`
	WriteTimeoutSec int `json:"write_timeout_seconds"`
}

// XtreamConfig holds Xtream API configuration.
type XtreamConfig struct {
	BaseURL           string `json:"base_url"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	RequestTimeoutSec int    `json:"request_timeout_seconds"`
}

// OutputConfig holds output file configuration.
type OutputConfig struct {
	DefaultFilename string `json:"default_filename"`
}

// LoggingConfig holds logging configuration.
type LoggingConfig struct {
	Level string `json:"level"` // DEBUG, INFO, WARN, ERROR
}
