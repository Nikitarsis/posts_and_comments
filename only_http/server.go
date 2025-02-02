package onlyhttp

import (
	"net/http"
)

func StartServer(config ServerCallbacks) {
	mux := http.NewServeMux()
	http.ListenAndServe(":8000", mux)
	mux.HandleFunc("/test", config.Test)
	mux.HandleFunc("/post", config.Post)
	mux.HandleFunc("/post/mute", config.PostMute)
	mux.HandleFunc("/post/unmute", config.PostUnmute)
	mux.HandleFunc("/posts", config.Posts)
	mux.HandleFunc("/comment", config.Comment)
}
