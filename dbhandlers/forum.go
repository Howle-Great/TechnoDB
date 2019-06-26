package dbhandlers

import (
	"github.com/jackc/pgx"
	"../models"
)

// /forum/create Создание форума
func CreateForum(f *models.Forum) (*models.Forum, error) {
	err := DB.pool.QueryRow(
		`
			INSERT INTO forum (slug, title, "user")
			VALUES ($1, $2, (
				SELECT nickname FROM users WHERE nickname = $3
			)) 
			RETURNING "user"
		`,
		&f.Slug,
		&f.Title,
		&f.User,
	).Scan(&f.User)

	switch ErrorCode(err) {
		case pgxOK:
			return f, nil
		case pgxErrUnique:
			forum, _ := GetForum(f.Slug)
			return forum, ForumIsExist
		case pgxErrNotNull:
			return nil, UserNotFound
		default:
			return nil, err
	}
}

// /forum/{slug}/details Получение информации о форуме
func GetForum(slug string) (*models.Forum, error) {
	f := models.Forum{}

	err := DB.pool.QueryRow(
		`
			SELECT slug, title, "user", posts, threads
			FROM forum
			WHERE slug = $1
		`, 
		slug,
		).Scan(
			&f.Slug,
			&f.Title,
			&f.User,
			&f.Posts,
			&f.Threads,
	)

	if err != nil {
		return nil, ForumNotFound
	}

	return &f, nil
}

// /forum/{slug}/create Создание ветки
func CreateForumThread(t *models.Thread) (*models.Thread, error) {
	if t.Slug != "" {
		thread, err := GetThread(t.Slug)
		if err == nil {
			return thread, ThreadIsExist
		}
	}

	err := DB.pool.QueryRow(
		`
			INSERT INTO threads (author, created, message, title, slug, forum)
			VALUES ($1, $2, $3, $4, $5, (SELECT slug FROM forums WHERE slug = $6)) 
			RETURNING author, created, forum, id, message, title
		`, 
		&t.Author,
		&t.Created, 
		&t.Message,
		&t.Title,
		&t.Slug,
		&t.Forum,
	).Scan(
		&t.Author,
		&t.Created, 
		&t.Forum,
		&t.ID,
		&t.Message,
		&t.Title,
	)
	
	switch ErrorCode(err) {
		case pgxOK:
			return t, nil
		case pgxErrNotNull:
			return nil, UserNotFound //UserNotFound
		case pgxErrForeignKey:
			return nil, ForumIsExist //ForumIsExist
		default:
			return nil, err
	}
}

var queryForumWithSience = map[string]string {
	"true": `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created <= $2::TEXT::TIMESTAMPTZ
		ORDER BY created DESC
		LIMIT $3::TEXT::INTEGER
	`,
	"false": `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1 AND created >= $2::TEXT::TIMESTAMPTZ
		ORDER BY created
		LIMIT $3::TEXT::INTEGER
	`,
}

var queryForumNoSience = map[string]string {
	"true": `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created DESC
		LIMIT $2::TEXT::INTEGER
	`,
	"false": `
		SELECT author, created, forum, id, message, slug, title, votes
		FROM threads
		WHERE forum = $1
		ORDER BY created
		LIMIT $2::TEXT::INTEGER
	`,
}

// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreads(slug, limit, since, desc string) (*models.Threads, error) {
	var rows *pgx.Rows
	var err error

	if since != "" {
		rows, err = DB.pool.Query(queryForumWithSience[desc], slug, since, limit)
	} else {
		rows, err = DB.pool.Query(queryForumNoSience[desc], slug, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, ForumNotFound
	}
	
	threads := models.Threads{}
	for rows.Next() {
		t := models.Thread{}
		err = rows.Scan(
			&t.Author,
			&t.Created,
			&t.Forum,
			&t.ID,
			&t.Message,
			&t.Slug,
			&t.Title,
			&t.Votes,
		)
		threads = append(threads, &t)
	}

	if len(threads) == 0 {
		_, err := GetForum(slug)
		if err != nil {
			return nil, ForumNotFound
		}
	}	
	return &threads, nil
}

var queryForumUserWithSience = map[string]string {
	"true": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			) 
			AND LOWER(nickname) < LOWER($2::TEXT)
		ORDER BY nickname DESC
		LIMIT $3::TEXT::INTEGER
	`,
	"false": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			) 
			AND LOWER(nickname) > LOWER($2::TEXT)
		ORDER BY nickname
		LIMIT $3::TEXT::INTEGER
	`,
}

var queryForumUserNoSience = map[string]string {
	"true": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			)
		ORDER BY nickname DESC
		LIMIT $2::TEXT::INTEGER
	`,
	"false": `
		SELECT nickname, fullname, about, email
		FROM users
		WHERE nickname IN (
				SELECT forum_user FROM forum_users WHERE forum = $1
			)
		ORDER BY nickname
		LIMIT $2::TEXT::INTEGER
	`,
}

// /forum/{slug}/users Пользователи данного форума
func GetForumUsers(slug, limit, since, desc string) (*models.Users, error) {
	var rows *pgx.Rows
	var err error
	
	if since != "" {
		rows, err = DB.pool.Query(queryForumUserWithSience[desc], slug, since, limit)
	} else {
		rows, err = DB.pool.Query(queryForumUserNoSience[desc], slug, limit)
	}
	defer rows.Close()

	if err != nil {
		return nil, ForumNotFound
	}
	
	users := models.Users{}
	for rows.Next() {
		u := models.User{}
		err = rows.Scan(
			&u.Nickname,
			&u.Fullname,
			&u.About,
			&u.Email,
		)
		users = append(users, &u)
	}

	if len(users) == 0 {
		_, err := GetForum(slug)
		if err != nil {
			return nil, ForumNotFound
		}
	}	
	return &users, nil
}
