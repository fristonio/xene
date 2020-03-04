package option

// XeneConfig is the global config structure for xene.
type XeneConfig struct {
	// APIServer contains the options for api server
	APIServer APIServerOpts `json:"apiserver"`
}

var (
	// Config contains the global config for xene
	Config = &XeneConfig{}

	// ConfigFile is the path of the config file to use for Xene
	// it will override the default config file.
	ConfigFile string
)
