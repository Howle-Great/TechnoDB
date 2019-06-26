package handlers

import (
	"strings"
	"../dbhandlers"
	"../models"
	"github.com/gin-gonic/gin"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/go-openapi/swag"
	"strconv"
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPost(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	q := c.Request.URL.Query()
	relatedQuery := q["related"]
	related := []string{}
	related = append(related, strings.Split(string(relatedQuery[0]), ",")...)

	result, err := dbhandlers.GetPostFull(id, related)

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.PostNotFound:
			c.JSON( 404, []byte(makeErrorPost(string(id))))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /post/{id}/details Изменение сообщения
func UpdatePost(c *gin.Context) {
	var postUpdate models.PostUpdate
	id, _ := strconv.Atoi(c.Param("id"))
	err := json.NewDecoder(c.Request.Body).Decode(&postUpdate)

	result, err := dbhandlers.UpdatePost(&postUpdate, id)

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.PostNotFound:
			c.JSON( 404, []byte(makeErrorPost(string(id))))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}
