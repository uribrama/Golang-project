package csv

import (
	"encoding/json"
	"mime/multipart"
	"strconv"

	"github.com/ansel1/merry"
	"github.com/uribrama/Golang-project/models/person"
)

type FileUpload struct {
	CSVFile *multipart.FileHeader `form:"file" binding:"required"`
}

const csvHeader = "text/csv"

var ErrIncorrectCsvFormat = merry.New("File format is not csv")
var ErrFileMaxSize = merry.New("File exceeds the limit")

func Validate(file FileUpload, maxSize int64) error {
	fileType := file.CSVFile.Header.Get("Content-Type")
	if fileType != csvHeader {
		return ErrIncorrectCsvFormat.Here().WithValue("format", fileType)
	}

	if file.CSVFile.Size > maxSize {
		return ErrFileMaxSize.Here().Append(strconv.FormatInt(maxSize, 10))
	}

	return nil
}

func DecodeCsvFields(records [][]string) []person.CSVPerson {
	var persons []person.CSVPerson
	//skip first line to read, as they have the headers
	for i := 1; i < len(records); i++ {
		m := make(map[string]string)
		recordLine := records[i]
		for u := 0; u < len(recordLine); u++ {
			m[records[0][u]] = recordLine[u]
		}
		per := decodeMapToStruct(m)
		persons = append(persons, per)
	}
	return persons
}

func decodeMapToStruct(mapValues map[string]string) person.CSVPerson {
	jsonString, _ := json.Marshal(mapValues)
	per := person.CSVPerson{}
	json.Unmarshal(jsonString, &per)
	return per
}
