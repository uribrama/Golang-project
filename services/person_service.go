package services

import (
	"encoding/csv"
	"fmt"

	"github.com/ansel1/merry"
	"github.com/go-playground/validator/v10"
	"github.com/uribrama/Golang-project/config"
	csvc "github.com/uribrama/Golang-project/csv"
	persondao "github.com/uribrama/Golang-project/daos"
	"github.com/uribrama/Golang-project/logger"
	person "github.com/uribrama/Golang-project/models/person"
)

type PersonService interface {
	GetPersons(person.PersonSearchArgs) ([]person.Person, error)
	UpdatePerson(person.PersonUpdateRequest) (person.Person, error)
	CreatePerson(person.PersonRequest) (person.Person, error)
	BatchPersons(csvc.FileUpload) ([]person.Person, error)
	InsertCsvPersons([]person.CSVPerson) []person.Person
}

type personServiceImpl struct {
	ctx       logger.Logger
	personDao persondao.PersonDAO
	cfg       config.Config
}

var (
	FailedToUpdate    = merry.New("Failed to update person")
	FailedToCreate    = merry.New("Failed to create person")
	ErrPersonNotFound = merry.New("Person/s not found")
	ErrReadingFile    = merry.New("Error reading file")
)

func PersonServiceInit(
	ctx logger.Logger,
	cfg config.Config,
	pd persondao.PersonDAO,
) PersonService {
	return &personServiceImpl{
		ctx:       ctx,
		cfg:       cfg,
		personDao: pd,
	}
}

func (s *personServiceImpl) GetPersons(search person.PersonSearchArgs) ([]person.Person, error) {
	s.ctx.Info("Getting person by: ", search)
	p, err := s.personDao.GetPersons(search)
	if err != nil && merry.Is(err, persondao.ErrPersonNotFound) {
		s.ctx.Error(err)
		return []person.Person{}, ErrPersonNotFound.Here()
	}
	return p, nil
}

func (s *personServiceImpl) UpdatePerson(p person.PersonUpdateRequest) (personResponse person.Person, err error) {
	if err != nil {
		return person.Person{}, err
	}

	s.ctx.Info("Fields receive to update person: ", p)

	per := person.FromPersonUpdateRequest(p)
	personResponse, err = s.personDao.Update(per)
	if err != nil && merry.Is(err, persondao.ErrUpdatingPerson) {
		s.ctx.Error(err)
		return person.Person{}, FailedToUpdate.Here()
	}
	return personResponse, err
}

func (s *personServiceImpl) CreatePerson(p person.PersonRequest) (person.Person, error) {
	s.ctx.Info("Fields receive to create person: ", p)

	per := person.FromPersonRequest(p)
	per, err := s.personDao.Create(per)
	if err != nil && merry.Is(err, persondao.ErrCreatingPerson) {
		s.ctx.Error(err)
		return person.Person{}, FailedToCreate.Here()
	}
	return per, err
}

func (s *personServiceImpl) BatchPersons(file csvc.FileUpload) ([]person.Person, error) {
	s.ctx.Info("Received a batch csv, starting validation")
	err := csvc.Validate(file, 10000)
	if err != nil {
		return nil, err
	}

	fileOpen, err := file.CSVFile.Open()
	if err != nil {
		s.ctx.Error(err)
		return nil, ErrReadingFile.Here()

	}
	defer fileOpen.Close()

	records, err := csv.NewReader(fileOpen).ReadAll()
	if err != nil {
		s.ctx.Error(err)
		return nil, ErrReadingFile.Here()
	}

	csvPersons := csvc.DecodeCsvFields(records)
	persons := s.InsertCsvPersons(csvPersons)

	s.ctx.Info(fmt.Sprintf("Inserted %d persons", len(persons)))
	return persons, nil
}

func (s *personServiceImpl) InsertCsvPersons(csvPersons []person.CSVPerson) []person.Person {
	var persons []person.Person
	for _, per := range csvPersons {
		validate := validator.New()
		err := validate.Struct(per)
		if err != nil {
			s.ctx.Debug(err.Error()+". person: ", per)
			continue
		}
		req := per.ToPersonRequest()
		//go routine here
		person, err := s.CreatePerson(req)
		if err != nil {
			s.ctx.Error(err)
		}
		persons = append(persons, person)
	}
	return persons
}
