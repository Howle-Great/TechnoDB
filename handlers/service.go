package handlers

import (
	"../dbhandlers"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/swag"
)


// /service/status Получение инфомарции о базе данных
func GetStatus(c *gin.Context) {

	result := dbhandlers.GetStatus()
	resp, err := swag.WriteJSON(result)

	switch err {
		case nil:
			c.JSON( 200, resp)
		default:
			c.JSON( 500, []byte(err.Error()))
	}
}

// /service/clear Очистка всех данных в базе
func Clear(c *gin.Context) {
	dbhandlers.Clear()
	c.JSON( 200, []byte("Очистка базы успешно завершена"))
}

