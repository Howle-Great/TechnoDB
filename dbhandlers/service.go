package dbhandlers

import (
	"../models"
)

// /service/status Получение инфомарции о базе данных
func GetStatus() *models.Status {
	status := &models.Status{}
	DB.pool.QueryRow(
		`
			SELECT 
			(SELECT COUNT(*) FROM users) AS users,
			(SELECT COUNT(*) FROM forums) AS forums,
			(SELECT COUNT(*) FROM posts) AS posts,
			(SELECT COALESCE(SUM(threads), 0) FROM forums WHERE threads > 0) AS threads
		`,
	).Scan(
		&status.User,
		&status.Forum,
		&status.Post,
		&status.Thread,
	)
	return status
}

// /service/clear Очистка всех данных в базе
func Clear() {
	DB.pool.Exec(`
		TRUNCATE users, forums, threads, posts, votes, forum_users;
	`)
}