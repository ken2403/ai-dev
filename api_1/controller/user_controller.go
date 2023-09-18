package controller

import (
	"main/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var _ IUserController = &UserController{}

type IUserController interface {
	GetUserByID(c echo.Context) error
}

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) GetUserByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	user, notFound, err := uc.userService.GetUserByID(uint(id))
	if notFound {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
