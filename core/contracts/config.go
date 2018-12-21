package contracts

// Config, main configuration data structure
type Config struct {
	HomeAssistant HomeAssistantConfig `yaml:"home_assistant"`
}

// HomeAssistantConfig is the configuration for the Home Assistant platform integration
type HomeAssistantConfig struct {
	IP    string `yaml:"ip"`
	SSL   bool   `yaml:"ssl"`
	Token string `yaml:"token"`
}
