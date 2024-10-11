package main

import (
	"os"
)

func main() {
	fs, err := os.ReadDir("../Tracks")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(fs); i++ {

		if fs[i].IsDir() {
			path := fs[i].Name()
			if IsExist(path + "/Illustration.png") {
				err := os.Rename(path+"/Illustration.png", "./"+path+".png")
				if err != nil {
					panic(err)
				}
			}
		}

	}
}
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
