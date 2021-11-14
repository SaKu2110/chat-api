package config

import "os"

func GetAPIPort() string {
	if port := os.Getenv("API_PORT"); port != "" {
		return ":" + port
	}
	return ":8080"
}
