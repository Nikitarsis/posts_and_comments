package onlyhttp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

// Организует получение комментария по Id
// Соответствует методу GET /Comment?comment_id=$comment$
func comment_get(
	w http.ResponseWriter,
	comment string,
	getComment func(comment uint64) (tdao.PROBLEM, *tdao.CommentDao),
) {
	//Парсинг CommentID
	commentId, errCom := strconv.ParseUint(comment, 16, 64)
	if errCom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex comment_id"))
		return
	}
	//Получение коммента
	problem, result := getComment(commentId)
	switch problem {
	case tdao.NO_SUCH_POST:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	message, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write(message)
}

// Организует запись комментария по Id
// Соответствует методам POST /Comment?comment_id=$comment$&user_id=$user$&post_id=$post_id$ и POST /Comment?user_id=$user$&post_id=$post_id$
// Возвращает ошибку, если user_id не соответствует user_id поста. Наивная реализация.
func comment_post(
	w http.ResponseWriter,
	reader io.ReadCloser,
	createComment func(user uint64, parent uint64, message *string) (uint64, tdao.PROBLEM),
	updateComment func(comment uint64, user uint64, message *string) tdao.PROBLEM,
	comment string,
	post string,
	user string,
) {
	//Проверяет, нужно ли создавать комментарий
	createNew := comment == ""
	//Парсит id пользователя
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
	//Читает тело комментария из поста
	bytes, err := io.ReadAll(reader)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cannot read body"))
		return
	}
	str := string(bytes)
	//Если comment_id нулевой, пост добавляется
	if createNew {
		retId, problem := createComment(userId, postId, &str)
		switch problem {
		case tdao.NO_SUCH_USER:
			id_str := strconv.FormatUint(userId, 16)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("User with id %s not found", id_str)))
			return
		case tdao.NO_SUCH_POST:
			id_str := strconv.FormatUint(postId, 16)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("Post with id %s not found", id_str)))
			return
		default:
			id_str := strconv.FormatUint(retId, 16)
			w.Write([]byte(fmt.Sprintf("\"post_id\":\"%s\"", id_str)))
			return
		}
	}
	//Парсит id комментария
	commentId, errCom := strconv.ParseUint(comment, 16, 64)
	if errCom != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	//Обновление комментария
	problem := updateComment(commentId, userId, &str)
	switch problem {
	case tdao.INCORRECT_USER:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Incorrect user"))
		return
	case tdao.NO_SUCH_POST:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Post with id %s not found", id_str)))
		return
	case tdao.NO_SUCH_USER:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("No such user %s", id_str)))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// Удаляет комментарий
// Соответствует методу DELETE /Comment?comment_id=$comment$&user_id=$user$
// Возвращает ошибку, если user_id не соответствует user_id поста. Наивная реализация.
func comment_delete(
	w http.ResponseWriter,
	deleteComment func(comment uint64, user uint64) tdao.PROBLEM,
	user string,
	comment string,
) {
	//Парист userId и commentId
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
	//Вызывает функцию удаления
	problem := deleteComment(commentId, userId)
	switch problem {
	case tdao.NO_SUCH_POST:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Comment with id %s doesn't exist", id_str)))
		return
	case tdao.INCORRECT_USER:
		id_str := strconv.FormatUint(commentId, 16)
		w.WriteHeader(http.StatusNotFound)
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(fmt.Sprintf("Incorrect user %s", id_str)))
		return
	default:
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

// Анализирует запросы к /Comment
func Comment(
	w http.ResponseWriter,
	r *http.Request,
	getComment func(comment uint64) (tdao.PROBLEM, *tdao.CommentDao),
	createComment func(user uint64, parent uint64, message *string) (uint64, tdao.PROBLEM),
	updateComment func(comment uint64, user uint64, message *string) tdao.PROBLEM,
	deleteComment func(comment uint64, user uint64) tdao.PROBLEM,
) {
	method := r.Method
	header := r.Header
	comment := header.Get("comment_id")
	switch method {
	case http.MethodGet:
		comment_get(w, comment, getComment)
	case http.MethodPost:
		user := header.Get("user_id")
		post := header.Get("post_id")
		message := r.Body
		defer message.Close()
		comment_post(w, message, createComment, updateComment, comment, post, user)
	case http.MethodDelete:
		user := header.Get("user_id")
		comment_delete(w, deleteComment, user, comment)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
