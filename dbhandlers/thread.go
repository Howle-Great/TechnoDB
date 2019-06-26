package dbhandlers

import (
	"fmt"
	"strconv"
	"time"
	"strings"
	"../models"
	"github.com/jackc/pgx"
)

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func GetThread(param string) (*models.Thread, error) {
	var err error
	var thread models.Thread

	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		err = DB.pool.QueryRow(
			`
				SELECT id, title, author, forum, message, votes, slug, created
				FROM threads
				WHERE id = $1
			`,
			id,
		).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
	} else {
		err = DB.pool.QueryRow(
			`
				SELECT id, title, author, forum, message, votes, slug, created
				FROM threads
				WHERE slug = $1
			`,
			param,
		).Scan(
			&thread.ID,
			&thread.Title,
			&thread.Author,
			&thread.Forum,
			&thread.Message,
			&thread.Votes,
			&thread.Slug,
			&thread.Created,
		)
	}

	if err != nil {
		return nil, ThreadNotFound
	}

	return &thread, nil
}

// /thread/{slug_or_id}/details Обновление ветки
func UpdateThread(thread *models.ThreadUpdate, param string) (*models.Thread, error) {
	threadFound, err := GetThread(param)
	if err != nil {
		return nil, PostNotFound
	}

	updatedThread := models.Thread{}

	err = DB.pool.QueryRow(`
			UPDATE threads
			SET title = coalesce(nullif($2, ''), title),
				message = coalesce(nullif($3, ''), message)
			WHERE slug = $1
			RETURNING id, title, author, forum, message, votes, slug, created
		`,
		&threadFound.Slug,
		&thread.Title,
		&thread.Message,
	).Scan(
		&updatedThread.ID,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)

	if err != nil {
		return nil, err
	}

	return &updatedThread, nil
}

func authorExists(nickname string) bool {
	var user models.User
	err :=  DB.pool.QueryRow(
		`
			SELECT "nickname", "fullname", "email", "about"
			FROM users
			WHERE "nickname" = $1
		`,
		nickname,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	)

	if err != nil && err.Error() == noRowsInResult {
		return true
	}
	return false
}

const postID = `
	SELECT id
	FROM posts
	WHERE id = $1 AND thread IN (SELECT id FROM threads WHERE thread <> $2)
`

func parentExitsInOtherThread(parent int64, threadID int32) bool {
	var t int64
	err := DB.pool.QueryRow(postID, parent, threadID).Scan(&t)

	if err != nil && err.Error() == noRowsInResult {
		return false
	}
	return true
}

func parentNotExists(parent int64) bool {
	if parent == 0 {
		return false
	}

	var t int64
	err := DB.pool.QueryRow(`SELECT id FROM posts WHERE id = $1`, parent).Scan(&t); 
	
	if err != nil {
		return true
	}
	return false
}

func checkPost(p *models.Post, t *models.Thread) error {
	if authorExists(p.Author) {
		return UserNotFound
	}
	if parentExitsInOtherThread(p.Parent, t.ID) || parentNotExists(p.Parent) {
		return PostParentNotFound
	}
	return nil
}

