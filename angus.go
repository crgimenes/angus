package angus

import (
	"crypto/rand"
	"net/http"
	"sync"
	"time"
)

const (
	cookieName = "angus"
)

var (
	mux *http.ServeMux  = http.NewServeMux()
	sc  *SessionControl = NewSession(cookieName)
	mx  sync.Mutex
)

func RandomID() string {
	const (
		length  = 16
		charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	)
	lenCharset := byte(len(charset))
	b := make([]byte, length)
	_, _ = rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = charset[b[i]%lenCharset]
	}
	return string(b)
}

func ServeMux() *http.ServeMux {
	return mux
}

func ListenAndServe(addr string) error {

	mux.HandleFunc("/angus/ws", wsHandler)
	mux.HandleFunc("/angus/", mainHandler)

	s := &http.Server{
		Handler:        mux,
		Addr:           addr,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.ListenAndServe()
}
