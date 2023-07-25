package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Replacement struct {
	Original    string
	Replacement string
}

func readReplacementsFromCSV(filePath string) ([]Replacement, error) {
	var replacements []Replacement

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if len(record) >= 2 {
			replacement := Replacement{
				Original:    record[0],
				Replacement: record[1],
			}
			replacements = append(replacements, replacement)
		}
	}

	return replacements, nil
}

func replaceTextInFile(inputFilePath, outputFilePath string, replacements []Replacement) error {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()

	for scanner.Scan() {
		line := scanner.Text()
		for _, replacement := range replacements {
			line = strings.ReplaceAll(line, replacement.Original, "\""+replacement.Replacement+"\",")
		}
		fmt.Fprintln(writer, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func main() {
	replacementFilePath := "replacements.csv" // Replace with the path to your CSV file containing replacements
	inputFilePath := "input.txt"              // Replace with the path to your input .txt file
	outputFilePath := "output.txt"            // Replace with the desired path for the output .txt file

	replacements, err := readReplacementsFromCSV(replacementFilePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = replaceTextInFile(inputFilePath, outputFilePath, replacements)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Text replacement completed successfully.")
}
