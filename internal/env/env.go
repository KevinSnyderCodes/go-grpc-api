package env

import (
	"os"
	"strconv"
)

func Has(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

func GetString(key string) (string, error) {
	return os.Getenv(key), nil
}

func MustGetString(key string) string {
	v, err := GetString(key)
	if err != nil {
		panic(err)
	}

	return v
}

func GetStringOrDefault(key string, value string) (string, error) {
	if !Has(key) {
		return value, nil
	}

	return GetString(key)
}

func MustGetStringOrDefault(key string, value string) string {
	v, err := GetStringOrDefault(key, value)
	if err != nil {
		panic(err)
	}

	return v
}

func GetUint(key string) (uint, error) {
	s := os.Getenv(key)

	v, err := strconv.ParseUint(s, 0, 0)

	return uint(v), err
}

func MustGetUint(key string) uint {
	v, err := GetUint(key)
	if err != nil {
		panic(err)
	}

	return v
}

func GetUintOrDefault(key string, value uint) (uint, error) {
	if !Has(key) {
		return value, nil
	}

	return GetUint(key)
}

func MustGetUintOrDefault(key string, value uint) uint {
	v, err := GetUintOrDefault(key, value)
	if err != nil {
		panic(err)
	}

	return v
}
