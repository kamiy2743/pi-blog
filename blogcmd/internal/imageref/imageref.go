package imageref

import "fmt"

func For(serviceName string, env string) string {
	return fmt.Sprintf("ghcr.io/kamiy2743/blog/%s:%s", serviceName, env)
}

func JobImages(env string) []string {
	return []string{
		For("migration", env),
		For("seed", env),
	}
}

func IsImageEnv(env string) bool {
	return env == "stg" || env == "prd"
}
