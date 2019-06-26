package main

import ()
// import (
// 	"github.com/gin-gonic/gin"
// 	"fmt"
// )

// func main() {
// 	r := CreateRouter()
// 	r.Run(":3000")
// 	fmt.Println("Start server on port 3000")
// }

// func CreateRouter() (*gin.Engine) {
// 	r := gin.New()
// 	forumApi := r.Group("api/")
// 	forumApi.POST(":slug/create", CreateSlug)
// 	forumApi.GET(":slug/details", GetForum)
// 	return r
// }

// func CreateForum(c *gin.Context) {
// 	c.JSON(200, "Qwe")
// 	panic("1")
// }

// func CreateSlug(c *gin.Context) {
// 	slug := c.Param("slug")
// 	c.JSON(201, gin.H{"slug": slug})
// 	// panic("2")
// }

// func GetForum(c *gin.Context) {
// 	slug := c.Param("slug")
// 	c.JSON(202, gin.H{"slug": slug})
// 	// panic("3")
// }

// /*
// curl -v -d "" http://localhost:3000/api/zxc/create
// curl -v http://localhost:3000/api/zSSD12c/details
// */