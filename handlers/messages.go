package handlers

import(
	"fmt"
)

func makeErrorUser(s string) string {
	return fmt.Sprintf(`{"message": "Can't find user by nickname: %s"}`, s)
}

func makeErrorEmail(s string) string {
	return fmt.Sprintf(`{"message": "This email is already registered by user: %s"}`, s)
}

func makeErrorForum(s string) string {
	return fmt.Sprintf(`{"message": "Can't find forum with slug: %s"}`, s)
}

func makeErrorThread(s string) string {
	return fmt.Sprintf(`{"message": "Can't find thread by slug: %s"}`, s)
}

func makeErrorThreadConflict() string {
	return `{"message": "Parent post was created in another thread"}`
}

func makeErrorThreadID(s string) string {
	return fmt.Sprintf(`{"message": "Can't find thread by slug: %s"}`, s)
}

func makeErrorPost(s string) string {
	return fmt.Sprintf(`{"message": "Can't find post with id: %s"}`, s)
}

func makeErrorPostAuthor(s string) string {
	return fmt.Sprintf(`{"message": "Can't find post author by nickname: %s"}`, s)
}