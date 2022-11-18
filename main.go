package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	occ, lines, err := FindReplaceFile("data.txt", "Python", "Go", "newdata.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("=== RESUME ===")
	defer fmt.Println("===  FIN  ===")

	fmt.Println("nombre d'occurence :", occ)
	fmt.Println("les lignes qui ont changées :", lines)
}

func ProcessLine(text, old, new string) (found bool, res string, count int) {
	oldLower := strings.ToLower(old)
	newLower := strings.ToLower(new)

	res = text

	if strings.Contains(text, old) || strings.Contains(text, oldLower) {
		found = true

		count += strings.Count(text, old)
		count += strings.Count(text, oldLower)

		res = strings.ReplaceAll(text, old, new)
		//must take res as first parameter to keep change
		res = strings.ReplaceAll(res, oldLower, newLower)
	}

	return found, res, count
}

func FindReplaceFile(src, old, new, dst string) (occ int, lines []int, err error) {
	file, err := os.Open(src)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	if err != nil {
		return occ, lines, err
	}

	newFile, err := os.Create(dst)
	defer func(newFile *os.File) {
		_ = newFile.Close()
	}(newFile)

	if err != nil {
		return occ, lines, err
	}

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(newFile)
	defer func(writer *bufio.Writer) {
		_ = writer.Flush()
	}(writer)

	line := 1
	for scanner.Scan() {
		found, res, count := ProcessLine(scanner.Text(), old, new)
		if found {
			occ += count
			lines = append(lines, line)
		}

		_, err = fmt.Fprintln(writer, res)
		if err != nil {
			return 0, nil, err
		}

		line++
	}

	return occ, lines, err

}
