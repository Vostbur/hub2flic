package main

import (
	"log"
	"os"
)

func isPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func cleanUp(path string) {
	isClonePathExists, _ := isPathExists(path)
	if isClonePathExists {
		err := os.RemoveAll(path)
		if err != nil {
			log.Printf("\033[31;1m%s\033[0m\n", err)
		}
	}
}
