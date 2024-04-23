package packages

import (
	"fmt"
	"os"
)

func WriteToFile(filename string, data []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range data {
		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}
