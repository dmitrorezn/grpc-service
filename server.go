package main 

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"flag"
	"tcp"
	service "github.com/dmitrorezn/grpc-service/gen"

)

// set PATH=%PATH%;C:\protoc-21.6-win64\bin;%GOPATH%/bin
// 
//protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  service/proto/service.proto

var (
	httpPort = flag.String("http_port","8080","p1")
	grpcPort = flag.String("grpc_port","9090","p2")
)

func main()  {

	flag.Parse()


	go httpServer()

	go grpcServer()

}

type articleServer {
	service.ArticleServer
}

func newServer() articleServer {
	return &articleServer{}
}

func(s *articleServer) GetArticleByID(ctx context.Context,request *service.GetArticleRequest) (*service.ArticleResponce, error) {

	return 
}

func grpcServer() {

	lis, err := tcp.Lisener("tcp", ":"+grpcPort)

	if err != nil {
		panic("grpc Lisener")
	}

	s := grpc.NewServer()


	service.RegisterArticleServer(s, newServer())

	err = s.Serve(lis)

	if err != nil {
		panic("grpc liServesener")
	}
}





func httpServer() {
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

	http.ListenAndServe(":"+httpPort, r)
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