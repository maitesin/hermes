package sql

import "fmt"

// Config defines the SQL connection configuration
type Config struct {
	URL          string
	SSLMode      string
	BinaryParams string
}

// DatabaseURL returns the url prepared with its param values.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("%s?sslmode=%s&binary_parameters=%s", c.URL, c.SSLMode, c.BinaryParams)
}
