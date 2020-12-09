package Drivers

import (
	"encoding/csv"
	"os"
)

type CSV struct {
	IReader
	IWriter
}

func New_CSV() *CSV {
	return &CSV{}
}

func (doc *CSV) Reading(files *os.File) ([][]string, error) {

	r := csv.NewReader(files)

	records, err := r.ReadAll()

	return records, err
}

func (doc *CSV) Writing(files *os.File, records [][]string) error {
	w := csv.NewWriter(files)
	w.WriteAll(records) // calls Flush internally

	return w.Error()
}
