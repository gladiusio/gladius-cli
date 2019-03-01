package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// GetString - Wrapper around viper GetString
func GetString(key string) string {
	return viper.GetString(key)
}

// SetupConfig - Sets up, watches, and registers default config
func SetupConfig(configName string, defaults map[string]string) {
	viper.SetConfigName(configName)

	base, err := GetGladiusBase()
	if err != nil {
		viper.AddConfigPath(".") // Search only for local config
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath(base) // OS specifc
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	err = viper.ReadInConfig() // Find and read the config file
	// Should probably fix this...
	if err != nil {
		if strings.HasPrefix(err.Error(), "Config File") {
		} else { // Handle errors reading the config file
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	} else {
		viper.WatchConfig()
		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	}

}

// GetGladiusBase - Returns the base directory
func GetGladiusBase() (string, error) {
	var m string
	var err error

	if os.Getenv("GLADIUSBASE") != "" {
		m = os.Getenv("GLADIUSBASE")
	} else {
		switch runtime.GOOS {
		case "windows":
			m = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), ".gladius")
		case "linux":
			m = os.Getenv("HOME") + "/.gladius"
		case "darwin":
			m = os.Getenv("HOME") + "/.gladius"
		default:
			m = ""
			err = errors.New("unknown operating system, can't find gladius base directory. Set the GLADIUSBASE environment variable, or supply the directory as the first argument to add it manually")
		}
	}

	return m, err
}

func CLIDefaults() map[string]string {
	m := make(map[string]string)
	base, err := GetGladiusBase()
	if err != nil {
		log.Fatal("Could not retrieve gladius base")
	}
	m["DirLogs"] = filepath.Join(base, "logs")
	viper.SetDefault("Ports.Guardian", 7791)
	viper.SetDefault("Ports.EdgeD", 8081)
	viper.SetDefault("Ports.NetworkGateway", 3001)

	return m
}
