package utils

import (
	"GoRestify/pkg/pkg_types"
	"log"
	"os"
)

// SetENVs set env value to env map
func SetENVs(envList []pkg_types.Envkey) pkg_types.Envs {

	envs := make(pkg_types.Envs)

	for _, v := range envList {
		envs[v] = setEnvValue(string(v), true)
	}

	return envs
}

func setEnvValue(envName string, isRequired bool) (envValue string) {
	envValue = os.Getenv(envName)
	if envValue == "" && isRequired {
		log.Fatalf("env required: %v", envName)
	}
	return
}
