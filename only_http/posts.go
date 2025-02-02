package onlyhttp

import (
	"encoding/json"
	"net/http"

	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
)

func Posts(
	w http.ResponseWriter,
	r *http.Request,
	listPosts func() []tdao.PostDao,
) {
	method := r.Method
	if method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
		ret, err := json.Marshal(listPosts())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(ret)
			return
		}
		w.Write(ret)
	}
}
