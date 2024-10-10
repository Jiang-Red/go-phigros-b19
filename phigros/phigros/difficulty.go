package phigros

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

var difficulty = map[string][]float32{}

func LoadDifficult(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), "\t")
		diff := []float32{}
		for _, v := range lines[1:] {
			floatVal, _ := strconv.ParseFloat(v, 32)
			diff = append(diff, float32(floatVal))
		}
		difficulty[lines[0]]=diff
	}
	return scanner.Err()

}
