package onlyhttp

import (
	"fmt"
	"net/http"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

type PostsCallback struct {
	//Функция, просматривающая посты
	ListPosts func() []tdao.PostDao
}

type PostCallback struct {
	//Получает пост с комментариями и
	GetPost func(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM)
	//Загружает новый пост
	CreatePost func(user uint64, message *string) (uint64, tdao.PROBLEM)
	//Обновляет пост, если uid не совпадают, ошибка
	UpdatePost func(post uint64, user uint64, message *string) tdao.PROBLEM
	//Удаляет пост, если uid не совпадают, ошибка
	DeletePost func(post uint64, user uint64) tdao.PROBLEM
}

type MutePostCallback struct {
	//Запрещает добавлять комментарии
	MutePost func(post uint64, user uint64) tdao.PROBLEM
}

type UnmutePostCallback struct {
	//Разрешает добавлять комментарии
	UnmutePost func(post uint64, user uint64) tdao.PROBLEM
}

type CommentCallback struct {
	//Получает комментарий
	GetComment func(comment uint64) (tdao.PROBLEM, *tdao.CommentDao)
	//Загружает новый комментарий
	CreateComment func(user uint64, parent uint64, message *string) (uint64, tdao.PROBLEM)
	//Обновляет комментарий, если uid не совпадают, ошибка
	UpdateComment func(comment uint64, user uint64, message *string) tdao.PROBLEM
	//Удаляет комментарий, если uid не совпадают, ошибка
	DeleteComment func(comment uint64, user uint64) tdao.PROBLEM
}

// Структура, содержащая колбеки
type ServerCallbacks struct {
	//Логгер
	Log          func(string)
	Posts_       PostsCallback
	Post_        PostCallback
	Post_mute_   MutePostCallback
	Post_unmute_ UnmutePostCallback
	Comment_     CommentCallback
}

// Тестовый метод
func (s ServerCallbacks) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Hiiiiiii :3")
}

// Анализирует запросы к /Post
func (s ServerCallbacks) Post(w http.ResponseWriter, r *http.Request) {
	Post(
		w,
		r,
		s.Post_.GetPost,
		s.Log,
		s.Post_.CreatePost,
		s.Post_.UpdatePost,
		s.Post_.DeletePost,
	)
}

func (s ServerCallbacks) Comment(w http.ResponseWriter, r *http.Request) {
	Comment(
		w,
		r,
		s.Comment_.GetComment,
		s.Comment_.CreateComment,
		s.Comment_.UpdateComment,
		s.Comment_.DeleteComment,
	)
}

func (s ServerCallbacks) Posts(w http.ResponseWriter, r *http.Request) {
	Posts(
		w,
		r,
		s.Posts_.ListPosts,
	)
}

func (s ServerCallbacks) PostMute(w http.ResponseWriter, r *http.Request) {
	PostMute(
		w,
		r,
		s.Post_mute_.MutePost,
	)
}

func (s ServerCallbacks) PostUnmute(w http.ResponseWriter, r *http.Request) {
	PostUnmute(
		w,
		r,
		s.Post_unmute_.UnmutePost,
	)
}
