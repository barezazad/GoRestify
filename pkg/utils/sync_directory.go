package utils

import "os"

// SyncDirectories to check basic directories is exist. if they are not exist,it will be create.
func SyncDirectories(directories []string) {

	for _, v := range directories {
		isExistDir(v)
	}
}

func isExistDir(dirPath string) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}
}
