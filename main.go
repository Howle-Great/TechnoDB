package main

import (
	"./dbhandlers"
	"./handlers"
	// "github.com/jackc/pgx"
	// "fmt"
	// "./models"
)

func main() {
	err := dbhandlers.DB.Connetc()
	router := handlers.CreateRouter()	
	router.Run(":5000")
	if err != nil {
		panic(err)
	}
	// err = nil
	// fmt.Println("Starting server on 127.0.0.1:5000")
	// var tmp models.User
	// _ = dbhandlers.DB.Pool.QueryRow(`CREATE EXTENSION IF NOT EXISTS citext;`)
	// row := dbhandlers.DB.Pool.QueryRow(`
	// CREATE TABLE fdsa (
	// 	"nickname1" CITEXT UNIQUE PRIMARY KEY,
	// 	"email1"    CITEXT UNIQUE NOT NULL,
	// 	"fullname1" CITEXT NOT NULL,
	// 	"about1"    TEXT
	// );
	// `)
	// fmt.Println(row)
	// row = dbhandlers.DB.Pool.QueryRow(`
	// INSERT INTO fdsa (nickname1, email1, fullname1, about1)
	// 			VALUES ($1, $2, $3, $4) 
	// `,
	// "asd",
	// "email",
	// "name",
	// "about")
	// fmt.Println(row)
	// row3 := dbhandlers.DB.Pool.QueryRow(`
	// 	SELECT * FROM fdsa
	// `).Scan(
	// 	tmp.Nickname,
	// 	tmp.Email,
	// 	tmp.Fullname,
	// 	tmp.About)
	// fmt.Printf("\nDB say: %d", row3)
	// fmt.Println(err)

	// switch ErrorCode(err) {
	// 	case pgxOK:
	// 		return f, nil
	// 	case pgxErrUnique:
	// 		forum, _ := GetForum(f.Slug)
	// 		return forum, ForumIsExist
	// 	case pgxErrNotNull:
	// 		return nil, UserNotFound
	// 	default:
	// 		return nil, err
	// }
}