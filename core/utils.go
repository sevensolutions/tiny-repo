package core

import (
	"os"
	"strconv"
)

func GetRequiredEnvVar(name string) string {
	value := os.Getenv(name)

	if value == "" {
		panic("Missing value for required environment variable " + name)
	}

	return value
}

func GetRequiredEnvVarBool(name string) bool {
	value := GetRequiredEnvVar(name)

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		panic("Invalid value for boolean environment variable " + name)
	}

	return boolValue
}

func FilterArray[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func MapArray[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
