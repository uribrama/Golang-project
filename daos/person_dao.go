package person

import (
	"strconv"

	"github.com/ansel1/merry"
	"github.com/uribrama/Golang-project/logger"
	"github.com/uribrama/Golang-project/models/database"
	person "github.com/uribrama/Golang-project/models/person"
	"gorm.io/gorm"
)

var (
	ErrPersonNotFound = merry.New("Person not found")
	ErrCreatingPerson = merry.New("Error during person creation")
	ErrUpdatingPerson = merry.New("Error updating person")
)

const tableName = "persons"

type PersonDAO interface {
	Create(person.Person) (person.Person, error)
	Update(person.Person) (person.Person, error)
	GetPersons(person.PersonSearchArgs) ([]person.Person, error)
}

type personDaoImpl struct {
	ctx   logger.Logger
	table string
	db    *database.Database
}

func PersonDaoInit(ctx logger.Logger, db *database.Database) PersonDAO {
	return &personDaoImpl{
		ctx:   ctx,
		table: tableName,
		db:    db,
	}
}

func (dao *personDaoImpl) Create(p person.Person) (person.Person, error) {
	tx := dao.db.Create(&p)
	if tx.Error != nil {
		return person.Person{}, ErrCreatingPerson.Here()
	}
	return p, nil
}

func (dao *personDaoImpl) Update(p person.Person) (person.Person, error) {
	tx := dao.db.Updates(&p)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return person.Person{}, ErrUpdatingPerson.Here().WithValue("id", p.ID)
	}

	dao.db.Find(&p)
	return p, nil
}

func (dao *personDaoImpl) GetPersons(search person.PersonSearchArgs) (persons []person.Person, err error) {
	var result *gorm.DB
	filter := createUserFilter(search)

	queryConditions := dao.buildQueryFilter(filter, search)
	result = queryConditions.Find(&persons)

	if result.RowsAffected == 0 || merry.Is(result.Error, gorm.ErrRecordNotFound) {
		return []person.Person{}, ErrPersonNotFound.Here().WithValue("args", search)
	}
	return persons, nil
}

func createUserFilter(searchArgs person.PersonSearchArgs) *person.Person {
	filterPerson := &person.Person{}
	if searchArgs.DNI != nil {
		filterPerson.DNI = searchArgs.DNI
	}
	if searchArgs.Email != "" {
		filterPerson.Email = &searchArgs.Email
	}
	return filterPerson
}

func (dao *personDaoImpl) buildQueryFilter(filter *person.Person, args person.PersonSearchArgs) *gorm.DB {
	query := dao.db.Where(filter)
	if args.FromSalary != nil {
		query.Where("salary >= " + strconv.Itoa(*args.FromSalary))
	}
	if args.ToSalary != nil {
		query.Where("salary <= " + strconv.Itoa(*args.ToSalary))
	}
	return query
}
