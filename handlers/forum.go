package handlers

import (
	"../dbhandlers"
	"../models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
	_ "github.com/lib/pq"
)

// /forum/create Создание форума
func CreateForum(c *gin.Context) {
	var forum models.Forum
	slug := c.Param("slug")
	if slug != "create" {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
		return
	}

	err := json.NewDecoder(c.Request.Body).Decode(&forum)
	if err != nil {
		c.JSON(500, []byte(err.Error()))
		return
	}

	result, err := dbhandlers.CreateForum(&forum)

	switch err {
		case nil:
			resp, _ := json.Marshal(result)
			c.JSON(201, resp)
		case dbhandlers.UserNotFound:
			c.JSON(404, []byte(makeErrorUser(forum.User)))
		case dbhandlers.ForumIsExist:
			resp, _ := json.Marshal(result)
			c.JSON(409, resp)
		default:		
			c.JSON(500, []byte(err.Error()))
	}
}

// /forum/{slug}/create Создание ветки
func CreateSlug(c *gin.Context) {
	var slug models.Thread
	err := json.NewDecoder(c.Request.Body).Decode(&slug)
	slugParam := c.Param("slug")

	if err != nil {
		c.JSON(500, []byte(err.Error()))
		return
	}
	// thread := &models.Thread{}
	// thread.Forum = slug // иначе не знаю как
	slug.Forum = slugParam

	result, err := dbhandlers.CreateForumThread(&slug)

	switch err {
		case nil:
			resp, _ := json.Marshal(result)
			c.JSON(201, resp)
		case dbhandlers.UserNotFound:
			c.JSON(404, []byte(makeErrorUser(slug.Author)))
		case dbhandlers.ForumIsExist:
			resp, _ := json.Marshal(result)
			c.JSON(409, resp)
		default:		
			c.JSON(500, []byte(err.Error()))
	}
}

// /forum/{slug}/details Получение информации о форуме
func GetForum(c *gin.Context) {
	slug := c.Param("slug")

	result, err := dbhandlers.GetForum(slug)

	switch err {
		case nil:
			resp, _ := json.Marshal(result)
			c.JSON( 200, resp)
		case dbhandlers.ForumNotFound:
			c.JSON( 404, []byte(makeErrorForum(slug)))
		default:
			c.JSON( 500, []byte(err.Error()))
	}
}

// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreads(c *gin.Context) {
	slug := c.Param("slug")
	args := c.Request.URL.Query()
	limit := args.Get("limit")
	since := args.Get("since")
	desc := args.Get("desc")

	if limit == "" { limit = "1" }
	if desc == "" { desc = "false" }

	result, err := dbhandlers.GetForumThreads(slug, limit, since, desc)
	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.ForumNotFound:
			c.JSON( 404, []byte(makeErrorForum(slug)))
		default:
			c.JSON( 500, []byte(err.Error()))
	}
}

// /forum/{slug}/users Пользователи данного форума
func GetForumUsers(c *gin.Context) {
	slug := c.Param("slug")
	q := c.Request.URL.Query()
	var limit, since, desc []string
	limit = q["limit"]
	since = q["since"]
	desc = q["desc"]

	if limit[0] == "" { limit[0] = "1" }
	if desc[0] == "" { desc[0] = "false" }

	result, err := dbhandlers.GetForumUsers(slug, limit[0], since[0], desc[0])

	switch err {
		case nil:
			resp, _ := swag.WriteJSON(result)
			c.JSON( 200, resp)
		case dbhandlers.ForumNotFound:
			c.JSON( 404, []byte(makeErrorUser(slug)))
		default:
			c.JSON( 500, []byte(err.Error()))
	}
}