package handlers

import (
	"encoding/json"
	// "strconv"
	"../dbhandlers"
	"../models"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
)

// /thread/{slug_or_id}/details Получение информации о ветке обсуждения
func GetThread(c *gin.Context) {
	param := c.Param("slug_or_id")

	result, err := dbhandlers.GetThread(param)

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.ThreadNotFound:
			c.JSON( 404, []byte(makeErrorThread(param)))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThread(c *gin.Context) {
	param := c.Param("slug_or_id")

	var threadUpdate models.ThreadUpdate
	err := json.NewDecoder(c.Request.Body).Decode(&threadUpdate)

	if err != nil {
		c.JSON( 500, []byte(err.Error()))
		return
	}

	result, err := dbhandlers.UpdateThread(&threadUpdate, param)

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.PostNotFound:
			c.JSON( 404, []byte(makeErrorThread(param)))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/create Создание новых постов
func CreatePost(c *gin.Context) {
	param := c.Param("slug_or_id")

	var post models.Posts
	err := json.NewDecoder(c.Request.Body).Decode(&post)
	if err != nil {
		c.JSON( 500, []byte(err.Error()))
		return
	}

	result, err := dbhandlers.CreateThread(&post, param)

	resp, _ := swag.WriteJSON(result)

	switch err {
		case nil:
			c.JSON( 201, resp)
		case dbhandlers.ThreadNotFound:
			c.JSON( 404, []byte(makeErrorThreadID(param)))
		case dbhandlers.UserNotFound:
			c.JSON( 404, []byte(makeErrorPostAuthor(param)))
		case dbhandlers.PostParentNotFound:
			c.JSON( 409, []byte(makeErrorThreadConflict()))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPosts(c *gin.Context) {
	param := c.Param("slug_or_id")
	q := c.Request.URL.Query()
	limit := q["limit"][0]
	sort := q["sort"][0]
	since := q["since"][0]
	desc := q["desc"][0]

	if limit == "" { limit = "1";	}
	if sort == "" { sort = "flat"; }
	if desc == "" { desc = "false";	}

	result, err := dbhandlers.GetThreadPosts(param, limit, since, sort, desc)
	
	resp, _ := swag.WriteJSON(result)

	switch err {
		case nil:
			c.JSON( 200, resp)
		case dbhandlers.ForumNotFound:
			c.JSON( 404, []byte(makeErrorThread(param)))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}

// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVote(c *gin.Context) {
	param := c.Param("slug_or_id")
	var votes models.Vote
	err := json.NewDecoder(c.Request.Body).Decode(&votes)
	if err != nil {
		c.JSON( 500, []byte(err.Error()))
		return
	}	

	var vote models.Vote
	result, err := dbhandlers.MakeThreadVote(&vote, param)	
	
	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.ForumNotFound:
			c.JSON( 404, []byte(makeErrorThread(param)))
		case dbhandlers.UserNotFound:
			c.JSON( 404, []byte(makeErrorUser(param)))
		default:		
			c.JSON( 500, []byte(err.Error()))
	}
}
