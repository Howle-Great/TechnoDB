package dbhandlers

import (
	"strconv"
	"../models"
)

// /post/{id}/details Получение информации о ветке обсуждения
func GetPost(id int) (*models.Post, error) {
	post := models.Post{}

	err := DB.Pool.QueryRow(
		`
			SELECT id, author, message, forum, thread, created, "isEdited", parent
			FROM posts 
			WHERE id = $1
		`,
		id,
	).Scan(
		&post.ID,
		&post.Author,
		&post.Message,
		&post.Forum,
		&post.Thread,
		&post.Created,
		&post.IsEdited,
		&post.Parent,
	)

	if err == nil {
		return &post, nil
	} else if (err.Error() == noRowsInResult) {
		return nil, PostNotFound
	} else {
		return nil, err
	}
}


// /post/{id}/details Получение информации о ветке обсуждения
func GetPostFull(id int, related []string) (*models.PostFull, error) {
	postFull := models.PostFull{}
	var err error
	postFull.Post, err = GetPost(id)
	if err != nil {
		return nil, err
	}

	for _, model := range related {
		switch model {
		case "thread":
			postFull.Thread, err = GetThread(strconv.Itoa(int(postFull.Post.Thread)))
		case "forum":
			postFull.Forum, err = GetForum(postFull.Post.Forum)
		case "user":
			postFull.Author, err = GetUser(postFull.Post.Author)
		}

		if err != nil {
			return nil, err
		}
	}

	return &postFull, nil
}

// /post/{id}/details Изменение сообщения
func UpdatePost(postUpdate *models.PostUpdate, id int) (*models.Post, error) {
	post, err := GetPost(id)
	if err != nil {
		return nil, PostNotFound
	}

	if len(postUpdate.Message) == 0 {
		return post, nil
	}

	rows := DB.Pool.QueryRow(`
			UPDATE posts 
			SET message = COALESCE($2, message), "isEdited" = ($2 IS NOT NULL AND $2 <> message) 
			WHERE id = $1 
			RETURNING author::text, created, forum, "isEdited", thread, message, parent
		`, 
		strconv.Itoa(id), 
		&postUpdate.Message)

	err = rows.Scan(
		&post.Author,
		&post.Created,
		&post.Forum,
		&post.IsEdited,
		&post.Thread,
		&post.Message,
		&post.Parent,
	)

	if err == nil {
		return post, nil
	} else if (err.Error() == noRowsInResult) {
		return nil, PostNotFound
	} else {
		return nil, err
	}
}
