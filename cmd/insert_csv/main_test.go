package main

import (
	"os"
	"path"

	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/uribrama/Golang-project/config"
	dao "github.com/uribrama/Golang-project/daos"
	"github.com/uribrama/Golang-project/logger"
	"github.com/uribrama/Golang-project/models/database"
	"github.com/uribrama/Golang-project/models/person"
)

var _ = Describe("Test cmd insert data functions", func() {
	var (
		ctx = *logger.New(true)
	)

	Context("Test read persons from file", func() {
		It("Read csv file correctly", func() {
			currentDir, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			persons, _ := readAndParseFile(ctx, path.Join(currentDir, "/test.csv"))
			Ω(len(persons)).Should(Equal(5))
		})

		It("Read csv file fails", func() {
			persons, err := readAndParseFile(ctx, "non.csv")
			Ω(err).ShouldNot(BeNil())
			Ω(err.Error()).Should(Equal("no such file or directory"))
			Ω(len(persons)).Should(Equal(0))

		})
	})

	Context("Test insert persons", func() {
		var persons []person.CSVPerson
		var db mockDB

		BeforeEach(func() {
			persons, _ = readAndParseFile(ctx, "testdata/test.csv")
			db = mockDB{}

		})

		It("Insert persons correctly", func() {
			var personDao = dao.PersonDaoInit(ctx, &db.Database)
			perReq := person.FromPersonRequest(persons[0].ToPersonRequest())
			_, err := personDao.Create(perReq)
			Ω(err).Should(BeNil())
		})
	})
})

type mockDB struct {
	database.Database
	mock.Mock
}

func (d *mockDB) Start(ctx logger.Logger, cfg config.Config) error {
	return d.Called(ctx, cfg).Error(0)
}
