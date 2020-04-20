package option

// XeneConfig is the global config structure for xene.
type XeneConfig struct {
	// APIServer contains the options for api server
	APIServer APIServerOpts `yaml:"apiserver"`

	// Agent contains configuration for xene agent.
	Agent AgentOpts `yaml:"agent"`

	// Store is the store configuration that xene is configured with.
	Store Store `yaml:"store"`
}

var (
	// Config contains the global config for xene
	Config = &XeneConfig{}

	// ConfigFile is the path of the config file to use for Xene
	// it will override the default config file.
	ConfigFile string
)
