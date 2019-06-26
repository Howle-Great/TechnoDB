package dbhandlers

import (
	"errors"
)

// Ошибки БД
const (
	pgxOK				= ""
	pgxErrNotNull		= "23502"
	pgxErrForeignKey 	= "23503"	// может возникнуть при добавлении дубликата
	pgxErrUnique 		= "23505"
	noRowsInResult 		= "no rows in result set"
)

// Ошибки запросов
var (
	ForumIsExist			 = errors.New("Forum was created earlier")
	ForumNotFound			 = errors.New("Forum not found")
	ForumOrAuthorNotFound	 = errors.New("Forum or Author not found")
	UserNotFound			 = errors.New("User not found")
	UserIsExist				 = errors.New("User was created earlier")
	UserUpdateConflict		 = errors.New("User not updated")
	ThreadIsExist			 = errors.New("Thread was created earlier")
	ThreadNotFound			 = errors.New("Thread not found")
	PostParentNotFound		 = errors.New("No parent for thread")
	PostNotFound			 = errors.New("Post not found")
)