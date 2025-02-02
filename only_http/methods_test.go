package onlyhttp

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"testing"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

type TestReaderCloser struct {
	InnerReader io.Reader
	finished    bool
}

func (t TestReaderCloser) Close() error {
	return nil
}

func (t TestReaderCloser) Read(p []byte) (n int, err error) {
	return t.InnerReader.Read(p)
}

func getTestReaderCloser(str string) TestReaderCloser {
	reader := strings.NewReader(str)
	return TestReaderCloser{
		InnerReader: reader,
		finished:    false,
	}
}

type TestResponse struct {
	out     func([]byte)
	counter func(int)
}

func (t TestResponse) Header() http.Header {
	return make(http.Header)
}

func (t *TestResponse) Write(bytes []byte) (int, error) {
	t.out(bytes)
	return 0, nil
}

func (t *TestResponse) WriteHeader(statusCode int) {
	t.counter(statusCode)
}

func getTestHttpWriter(out func([]byte), counter func(int)) http.ResponseWriter {
	return &TestResponse{out, counter}
}

func getTestPost() tdao.PostDao {
	uid := rand.Uint64()
	pid := rand.Uint64()
	str := "test"
	return tdao.PostDao{
		PostId:  pid,
		UserId:  uid,
		Message: &str,
	}
}

func getTestComment() []tdao.CommentDao {
	ret := make([]tdao.CommentDao, 5)
	for i := range ret {
		user := rand.Uint64()
		post := rand.Uint64()
		parent := rand.Uint64()
		children := make([]uint64, 5)
		for i := range children {
			children[i] = rand.Uint64()
		}
		str := "test"
		ret[i] = tdao.CommentDao{
			UserId:      user,
			PostId:      post,
			ParentId:    parent,
			ChildrenIds: children,
			Message:     &str,
		}
	}
	return ret
}

func getReadCloser(str string) io.ReadCloser {
	return getTestReaderCloser(str)
}

func TestGetPost(t *testing.T) {
	rand.New(rand.NewSource(int64(12234)))
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
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	w := getTestHttpWriter(out, counter)
	post_get(
		w,
		getPost,
		func(s string) { fmt.Println(s) },
		strconv.FormatInt(int64(postId), 16),
		strconv.FormatInt(int64(fromId), 10),
		strconv.FormatInt(int64(toId), 10),
	)
}

func TestPostPost(t *testing.T) {
	rand.New(rand.NewSource(int64(12234)))
	userId := rand.Uint64()
	postId := rand.Uint64()
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	w := getTestHttpWriter(out, counter)
	body := "Test string"
	reader := getReadCloser(body)
	createPost := func(user uint64, message *string) (uint64, tdao.PROBLEM) {
		if user != userId {
			t.Error()
		}
		if *message != body {
			t.Error()
		}
		return rand.Uint64(), tdao.NO_PROBLEM
	}
	updatePost := func(post uint64, user uint64, message *string) tdao.PROBLEM {
		if post != postId {
			fmt.Printf("arg:%d\nneed:%d\n", post, postId)
			t.Error()
		}
		if user != userId {
			fmt.Printf("arg:%d\nneed:%d\n", user, userId)
			t.Error()
		}
		if *message != body {
			fmt.Printf("arg:%s\nneed:%s\n", *message, body)
			t.Error()
		}
		return tdao.NO_PROBLEM
	}
	post_post(
		w,
		reader,
		createPost,
		updatePost,
		strconv.FormatUint(postId, 16),
		strconv.FormatUint(userId, 16),
	)
}

func TestPostDelete(t *testing.T) {
	rand.New(rand.NewSource(int64(12234)))
	userId := rand.Uint64()
	postId := rand.Uint64()
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	w := getTestHttpWriter(out, counter)
	deletePost := func(post uint64, user uint64) tdao.PROBLEM {
		if post != postId {
			t.Error()
		}
		if user != userId {
			t.Error()
		}
		return tdao.NO_PROBLEM
	}
	post_delete(
		w,
		deletePost,
		strconv.FormatUint(userId, 16),
		strconv.FormatUint(postId, 16),
	)
}

func TestPostsList(t *testing.T) {
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	check := false
	listPosts := func() []tdao.PostDao {
		check = true
		return make([]tdao.PostDao, 0)
	}
	w := getTestHttpWriter(out, counter)
	posts_get(w, listPosts)
	if !check {
		t.Error()
	}
}

func TestCommentGet(t *testing.T) {
	rand.New(rand.NewSource(int64(12234)))
	commentId := rand.Uint64()
	getComment := func(comment uint64) (tdao.PROBLEM, *tdao.CommentDao) {
		if comment != uint64(commentId) {
			t.Error()
		}
		retCom := getTestComment()
		problem := tdao.NO_PROBLEM
		return problem, &retCom[1]
	}
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	w := getTestHttpWriter(out, counter)
	comment_get(
		w,
		strconv.FormatUint(commentId, 16),
		getComment,
	)
}

func TestCommentPost(t *testing.T) {
	rand.New(rand.NewSource(int64(12234)))
	userId := rand.Uint64()
	postId := rand.Uint64()
	out := func(s []byte) {
		fmt.Print(string(s))
	}
	counter := func(c int) {
		if c/100 != 2 {
			t.Error(c)
		}
	}
	w := getTestHttpWriter(out, counter)
	body := "Test string"
	reader := getReadCloser(body)
	createPost := func(user uint64, message *string) (uint64, tdao.PROBLEM) {
		if user != userId {
			t.Error()
		}
		if *message != body {
			t.Error()
		}
		return rand.Uint64(), tdao.NO_PROBLEM
	}
	updatePost := func(post uint64, user uint64, message *string) tdao.PROBLEM {
		if post != postId {
			fmt.Printf("arg:%d\nneed:%d\n", post, postId)
			t.Error()
		}
		if user != userId {
			fmt.Printf("arg:%d\nneed:%d\n", user, userId)
			t.Error()
		}
		if *message != body {
			fmt.Printf("arg:%s\nneed:%s\n", *message, body)
			t.Error()
		}
		return tdao.NO_PROBLEM
	}
	comment_post(
		w,
		reader,
		createPost,
		updatePost,
		strconv.FormatUint(postId, 16),
		strconv.FormatUint(userId, 16),
	)
}
