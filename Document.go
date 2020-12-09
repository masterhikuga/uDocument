package Document

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strconv"

	"github.com/mashingan/smapping"
)

const (
	defaultTag = "Doc"
)

type Document struct {
	format   IDocument
	namefile string
}

/*
===========================
	Public
===========================
*/

func New_Document(Format IDocument, namefile string) *Document {
	doc := Document{
		format:   Format,
		namefile: namefile,
	}

	return &doc
}

func (doc *Document) ReadToMAP() ([]string, []map[string]interface{}) {
	files, err := os.Open(doc.namefile)
	if err != nil {
		log.Fatal(err)
	}

	records, err := doc.format.Reading(files)
	if err != nil {
		log.Fatal(err)
	}

	if err := files.Close(); err != nil {
		log.Fatal(err)
	}

	heads, datas := doc.convTwoDarrayToMAP(records)

	return heads, datas
}

func (doc *Document) WriteFromMAP(heads []string, datas []map[string]interface{}) {
	files, err := os.Create(doc.namefile)
	if err != nil {
		log.Fatal(err)
	}

	err = doc.format.Writing(files, doc.convMAPToTwoDarray(datas, heads))

	if err != nil {
		log.Fatal(err)
	}

	if err := files.Close(); err != nil {
		log.Fatal(err)
	}
}

func (doc *Document) Reads(Data interface{}) (*Document) {
	_, records := doc.ReadToMAP()
	err := doc.convMAPToInterface(records, Data)
	if err != nil {
		panic(err)
	}
	return doc
}

func (doc *Document) Writes(Data interface{}) (*Document) {
	err, Head, Values := doc.convInterfaceToMAP(Data)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	doc.WriteFromMAP(Head, Values)
	return doc
}

/*
===========================
	Private
===========================
*/

func (doc *Document) convMAPToTwoDarray(records []map[string]interface{}, heads []string) [][]string {
	var WriteRecords [][]string

	WriteRecords = append(WriteRecords, heads)
	for _, record := range records {
		var valueField []string

		for _, name := range heads {
			valueField = append(valueField, doc.convAllInString(record[name]))
		}

		WriteRecords = append(WriteRecords, valueField)
	}

	return WriteRecords
}

func (doc *Document) convTwoDarrayToMAP(WriteRecords [][]string) (heads []string, records []map[string]interface{}) {
	for index, recLines := range WriteRecords {
		if index == 0 {
			heads = recLines
			continue
		}

		temp := make(map[string]interface{})
		for index, recValues := range recLines {
			temp[heads[index]] = doc.convToTypeData(recValues)
		}

		records = append(records, temp)
	}

	return
}

func (doc *Document) convAllInString(value interface{}) string {
	conv, check := value.(string)
	if check {
		return conv
	}

	var convInt int
	convInt, check = value.(int)
	if check {
		conv = strconv.Itoa(convInt)
		return conv
	}

	var convFloat float64
	convFloat, check = value.(float64)
	if check {
		conv = strconv.FormatFloat(convFloat, 'f', -1, 64)
		return conv
	}

	var convBool bool
	convBool, check = value.(bool)
	if check {
		conv = strconv.FormatBool(convBool)
		return conv
	}
	return conv
}

func (doc *Document) convToTypeData(raw string) interface{} {
	if process, err := strconv.Atoi(raw); err == nil {
		return process
	}

	if process, err := strconv.ParseBool(raw); err == nil {
		return process
	}

	if process, err := strconv.ParseFloat(raw, 64); err == nil {
		return process
	}

	if process, err := strconv.ParseUint(raw, 10, 64); err == nil {
		return process
	}

	return raw
}

func (doc *Document) convInterfaceToMAP(data interface{}) (err error, Head []string, Value []map[string]interface{}) {
	Value = nil
	DataValue := reflect.ValueOf(data)
	DataValue = reflect.Indirect(DataValue)

	if DataValue.Kind() != reflect.Slice {
		err = errors.New("Must Slice or array")
		return
	}

	len := DataValue.Len()
	for index := 0; index < len; index++ {
		var MapValues map[string]interface{}
		MainStruct := reflect.New(DataValue.Type().Elem()).Interface()
		MainStructValue := reflect.ValueOf(MainStruct)
		MainStructValue = reflect.Indirect(MainStructValue)
		MainStructValue.Set(DataValue.Index(index))

		MapValues = smapping.MapFields(MainStruct)
		Value = append(Value, MapValues)

		if index == 0 {
			DataType := DataValue.Index(index).Type()
			for numfield := 0; numfield < DataValue.Index(index).NumField(); numfield++ {
				Head = append(Head, DataType.Field(numfield).Name)
			}
		}
	}
	return
}

func (doc *Document) convMAPToInterface(Datas []map[string]interface{}, records interface{}) error {
	var err error

	recordsValue := reflect.ValueOf(records)
	recordsValue = reflect.Indirect(recordsValue)

	if !recordsValue.CanSet() {
		return errors.New("Must Can Set")
	}

	if recordsValue.Kind() != reflect.Slice {
		return errors.New("Must be a slice")
	}

	recordsDump := reflect.New(recordsValue.Type())
	recordsDump = reflect.Indirect(recordsDump)

	for _, data := range Datas {
		mainstruct := reflect.New(recordsValue.Type().Elem())
		rec := mainstruct.Interface()
		err = smapping.FillStruct(rec, smapping.Mapped(data))
		mainstruct = reflect.Indirect(reflect.ValueOf(rec))
		recordsDump = reflect.Append(recordsDump, mainstruct)
	}

	recordsValue.Set(recordsDump)
	return err
}
