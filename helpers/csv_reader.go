package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCsv reads csv file and returns slice of string
func ReadCsv(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Csv file could not open.")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvSlice, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Csv file could not read.")
	}

	return csvSlice
}
