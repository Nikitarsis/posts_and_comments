package onlyhttp

import (
	"net/http"
	"strconv"
	"testing"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

type TestResponse struct {
	out chan<- string
}

func (t TestResponse) Header() http.Header {
	return make(http.Header)
}

func (t *TestResponse) Write(bytes []byte) (int, error) {
	t.out <- string(bytes)
	return 0, nil
}

func (t *TestResponse) WriteHeader(statusCode int) {
	t.out <- strconv.FormatInt(int64(statusCode), 10)
}

func getTestHttpWriter(out chan<- string) http.ResponseWriter {
	return &TestResponse{out}
}

func getTestPost() tdao.PostDao {
	uid := uint64(111)
	pid := uint64(111)
	str := "test"
	return tdao.PostDao{pid, uid, &str}
}

func getTestComment() []tdao.CommentDao {
	ret := make([]tdao.CommentDao, 5)
	for i, _ := range ret {
		user := uint64(111 * i)
		post := uint64(111 * i)
		parent := uint64(111 * i)
		children := make([]uint64, 5)
		str := "test"
		ret[i] = tdao.CommentDao{user, post, parent, children, &str}
	}
	return ret
}

func TestGetPost(t *testing.T) {
	postId := 124
	fromId := 0
	toId := 100
	getPost := func(post uint64, from uint64, to uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM) {
		retPost := getTestPost()
		if post != uint64(postId) {
			t.Error()
		}
		if from != uint64(fromId) {
			t.Error()
		}
		if to != uint64(toId) {
			t.Error()
		}
		retCom := getTestComment()
		problem := tdao.NO_PROBLEM
		return retPost, retCom, problem
	}
	var out chan string = make(chan string)
	w := getTestHttpWriter(out)
	post_get(
		w,
		getPost,
		func(s string) {},
		strconv.FormatInt(int64(postId), 16),
		strconv.FormatInt(int64(fromId), 10),
		strconv.FormatInt(int64(toId), 10),
	)
}
