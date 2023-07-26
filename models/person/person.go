package person

import (
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Person struct {
	gorm.Model
	FirstName *string
	LastName  *string
	DNI       *int `gorm:"unique;not null"`
	Email     *string
	Age       *int8
	Salary    *int
	CreatedAt time.Time `gorm:"type:timestamp"`
	UpdatedAt time.Time `gorm:"type:timestamp"`
}

type CSVPerson struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required" `
	DNI       string `json:"dni" validate:"required,gte=1"`
	Age       string `json:"age"  validate:"required,gte=1"`
	Email     string `json:"email" validate:"required,email"`
	Salary    string `json:"salary" validate:"omitempty,numeric,gte=1"`
}

type PersonRequest struct {
	Id        uint
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	DNI       int    `json:"DNI" binding:"required,gte=1"`
	Age       int8   `json:"age" binding:"required,numeric,gte=18"`
	Email     string `json:"email" binding:"required,email"`
	Salary    *int   `json:"salary" binding:"omitempty,numeric,gte=1"`
}

type PersonUpdateRequest struct {
	Id        uint
	FirstName *string
	LastName  *string
	Age       *int8   `json:"age" binding:"omitempty,numeric,gte=18"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Salary    *int    `json:"salary" binding:"omitempty,numeric,gte=1"`
}

type PersonSearchArgs struct {
	DNI        *int   `form:"dni"`
	Email      string `form:"email"`
	FromSalary *int   `form:"fromSalary"`
	ToSalary   *int   `form:"toSalary"`
}

type PersonResponse struct {
	Id        int64   `json:"id"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	DNI       *int    `json:"DNI"`
	Age       *int8   `json:"age"`
	Email     *string `json:"email"`
	Salary    *int    `json:"salary"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func FromPersonUpdateRequest(p PersonUpdateRequest) Person {

	person := Person{
		FirstName: p.FirstName,
		LastName:  p.LastName,
		Age:       p.Age,
		Email:     p.Email,
		Salary:    p.Salary,
	}
	person.ID = p.Id
	return person
}

func FromPersonRequest(p PersonRequest) Person {
	var (
		firstname string
		lastname  string
		dni       int
		age       int8
		email     string
	)

	firstname = p.FirstName
	lastname = p.LastName
	dni = p.DNI
	age = p.Age
	email = p.Email
	person := Person{
		FirstName: &firstname,
		LastName:  &lastname,
		DNI:       &dni,
		Age:       &age,
		Email:     &email,
		Salary:    p.Salary,
	}
	person.ID = p.Id
	return person
}

func (p Person) ToPersonResponse() PersonResponse {
	var updatedAt string
	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	response := PersonResponse{
		Id:        int64(p.ID),
		FirstName: p.FirstName,
		LastName:  p.LastName,
		DNI:       p.DNI,
		Age:       p.Age,
		Email:     p.Email,
		Salary:    p.Salary,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: updatedAt,
	}
	return response
}

func ToPersonsResponse(persons []Person) []PersonResponse {
	var personsResponse []PersonResponse
	for _, p := range persons {
		personsResponse = append(personsResponse, p.ToPersonResponse())
	}
	return personsResponse
}

func (csvPer CSVPerson) ToPersonRequest() PersonRequest {
	var (
		age    int8
		dni    int
		salary *int
	)

	a, _ := strconv.ParseInt(csvPer.Age, 10, 64)
	age = int8(a)

	d, _ := strconv.ParseInt(csvPer.DNI, 10, 64)
	dni = int(d)

	if csvPer.Salary != "" {
		s, _ := strconv.ParseInt(csvPer.Salary, 10, 64)
		parse := int(s)
		salary = &parse
	}

	perRequest := PersonRequest{FirstName: csvPer.FirstName, LastName: csvPer.LastName, DNI: dni, Age: age, Email: csvPer.Email, Salary: salary}
	return perRequest
}
