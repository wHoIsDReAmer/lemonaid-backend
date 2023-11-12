package dotenv

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func Load(name string) error {
	file, err := os.ReadFile(name)
	if err != nil {
		return errors.New("Cannot found environment file")
	}

	read := string(file)
	envs := strings.Split(read, "\n")

	for _, row := range envs {
		row = strings.ReplaceAll(row, "\r", "")
		splits := strings.Split(row, "=")
		if len(splits) == 1 {
			continue
		}

		os.Setenv(splits[0], strings.Join(splits[1:], "="))
	}

	fmt.Println("Load all environment in the file")

	return nil
}