// thread/{slug_or_id}/create Создание новых постов
func CreateThread(posts *models.Posts, param string) (*models.Posts, error) {
	thread, err := GetThread(param)
	if err != nil {
		return nil, err
	}

	postsNumber := len(*posts)
	if postsNumber == 0 {
		return posts, nil
	}

	dateTimeTemplate := "2006-01-02 15:04:05"
	created := time.Now().Format(dateTimeTemplate)
	query := strings.Builder{}
	query.WriteString("INSERT INTO posts (author, created, message, thread, parent, forum, path) VALUES ")
	queryBody := "('%s', '%s', '%s', %d, %d, '%s', (SELECT path FROM posts WHERE id = %d) || (SELECT last_value FROM posts_id_seq)),"
	for i, post := range *posts {
		err = checkPost(post, thread)
		if err != nil {
			return nil, err
		}

		temp := fmt.Sprintf(queryBody, post.Author, created, post.Message, thread.ID, post.Parent, thread.Forum, post.Parent)
		// удаление запятой в конце queryBody для последнего подзапроса
		if i == postsNumber - 1 {
			temp = temp[:len(temp) - 1]
		}
		query.WriteString(temp)
	}
	query.WriteString("RETURNING author, created, forum, id, message, parent, thread")

	tx, txErr := DB.pool.Begin()
	if txErr != nil {
		return nil, txErr
	}
	defer tx.Rollback()

	rows, err := tx.Query(query.String())
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	insertPosts := models.Posts{}
	for rows.Next() {
		post := models.Post{}
		rows.Scan(
			&post.Author,
			&post.Created,
			&post.Forum,
			&post.ID,
			&post.Message,
			&post.Parent,
			&post.Thread,
		)
		insertPosts = append(insertPosts, &post) 
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// по хорошему это впихнуть в хранимые процедуры, но нормальные ребята предпочитают костылить
	tx.Exec(`UPDATE forums SET posts = posts + $1 WHERE slug = $2`, len(insertPosts), thread.Forum)
	for _, p := range insertPosts {
		tx.Exec(`INSERT INTO forum_users VALUES ($1, $2) ON CONFLICT DO NOTHING`, p.Author, p.Forum)
	}

	tx.Commit()

	return &insertPosts, nil
}

var queryPostsWithSience = map[string]map[string]string {
	"true": map[string]string {
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND (path < (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
			ORDER BY path DESC
			LIMIT $3::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts p
			WHERE p.thread = $1 and p.path[1] IN (
				SELECT p2.path[1]
				FROM posts p2
				WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] < (SELECT p3.path[1] from posts p3 where p3.id = $2)
				ORDER BY p2.path DESC
				LIMIT $3
			)
			ORDER BY p.path[1] DESC, p.path[2:]
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND id < $2::TEXT::INTEGER
			ORDER BY id DESC
			LIMIT $3::TEXT::INTEGER
		`,
	},
	"false": map[string]string {
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND (path > (SELECT path FROM posts WHERE id = $2::TEXT::INTEGER))
			ORDER BY path
			LIMIT $3::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts p
			WHERE p.thread = $1 and p.path[1] IN (
				SELECT p2.path[1]
				FROM posts p2
				WHERE p2.thread = $1 AND p2.parent = 0 and p2.path[1] > (SELECT p3.path[1] from posts p3 where p3.id = $2::TEXT::INTEGER)
				ORDER BY p2.path
				LIMIT $3::TEXT::INTEGER
			)
			ORDER BY p.path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND id > $2::TEXT::INTEGER
			ORDER BY id
			LIMIT $3::TEXT::INTEGER
		`,
	},
}

var queryPostsNoSience = map[string]map[string]string {
	"true": map[string]string {
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY path DESC
			LIMIT $2::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND path[1] IN (
				SELECT path[1]
				FROM posts
				WHERE thread = $1
				GROUP BY path[1]
				ORDER BY path[1] DESC
				LIMIT $2::TEXT::INTEGER
			)
			ORDER BY path[1] DESC, path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1
			ORDER BY id DESC
			LIMIT $2::TEXT::INTEGER
		`,
	},
	"false": map[string]string {
		"tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY path
			LIMIT $2::TEXT::INTEGER
		`,
		"parent_tree": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 AND path[1] IN (
				SELECT path[1] 
				FROM posts 
				WHERE thread = $1 
				GROUP BY path[1]
				ORDER BY path[1]
				LIMIT $2::TEXT::INTEGER
			)
			ORDER BY path
		`,
		"flat": `
			SELECT id, author, parent, message, forum, thread, created
			FROM posts
			WHERE thread = $1 
			ORDER BY id
			LIMIT $2::TEXT::INTEGER
		`,
	},
}

// /thread/{slug_or_id}/posts Сообщения данной ветви обсуждения
func GetThreadPosts(param, limit, since, sort, desc string) (*models.Posts, error) {
	thread, err := GetThread(param)
	if err != nil {
		return nil, ForumNotFound
	}

	var rows *pgx.Rows

	if since != "" {
		query := queryPostsWithSience[desc][sort]
		rows, err = DB.pool.Query(query, thread.ID, since, limit)
	} else {
		query := queryPostsNoSience[desc][sort]
		rows, err = DB.pool.Query(query, thread.ID, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	posts := models.Posts{}
	for rows.Next() {
		post := models.Post{}

		err = rows.Scan(
			&post.ID,
			&post.Author,
			&post.Parent,
			&post.Message,
			&post.Forum,
			&post.Thread,
			&post.Created,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

// /thread/{slug_or_id}/vote Проголосовать за ветвь обсуждения
func MakeThreadVote(vote *models.Vote, param string) (*models.Thread, error) {
	var err error

	tx, txErr := DB.pool.Begin()
	if txErr != nil {
		return nil, txErr
	}
	defer tx.Rollback()

	var thread models.Thread
	if isNumber(param) {
		id, _ := strconv.Atoi(param)
		err = tx.QueryRow(`SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE id = $1`, id).Scan(
			&thread.ID,
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
	} else {
		err = tx.QueryRow(`SELECT id, author, created, forum, message, slug, title, votes FROM threads WHERE slug = $1`, param).Scan(
			&thread.ID,
			&thread.Author,
			&thread.Created,
			&thread.Forum,
			&thread.Message,
			&thread.Slug,
			&thread.Title,
			&thread.Votes,
		)
	}	
	if err != nil {
		return nil, ForumNotFound
	}

	var nick string
	err = tx.QueryRow(`SELECT nickname FROM users WHERE nickname = $1`, vote.Nickname).Scan(&nick)
	if err != nil {
		return nil, UserNotFound
	}

	rows, err := tx.Exec(`UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3;`, vote.Voice, thread.ID, vote.Nickname)
	if rows.RowsAffected() == 0 {
		_, err := tx.Exec(`INSERT INTO votes (nickname, thread, voice) VALUES ($1, $2, $3);`, vote.Nickname, thread.ID, vote.Voice)
		if err != nil {
			return nil, UserNotFound
		}
	}
	// если возник вопрос - в какой мемент делаем +1 к voice -> смотри init.sql

	err = tx.QueryRow(`SELECT votes FROM threads WHERE id = $1`, thread.ID).Scan(&thread.Votes)
	if err != nil {
		return nil, err
	}
	
	tx.Commit()

	return &thread, nil
}