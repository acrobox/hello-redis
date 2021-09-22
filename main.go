package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
)

func main() {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":" + os.Getenv("ABX_PORT"), Handler: mux}
	mux.HandleFunc("/", index)
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	n, err := incr()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	fmt.Fprintf(w, "%d\n", n)
}

func incr() (int64, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("ABX_CACHE_DSN")})
	defer rdb.Close()
	name := os.Getenv("ABX_NAME")
	return rdb.Incr(ctx, name).Result()
}
