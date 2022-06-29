package main

type KubeConfig struct {
	DryRun          bool          `mapstructure:"DryRun"`
	InCluster       bool          `mapstructure:"InCluster"`
	ConfigEntryList []ConfigEntry `mapstructure:"Jobs, mapstructure"`
}

type ConfigEntry struct {
	Name    string  `mapstructure:"Name"`
	GitItem GitItem `mapstructure:"Git"`
	PodItem PodItem `mapstructure:"Pod"`
}

type GitItem struct {
	GitOrigin      string `mapstructure:"Origin"`
	GitDestination string `mapstructure:"Destination"`
	GitUsername    string `mapstructure:"Username"`
	GitAccessToken string `mapstructure:"AccessToken"`
	GitEmail       string `mapstructure:"Email"`
	GitAuthor      string `mapstructure:"Author"`
}

type PodItem struct {
	PodNamespace     string `mapstructure:"Namespace"`
	PodLabel         string `mapstructure:"Label"`
	PodContainer     string `mapstructure:"Container"`
	PodContainerPath string `mapstructure:"ContainerPath"`
	UseFirstPod      bool   `mapstructure:"UseFirst"`
}
