package settings

import "github.com/spf13/viper"

// AddConfigItems adds a new configuration item, and makes it overridable by env vars
func AddConfigItems(configItems []string) {
	for _, item := range configItems {
		viper.BindEnv(item)
	}
}

// ApplyWith gets a setting from viper, and passes it to a closure
func ApplyWith(item string, f func(interface{})) {
	if viper.IsSet(item) {
		f(viper.Get(item))
	}
}
