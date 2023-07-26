package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"os"
	"strings"

	"github.com/ansel1/merry"
	"github.com/uribrama/Golang-project/config"
	csvc "github.com/uribrama/Golang-project/csv"
	dao "github.com/uribrama/Golang-project/daos"
	"github.com/uribrama/Golang-project/logger"
	"github.com/uribrama/Golang-project/models/database"
	"github.com/uribrama/Golang-project/models/person"
	"github.com/uribrama/Golang-project/services"
)

var (
	cfg      config.Config
	ctx      logger.Logger
	pathFile = flag.String("file", "", "path to read json data")
)

func main() {
	os.Setenv("PROJECT", "Golang-project")
	os.Setenv("GO_ENV", "development")
	cfg = config.Instance()
	ctx = *logger.New(config.Instance().Get().Debug)

	flag.Parse()

	if *pathFile == "" {
		ctx.Error(merry.New("Please specify -file argument"))
		return
	}

	db := setup()

	if !strings.Contains(*pathFile, ".csv") {
		ctx.Error(merry.New("File " + *pathFile + " is not csv extension"))
		return
	}

	personDao := dao.PersonDaoInit(ctx, &db)
	personService := services.PersonServiceInit(ctx, cfg, personDao)

	persons, err := readAndParseFile(ctx, *pathFile)
	if err != nil {
		ctx.ErrorL("CSV data could not be read", err)
		return
	}
	personService.InsertCsvPersons(persons)
}

func readAndParseFile(ctx logger.Logger, filePath string) ([]person.CSVPerson, error) {
	if _, err := os.Stat(*pathFile); errors.Is(err, os.ErrNotExist) {
		ctx.Error(err)
		return nil, err
	}
	file, err := os.Open(*pathFile)
	if err != nil {
		ctx.Error(err)
		return nil, err

	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, merry.New("Error reading file").Here()
	}

	csvPersons := csvc.DecodeCsvFields(lines)
	return csvPersons, nil
}

func setup() database.Database {
	ctx.Info("Starting connection to db")

	db := database.Database{}
	_, err := db.Start(ctx, cfg)
	if err != nil {
		ctx.Error(err)
		panic(err)
	}
	return db
}
