package onlyhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

// Организует получение поста по Id
// Соответствует методу GET /Post?post_id=$postStr$&from=$from$&to=$to$
func post_get(
	w http.ResponseWriter,
	getPost func(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM),
	log func(string),
	postStr string,
	from string,
	to string,
) {
	//В случае неудачного парсинга, возвращается ошибка
	fromPos, errFrom := strconv.ParseUint(from, 10, 0)
	toPos, errTo := strconv.ParseUint(to, 10, 0)
	postId, errPost := strconv.ParseUint(postStr, 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	if errFrom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from dec from="))
		return
	}
	if errTo != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from dec to="))
		return
	}
	//Получение постов и комментов
	postRaw, commentsRaw, test := getPost(postId, fromPos, toPos)
	//Если поста нет, то возвращается 404
	switch test {
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("No such post %s", postStr)))
		return
	}
	//Маршализация ответа с ошибками
	post, errPost := json.Marshal(postRaw)
	comments, errCom := json.Marshal(commentsRaw)
	if errPost != nil {
		log(fmt.Sprintf("post wasn't parsed: %s", errPost.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
		return
	}
	if errCom != nil {
		log(fmt.Sprintf("comment wasn't parsed: %s", errCom.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(slices.Concat(post, []byte(","), comments)))
}

// Добавление поста
// Соответствует методу POST /Post?user_id=$user$&post_id=$post$ и POST /Post?user_id=$user$
// Возвращает ошибку, если user_id не соответствует user_id поста. Наивная реализация.
func post_post(
	w http.ResponseWriter,
	reader io.ReadCloser,
	createPost func(user uint64, message *string) (uint64, tdao.PROBLEM),
	updatePost func(post uint64, user uint64, message *string) tdao.PROBLEM,
	post string,
	user string,
) {
	//Проверяет, нужно ли создавать нужный пост
	createNew := post == ""
	//Парсит uid
	userId, errUser := strconv.ParseUint(user, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	//Читает из Body
	bytes, err := io.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot read body"))
		return
	}
	str := string(bytes)
	//Если поле post_id пустое, создаётся новый пост и возвращается его ID
	if createNew {
		retId, problem := createPost(userId, &str)
		switch problem {
		case tdao.NO_SUCH_USER:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("no such user"))
			return
		case tdao.INCORRECT_USER:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			id_str := strconv.FormatUint(retId, 16)
			w.Write([]byte(fmt.Sprintf("\"post_id\":\"%s\"", id_str)))
			return
		}
	}
	//Парсит не пустое поле post_id
	postId, errPost := strconv.ParseUint(post, 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	//Обновляет содержимое
	problem := updatePost(postId, userId, &str)
	switch problem {
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Post not found"))
		return
	case tdao.NO_SUCH_USER:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("no such user"))
		return
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

// Удаление поста
// Соответствует методу DELETE /Post?user_id=$user$&post_id=$post$
// Возвращает ошибку, если user_id не соответствует user_id поста. Наивная реализация.
func post_delete(
	w http.ResponseWriter,
	deletePost func(post uint64, user uint64) tdao.PROBLEM,
	user string,
	post string,
) {
	//Парсит user_id и post_id
	userId, errUser := strconv.ParseUint(user, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	postId, errPost := strconv.ParseUint(post, 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	//Удаляет пост
	ret := deletePost(postId, userId)
	switch ret {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	case tdao.NO_SUCH_USER:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such user"))
		return
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such post"))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

// Анализирует запросы к /Post
func Post(
	w http.ResponseWriter,
	r *http.Request,
	getPost func(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM),
	log func(string),
	createPost func(user uint64, message *string) (uint64, tdao.PROBLEM),
	updatePost func(post uint64, user uint64, message *string) tdao.PROBLEM,
	deletePost func(post uint64, user uint64) tdao.PROBLEM,
) {
	method := r.Method
	header := r.Header
	post := header.Get("post_id")
	switch method {
	case http.MethodGet:
		from := header.Get("from")
		to := header.Get("to")
		post_get(w, getPost, log, post, from, to)
	case http.MethodPost:
		user := header.Get("user_id")
		reader := r.Body
		defer reader.Close()
		post_post(w, reader, createPost, updatePost, post, user)
	case http.MethodDelete:
		user := header.Get("user_id")
		post_delete(w, deletePost, user, post)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func PostMute(
	w http.ResponseWriter,
	r *http.Request,
	mutePost func(post uint64, user uint64) tdao.PROBLEM,
) {
	header := r.Header
	userId, errUser := strconv.ParseUint(header.Get("user_id"), 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	postId, errPost := strconv.ParseUint(header.Get("post_id"), 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := mutePost(postId, userId)
	switch problem {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	case tdao.NO_SUCH_USER:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such user"))
		return
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such post"))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func PostUnmute(
	w http.ResponseWriter,
	r *http.Request,
	unmutePost func(post uint64, user uint64) tdao.PROBLEM,
) {
	header := r.Header
	userId, errUser := strconv.ParseUint(header.Get("user_id"), 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	postId, errPost := strconv.ParseUint(header.Get("post_id"), 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := unmutePost(postId, userId)
	switch problem {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	case tdao.NO_SUCH_USER:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such user"))
		return
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("No such post"))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
		return
	}
}
