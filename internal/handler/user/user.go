package user

import (
	"errors"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/model"
	"gravitum-test-app/internal/service"
	"gravitum-test-app/pkg/helper"
	"gravitum-test-app/pkg/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	cfg     *config.Config
	service service.UserService
	log     *logger.Logger
}

func NewHandler(
	cfg *config.Config,
	service service.UserService,
	log *logger.Logger,
) *UserHandler {
	return &UserHandler{
		cfg:     cfg,
		service: service,
		log:     log,
	}
}

func (h *UserHandler) GetList(c *gin.Context) {
	result, err := h.service.GetList(c.Request.Context())
	if err != nil {
		h.log.Errorf("internal server error: %s", err)
		c.JSON(http.StatusInternalServerError, model.WrapError(http.StatusInternalServerError, err.Error()))
		return
	}

	h.log.Debug("get user list")
	c.JSON(http.StatusOK, model.WrapResponse(http.StatusOK, result))
}

func (h *UserHandler) Get(c *gin.Context) {

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Join(err, model.ErrRequestInvalidUrlParams)
		h.log.Errorf("bad request error: request param error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	result, err := h.service.Get(c.Request.Context(), uint(idInt))
	if err != nil {
		if errors.Is(err, model.ErrSqlNoRows) ||
			errors.Is(err, model.ErrNoUserWithSuchId) {
			h.log.Errorf("unprocessable entity error: %s", err)
			c.JSON(http.StatusUnprocessableEntity, model.WrapError(http.StatusUnprocessableEntity, err.Error()))
			return
		}

		h.log.Errorf("internal server error: %s", err)
		c.JSON(http.StatusInternalServerError, model.WrapError(http.StatusInternalServerError, err.Error()))
		return
	}

	h.log.Debugf("get user, id = %d", idInt)
	c.JSON(http.StatusOK, model.WrapResponse(http.StatusOK, result))
}

func (h *UserHandler) Create(c *gin.Context) {
	var bodyParams model.CreateUserRequest

	err := c.BindJSON(&bodyParams)
	if err != nil {
		err = errors.Join(err, model.ErrRequestInvalidBodyParams)
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	// name
	if bodyParams.Name == nil {
		err = model.ErrRequestNameRequired
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	name := helper.SanitizeInput(*bodyParams.Name)
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		err := model.ErrRequestNameRequired
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	// surname
	var surname *string
	if bodyParams.Surname != nil {
		surnameStr := helper.SanitizeInput(*bodyParams.Surname)
		surnameStr = strings.TrimSpace(surnameStr)
		if len(name) > 0 {
			surname = &surnameStr
		}
	}

	err = h.service.Create(c.Request.Context(), name, surname)
	if err != nil {
		h.log.Errorf("internal server error: %s", err)
		c.JSON(http.StatusInternalServerError, model.WrapError(http.StatusInternalServerError, err.Error()))
		return
	}

	if surname == nil {
		h.log.Debugf("user created, name=%s", name)
	} else {
		h.log.Debugf("user created, name=%s, surname=%s", name, *surname)
	}

	c.JSON(http.StatusCreated, model.WrapResponse(http.StatusCreated, nil))
}

func (h *UserHandler) Update(c *gin.Context) {

	idInt, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err = errors.Join(err, model.ErrRequestInvalidUrlParams)
		h.log.Errorf("bad request error: request param error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	var bodyParams model.UpdateUserRequest

	err = c.BindJSON(&bodyParams)
	if err != nil {
		err = errors.Join(err, model.ErrRequestInvalidBodyParams)
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	// name
	if bodyParams.Name == nil {
		err = model.ErrRequestNameRequired
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	name := helper.SanitizeInput(*bodyParams.Name)
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		err := model.ErrRequestNameRequired
		h.log.Errorf("bad request error: %s", err)
		c.JSON(http.StatusBadRequest, model.WrapError(http.StatusBadRequest, err.Error()))
		return
	}

	// surname
	var surname *string
	if bodyParams.Surname != nil {
		surnameStr := helper.SanitizeInput(*bodyParams.Surname)
		surnameStr = strings.TrimSpace(surnameStr)
		if len(name) > 0 {
			surname = &surnameStr
		}
	}

	err = h.service.Update(c.Request.Context(), uint(idInt), name, surname)
	if err != nil {
		h.log.Errorf("internal server error: %s", err)
		c.JSON(http.StatusInternalServerError, model.WrapError(http.StatusInternalServerError, err.Error()))
		return
	}

	if surname == nil {
		h.log.Debugf("user updated, id=%d, name=%s", idInt, name)
	} else {
		h.log.Debugf("user updated, id=%d, name=%s, surname=%s", idInt, name, *surname)
	}
	h.log.Debugf("get user, id = %d", idInt)

	c.JSON(http.StatusOK, model.WrapResponse(http.StatusOK, nil))
}
