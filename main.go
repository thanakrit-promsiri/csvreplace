package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Replacement struct {
	Original    string
	Replacement string
}

type DictionaryWord struct {
	Word      string
	Translate string
	Length    int
}

func sortCsv(csvTxt [][]string) []DictionaryWord {

	var DictionaryWordList []DictionaryWord
	for _, record := range csvTxt {
		if len(record) >= 2 {
			dictionaryWord := DictionaryWord{
				Word:      strings.TrimSpace(record[0]),
				Translate: strings.TrimSpace(record[1]),
				Length:    utf8.RuneCountInString(strings.TrimSpace(record[0])),
			}
			DictionaryWordList = append(DictionaryWordList, dictionaryWord)
		}
	}

	sort.SliceStable(DictionaryWordList, func(i, j int) bool {
		return DictionaryWordList[i].Length > DictionaryWordList[j].Length
	})

	return DictionaryWordList
}

func readReplacementsFromCSV(filePath string) ([]Replacement, error) {

	var replacements []Replacement

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	} else {
		defer file.Close()
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	DictionaryWordList := sortCsv(records)

	CsvNewWriter("replacementSort.csv", DictionaryWordList)

	file2, err := os.Open("replacementSort.csv")
	if err != nil {
		return nil, err
	} else {
		defer file.Close()
	}

	reader2 := csv.NewReader(file2)
	records2, err := reader2.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records2 {
		if len(record) >= 2 {
			replacement := Replacement{
				Original:    strings.TrimSpace(record[0]),
				Replacement: strings.TrimSpace(record[1]),
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

func CsvNewWriter(path string, DictionaryWordList []DictionaryWord) {
	csvFile, err := os.Create(path)
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	csvFile.Write(bomUtf8)
	csvWriter := csv.NewWriter(csvFile)

	defer csvFile.Close()
	r := make([]string, 0, 3)
	r = append(r, "Word")
	r = append(r, "Translate")
	r = append(r, "Length")
	csvWriter.Write(r)

	for _, DictionaryWord := range DictionaryWordList {
		r := make([]string, 0, 3)

		r = append(r, DictionaryWord.Word)
		r = append(r, DictionaryWord.Translate)
		r = append(r, strconv.Itoa(DictionaryWord.Length))
		err := csvWriter.Write(r)
		if err != nil {
			panic(err)
		}
		csvWriter.Flush()
	}

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
