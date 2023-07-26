package main

import (
	"net/http"
	"strconv"

	"github.com/ansel1/merry"
	"github.com/gin-gonic/gin"
	"github.com/uribrama/Golang-project/config"
	"github.com/uribrama/Golang-project/csv"
	persondao "github.com/uribrama/Golang-project/daos"
	"github.com/uribrama/Golang-project/logger"
	"github.com/uribrama/Golang-project/models/database"
	person "github.com/uribrama/Golang-project/models/person"
	"github.com/uribrama/Golang-project/services"
	"golang.org/x/net/context"
)

type PersonRequestHandler interface {
	healthcheck(c *gin.Context)
	configRequest(c *gin.Context)
	respondError(c *gin.Context, httpStatus int, err error)
	authHandler(c *gin.Context)
	getPersons(c *gin.Context)
	updatePerson(c *gin.Context)
	createPerson(c *gin.Context)
	batchCSV(c *gin.Context)
	errorHandler(err error) (int, ErrorResponse)
}

type personRequestHandlerImpl struct {
	ctx           logger.Logger
	requestCtx    context.Context
	cfg           config.Config
	personService services.PersonService
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func PersonRequestHandlerInit() PersonRequestHandler {
	return &personRequestHandlerImpl{}
}

func (h *personRequestHandlerImpl) healthcheck(c *gin.Context) {
	h.ctx.Info("Healthcheck requested")
	c.Status(http.StatusOK)
}

func (h *personRequestHandlerImpl) configRequest(c *gin.Context) {
	h.requestCtx = c.Request.Context()
	h.cfg = config.Instance()
	h.ctx = *h.cfg.GetLogging()

	h.ctx.Info("Starting lambda handler")

	db := database.Database{}
	_, err := db.Start(h.ctx, h.cfg)
	if err != nil {
		h.ctx.Error(err)
		c.AbortWithStatusJSON(500, ErrorResponse{Error: "Server cannot response"})
	}

	var personDao = persondao.PersonDaoInit(h.ctx, &db)
	h.personService = services.PersonServiceInit(h.ctx, h.cfg, personDao)

	c.Next()
}

func (h *personRequestHandlerImpl) respondError(c *gin.Context, httpStatus int, err error) {
	if httpStatus < http.StatusInternalServerError {
		h.ctx.Warn("httpStatus:", httpStatus, err)
	} else {
		h.ctx.ErrorL("httpStatus: ", httpStatus, err)
	}

	c.AbortWithStatusJSON(httpStatus, ErrorResponse{Error: err.Error()})
}

// TODO: cognito or Auth for lambdas
func (h *personRequestHandlerImpl) authHandler(c *gin.Context) {

	c.Next()
}

/*
Given the error set the corresponding status code, log the error and return the JSON response
*/
func (h *personRequestHandlerImpl) errorHandler(err error) (int, ErrorResponse) {
	var status = http.StatusInternalServerError

	if merry.Is(err, services.FailedToCreate, services.FailedToUpdate, services.ErrPersonNotFound) {
		status = http.StatusBadRequest
	}

	if merry.Is(err, csv.ErrFileMaxSize, csv.ErrIncorrectCsvFormat, services.ErrReadingFile) {
		status = http.StatusUnprocessableEntity
	}

	if status < http.StatusInternalServerError {
		h.ctx.ErrorL("httpStatus", status, err)
		return status, ErrorResponse{Error: err.Error()}
	} else {
		h.ctx.ErrorL("httpStatus", status, err)
		return status, ErrorResponse{Error: "Unexpected error"}
	}
}

func (h *personRequestHandlerImpl) getPersons(c *gin.Context) {
	var args person.PersonSearchArgs
	if err := c.ShouldBind(&args); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	persons, err := h.personService.GetPersons(args)
	if err != nil {
		c.AbortWithStatusJSON(h.errorHandler(err))
		return
	}

	personsResponse := person.ToPersonsResponse(persons)
	c.JSON(http.StatusOK, personsResponse)
}

func (h *personRequestHandlerImpl) updatePerson(c *gin.Context) {
	var personReq person.PersonUpdateRequest
	if err := c.ShouldBind(&personReq); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	pId := c.Param("id")
	a, err := strconv.ParseUint(pId, 10, 10)
	if err != nil {
		h.respondError(c, http.StatusBadRequest, merry.New("Invalid request parameter format"))
		return
	}
	personReq.Id = uint(a)

	perUpdated, err := h.personService.UpdatePerson(personReq)
	if err != nil {
		c.AbortWithStatusJSON(h.errorHandler(err))
		return
	}

	c.JSON(http.StatusOK, perUpdated.ToPersonResponse())
}

func (h *personRequestHandlerImpl) createPerson(c *gin.Context) {
	h.ctx.Info("Create person requested")
	var personReq person.PersonRequest
	if err := c.ShouldBind(&personReq); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	p, err := h.personService.CreatePerson(personReq)
	if err != nil {
		c.AbortWithStatusJSON(h.errorHandler(err))
		return
	}

	c.JSON(http.StatusOK, p.ToPersonResponse())
}

func (h *personRequestHandlerImpl) batchCSV(c *gin.Context) {
	var csvfile csv.FileUpload
	if err := c.ShouldBind(&csvfile); err != nil {
		h.respondError(c, http.StatusBadRequest, err)
		return
	}

	persons, err := h.personService.BatchPersons(csvfile)
	if err != nil {
		c.AbortWithStatusJSON(h.errorHandler(err))
		return
	}

	personsResponse := person.ToPersonsResponse(persons)
	c.JSON(http.StatusOK, personsResponse)
}
