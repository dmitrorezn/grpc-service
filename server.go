package main 

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"encoding/json"
)

func main()  {

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		w.Write([]byte("\n back"))
	})

	r.Route("/api", func(api chi.Router) {

		api.Get("/article/{articleID}", articleGet)

		api.Post("/save/article", articlePost)

	})
	
	fmt.Println("Start APP server")

	http.ListenAndServe(":8080", r)
}

type data struct {
	ID string `json:"id"`
}

func articleGet(w http.ResponseWriter, r *http.Request) {

	articleID := chi.URLParam(r, "articleID")


	d := data{
		ID: articleID,
	}

	repsonce, err := json.Marshal(d)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(repsonce)
}

func articlePost(w http.ResponseWriter, r *http.Request) {


}