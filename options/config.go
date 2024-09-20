package options

// Config manages default values for configuration settings and controls their mutability.
type Config struct {
	defaults map[string]any // Stores default values for configuration settings.
	sealed   bool           // Indicates if the config is immutable (sealed) or modifiable.
}

// NewConfig creates and returns a new Config instance with default values and is initially modifiable.
func NewConfig() *Config {
	return &Config{
		defaults: make(map[string]any),
		sealed:   false,
	}
}

// checkSealed panics if the config is sealed to prevent further modifications.
func (c *Config) checkSealed() {
	if c.sealed {
		panic("config is sealed and cannot be modified")
	}
}

// SetDefault sets a default value for a specified key if the config is not sealed.
// Panics if the config is sealed.
func (c *Config) SetDefault(k string, v any) {
	c.checkSealed()
	c.defaults[k] = v
}

// GetDefault retrieves the default value associated with the specified key.
// Returns nil if the key does not exist.
func (c *Config) GetDefault(k string) any {
	return c.defaults[k]
}

// Seal marks the config as immutable, preventing any further modifications.
func (c *Config) Seal() {
	c.sealed = true
}

// Global Config instance.
var (
	c = NewConfig()
)

// SetDefault updates the global config with a default value for a specified key.
// Panics if the global config is sealed.
func SetDefault(k string, v any) {
	c.SetDefault(k, v)
}

// GetDefault retrieves a default value from the global config for a specified key.
// Returns nil if the key does not exist.
func GetDefault(k string) any {
	return c.GetDefault(k)
}

// Seal marks the global config as immutable, preventing any further modifications.
func Seal() {
	c.Seal()
}
