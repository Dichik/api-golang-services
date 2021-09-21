package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func getHandlerCake(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cake"))
}

func main() {
	r := mux.NewRouter()

	users := NewInMemoryUserStorage()
	userService := UserService{
		repository: users,
	}
	r.HandleFunc("/cake", logRequest(getHandlerCake)).Methods(http.MethodGet)
	r.HandleFunc("/user/register", logRequest(userService.Register)).Methods(http.MethodPost)
	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	log.Println("Server started")
	err := srv.ListenAndServe()
	if err != nil {
		log.Println("error =", err)
	}
	log.Println("Bye :)")
}
