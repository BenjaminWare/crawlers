package shared_crawler_utils

import "os"

func GetFilenamesInDirectory(dir string) []string {
	file, err := os.Open(dir)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileNames, _ := file.Readdirnames(0)
	return fileNames
}