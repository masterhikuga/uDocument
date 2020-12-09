package Drivers

import "os"

type IWriter interface {
	Writing(files *os.File, records [][]string) error
}

type IReader interface {
	Reading(files *os.File) ([][]string, error)
}
