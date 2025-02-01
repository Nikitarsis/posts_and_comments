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
	warn          func(string)
	log           func(string)
	listPosts     func() []tdao.PostDao
	getPost       func(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, bool)
	postPost      func(user uint64, message *string) (uint64, tdao.PROBLEM)
	updatePost    func(post uint64, user uint64, message *string) tdao.PROBLEM
	deletePost    func(post uint64, user uint64) tdao.PROBLEM
	mutePost      func(post uint64, user uint64) tdao.PROBLEM
	unmutePost    func(post uint64, user uint64) tdao.PROBLEM
	getComment    func(comment uint64) tdao.PROBLEM
	postComment   func(user uint64, message *string) (uint64, tdao.PROBLEM)
	updateComment func(comment uint64, user uint64, message *string) tdao.PROBLEM
	deleteComment func(comment uint64, user uint64)
}

func (s ServerCallbacks) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
	fmt.Fprintf(w, "Hiiiiiii :3")
}

func (s ServerCallbacks) post_get(w http.ResponseWriter, postStr string, from string, to string) {
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
	postRaw, commentsRaw, check := s.getPost(postId, fromPos, toPos)
	if !check {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("No such post %d", strconv.FormatUint(postId, 16))))
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
		retId, response := s.postPost(userId, &str)
		id_str := strconv.FormatUint(retId, 16)
		w.Write([]byte(fmt.Sprint("\"post_id\":\"%s\"", id_str)))
		return
	}
	postId, errPost := strconv.ParseUint(postStr, 16, 64)
	if errUser != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse from hex post_id"))
		return
	}
	problem := s.updatePost(postId, userId, bytes)
	switch problem {
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
	case http.Method
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
