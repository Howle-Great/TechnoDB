package handlers

import (
	// "io/ioutil"
	"../dbhandlers"
	"encoding/json"
	"../models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/go-openapi/swag"
	"log"
)

// /user/{nickname}/create Создание нового пользователя
func CreateUser(c *gin.Context) {
	var newUser models.User
	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	newUser.Nickname = c.Param("nickname")

	if err != nil {
		c.JSON( 500, []byte(err.Error()))
		return
	}	

	result, err := dbhandlers.CreateUser(&newUser)
	log.Panicln(result);

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(newUser)
			c.JSON( 201, resp)
		case dbhandlers.UserIsExist:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 409, resp)
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /user/{nickname}/profile Получение информации о пользователе
func GetUser(c *gin.Context) {
	nickname := c.Param("nickname")

	result, err := dbhandlers.GetUser(nickname)

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.UserNotFound:
			c.JSON( 404, []byte(makeErrorUser(nickname)))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /user/{nickname}/profile Изменение данных о пользователе
func UpdateUser(c *gin.Context) {
	nickname := c.Param("nickname")

	var newUser models.User
	err := json.NewDecoder(c.Request.Body).Decode(&newUser)
	newUser.Nickname = c.Param("nickname")

	if err != nil {
		c.JSON( 500, []byte(err.Error()))
		return
	}	

	err = dbhandlers.UpdateUser(&newUser)

	switch err {
	case nil:
		resp, _ := swag.WriteJSON(newUser)
		c.JSON( 200, resp)
	case dbhandlers.UserNotFound:
		c.JSON( 404, []byte(makeErrorUser(nickname)))
	case dbhandlers.UserUpdateConflict:
		c.JSON( 409, []byte(makeErrorEmail(nickname)))
	default:		
		c.JSON( 500, []byte(err.Error()))
	}
}

