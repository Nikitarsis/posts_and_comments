package graphql

import (
	"net/http"

	"github.com/andrewwphillips/eggql"
)

func checkPosts(w http.ResponseWriter, r *http.Request) {
	eggql.InitialTimeout(1)
}

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	mux.HandleFunc("posts", checkPosts)

}
