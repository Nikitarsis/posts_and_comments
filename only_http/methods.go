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

// Структура, содержащая колбеки
type ServerCallbacks struct {
	//Логгер
	log func(string)
	//Функция, просматривающая посты
	listPosts func() []tdao.PostDao
	//Получает пост с комментариями и
	getPost func(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM)
	//Загружает новый пост
	createPost func(user uint64, message *string) (uint64, tdao.PROBLEM)
	//Обновляет пост, если uid не совпадают, ошибка
	updatePost func(post uint64, user uint64, message *string) tdao.PROBLEM
	//Удаляет пост, если uid не совпадают, ошибка
	deletePost func(post uint64, user uint64) tdao.PROBLEM
	//Запрещает добавлять комментарии
	mutePost func(post uint64, user uint64) tdao.PROBLEM
	//Разрешает добавлять комментарии
	unmutePost func(post uint64, user uint64) tdao.PROBLEM
	//Получает комментарий
	getComment func(comment uint64) (tdao.PROBLEM, *tdao.CommentDao)
	//Загружает новый комментарий
	createComment func(user uint64, message *string) (uint64, tdao.PROBLEM)
	//Обновляет комментарий, если uid не совпадают, ошибка
	updateComment func(comment uint64, user uint64, message *string) tdao.PROBLEM
	//Удаляет комментарий, если uid не совпадают, ошибка
	deleteComment func(comment uint64, user uint64) tdao.PROBLEM
}

// Тестовый метод
func (s ServerCallbacks) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Hiiiiiii :3")
}

// Организует получение поста по Id
func (s ServerCallbacks) post_get(w http.ResponseWriter, postStr string, from string, to string) {
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
	postRaw, commentsRaw, test := s.getPost(postId, fromPos, toPos)
	//Если пользователей нет, то возвращается 404
	switch test {
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("No such post %s", postStr)))
		return
	}
	post, errPost := json.Marshal(postRaw)
	comments, errCom := json.Marshal(commentsRaw)
	if errPost != nil {
		s.log("post wasn't parsed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
		return
	}
	if errCom != nil {
		s.log("comment wasn't parsed")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(slices.Concat(post, []byte(","), comments)))
}

func (s ServerCallbacks) post_post(w http.ResponseWriter, reader io.ReadCloser, post string, user string) {
	bytes := make([]byte, 0)
	createNew := post == ""
	userId, errUser := strconv.ParseUint(user, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	_, err := reader.Read(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot read body"))
		return
	}
	str := string(bytes)
	if createNew {
		retId, _ := s.createPost(userId, &str)
		id_str := strconv.FormatUint(retId, 16)
		w.Write([]byte(fmt.Sprint("\"post_id\":\"%s\"", id_str)))
		return
	}
	postId, errPost := strconv.ParseUint(post, 16, 64)
	if errPost != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := s.updatePost(postId, userId, &str)
	switch problem {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s ServerCallbacks) post_delete(w http.ResponseWriter, user string, post string) {
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
	ret := s.deletePost(postId, userId)
	switch ret {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s ServerCallbacks) Post(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	header := r.Header
	post := header.Get("post_id")
	switch method {
	case http.MethodGet:
		from := header.Get("from")
		to := header.Get("to")
		s.post_get(w, post, from, to)
	case http.MethodPost:
		user := header.Get("user_id")
		message := r.Body
		s.post_post(w, message, post, user)
	case http.MethodDelete:
		user := header.Get("user_id")
		s.post_delete(w, user, post)
	}
}

func (s ServerCallbacks) comment_get(w http.ResponseWriter, comment string) {
	commentId, errCom := strconv.ParseUint(comment, 16, 64)
	if errCom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex comment_id"))
		return
	}
	problem, result := s.getComment(commentId)
	switch problem {
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		return
	case tdao.NO_PROBLEM:
		message, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write(message)
		return
	}
}

func (s ServerCallbacks) comment_post(w http.ResponseWriter, reader io.ReadCloser, comment string, user string) {
	bytes := make([]byte, 0)
	createNew := comment == ""
	userId, errUser := strconv.ParseUint(user, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	_, err := reader.Read(bytes)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot read body"))
		return
	}
	str := string(bytes)
	if createNew {
		retId, _ := s.createPost(userId, &str)
		id_str := strconv.FormatUint(retId, 16)
		w.Write([]byte(fmt.Sprint("\"post_id\":\"%s\"", id_str)))
		return
	}
	commentId, errCom := strconv.ParseUint(comment, 16, 64)
	if errCom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := s.updatePost(commentId, userId, &str)
	switch problem {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	case tdao.NO_SUCH_POST:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint("Post with id %s not found", id_str)))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s ServerCallbacks) comment_delete(w http.ResponseWriter, user string, comment string) {
	userId, errUser := strconv.ParseUint(user, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex user_id"))
		return
	}
	commentId, errCom := strconv.ParseUint(comment, 16, 64)
	if errCom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := s.deleteComment(commentId, userId)
	switch problem {
	case tdao.NO_SUCH_POST:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprint("Comment with id %s doesn't exist", id_str)))
		return
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("Incorrect User"))
	}
	w.WriteHeader(http.StatusAccepted)
}

func (s ServerCallbacks) Comment(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	header := r.Header
	comment := header.Get("comment_id")
	switch method {
	case http.MethodGet:
		s.comment_get(w, comment)
	case http.MethodPost:
		user := header.Get("user_id")
		message := r.Body
		s.comment_post(w, message, comment, user)
	case http.MethodDelete:
		user := header.Get("user_id")
		s.comment_delete(w, user, comment)
	}
}

func (s ServerCallbacks) Posts(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		ret, err := json.Marshal(s.listPosts())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(ret)
			return
		}
		w.Write(ret)
	}
}
