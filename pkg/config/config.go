package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config structure les paramètres de l'application
type Config struct {
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

// Load charge la configuration depuis un fichier et/ou l'environnement
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Configuration par défaut
	v.SetDefault("url", "https://127.0.0.1/api")

	// Configuration du fichier
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}

	// Mapping des variables d'environnement
	v.SetEnvPrefix("REPOFLOW")
	v.AutomaticEnv()
	// Permet de mapper REPOFLOW_URL vers la clé "url"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Lecture
	if err := v.ReadInConfig(); err != nil {
		// On n'ignore l'erreur que si le fichier est manquant (les env vars peuvent suffire)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
