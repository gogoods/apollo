package apollo

type NamespaceConfig struct {
	conf      *Config
	Namespace string
}

func (config *Config) GlobalSettings() *NamespaceConfig {
	return config.GetNamespace("westudy.global.settings")
}

func (config *Config) GetNamespace(ns string) *NamespaceConfig {
	return &NamespaceConfig{
		conf:      config,
		Namespace: ns,
	}
}

func (nsConfig *NamespaceConfig) GetString(key string, defaultValue string) string {
	return nsConfig.conf.GetStringByNameSpace(nsConfig.Namespace, key, defaultValue)
}
