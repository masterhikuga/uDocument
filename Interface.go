package Document

import "Document/Document/Drivers"

type IDocument interface {
	Drivers.IReader
	Drivers.IWriter
}
