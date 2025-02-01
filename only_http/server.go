package onlyhttp

import (
	"net/http"
)

func StartServer(config ServerCallbacks) {
	mux := http.NewServeMux()
	http.ListenAndServe(":8000", mux)
	mux.HandleFunc("/test", config.Test)
	mux.HandleFunc("/post", config.Post)
	mux.HandleFunc("/post/mute", config.Mute)
	mux.HandleFunc("/post/unmute", config.Unmute)
	mux.HandleFunc("/posts", config.Posts)
	mux.HandleFunc("/comment", config.Comment)
}
