package onlyhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strconv"

	tdao "../translationdao"
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
	//Если пользователей нет, то возвращается 404
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
		log("post wasn't parsed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
		return
	}
	if errCom != nil {
		log("comment wasn't parsed")
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
	bytes := make([]byte, 0)
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
	_, err := reader.Read(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot read body"))
		return
	}
	str := string(bytes)
	//Если поле post_id пустое, создаётся новый пост и возвращается его ID
	if createNew {
		retId, _ := createPost(userId, &str)
		id_str := strconv.FormatUint(retId, 16)
		w.Write([]byte(fmt.Sprint("\"post_id\":\"%s\"", id_str)))
		return
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
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
	}
	w.WriteHeader(http.StatusAccepted)
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
	//Удаляет
	ret := deletePost(postId, userId)
	switch ret {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
	}
	w.WriteHeader(http.StatusAccepted)
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
		message := r.Body
		post_post(w, message, createPost, updatePost, post, user)
	case http.MethodDelete:
		user := header.Get("user_id")
		post_delete(w, deletePost, user, post)
	}
}
