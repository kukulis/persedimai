package env

import "github.com/joho/godotenv"

type EnvMap map[string]string

func (envFile EnvMap) Getenv(key string) string {

	// TODO padaryti, kad jeigu tuščia, tai kreiptis į os
	return envFile[key]
}

func (envFile EnvMap) Read(filename string) (envMap EnvMap, err error) {
	return godotenv.Read(filename)
}
