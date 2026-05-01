package config

func MustGetPort() string {
	return mustGetEnvString("PORT")
}

func MustGetRepoRoot() string {
	return mustGetEnvString("REPO_ROOT")
}

func MustGetServerName() string {
	return mustGetEnvString("SERVER_NAME")
}

func MustGetServerVersion() string {
	return mustGetEnvString("SERVER_VERSION")
}
