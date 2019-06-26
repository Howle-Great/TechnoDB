package dbhandlers

import (
	"../models"
)

// /user/{nickname}/create Создание нового пользователя
func CreateUser(u *models.User) (*models.Users, error)  {
	rows, err := DB.Pool.Exec(
		`
			INSERT
			INTO users ("nickname", "fullname", "email", "about")
			VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING
		`,
		&u.Nickname,
		&u.Fullname,
		&u.Email,
		&u.About,
	)
	if err != nil {
		return nil, err
	}

	if rows.RowsAffected() == 0 { // пользователь уже есть
		users := models.Users{}
		queryRows, err := DB.Pool.Query(`
				SELECT "nickname", "fullname", "email", "about"
				FROM users
				WHERE "nickname" = $1 OR "email" = $2
			`, 
			&u.Nickname, 
			&u.Email)
		defer queryRows.Close()

		if err != nil {
			return nil, err
		}

		for queryRows.Next() {
			user := models.User{}
			queryRows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
			users = append(users, &user)
		}
		return &users, UserIsExist
	}

	return nil, nil
}

// /user/{nickname}/profile Получение информации о пользователе
func GetUser(nickname string) (*models.User, error) {
	user := models.User{}

	err := DB.Pool.QueryRow(`
			SELECT "nickname", "fullname", "email", "about"
			FROM users
			WHERE "nickname" = $1
		`, 
		nickname).Scan(
			&user.Nickname,
			&user.Fullname,
			&user.Email,
			&user.About,
		)

	if err != nil {
		return nil, UserNotFound
	}

	return &user, nil
}

// /user/{nickname}/profile Изменение данных о пользователе
func UpdateUser(user *models.User) error {
	err := DB.Pool.QueryRow(
		`
			UPDATE users
			SET fullname = coalesce(nullif($2, ''), fullname),
				email = coalesce(nullif($3, ''), email),
				about = coalesce(nullif($4, ''), about)
			WHERE "nickname" = $1
			RETURNING nickname, fullname, email, about
		`,
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
	).Scan(
		&user.Nickname,
		&user.Fullname,
		&user.Email,
		&user.About,
	)

	if err != nil {
		if ErrorCode(err) != pgxOK {
			return UserUpdateConflict
		}
		return UserNotFound
	}

	return nil
}
