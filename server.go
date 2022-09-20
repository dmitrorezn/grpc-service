package main 

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"flag"

	"net"
	"context"

	service "github.com/dmitrorezn/grpc-service/gen/service"

	"google.golang.org/grpc"
)

// set PATH=%PATH%;C:\protoc-21.6-win64\bin;%GOPATH%/bin
// 
//protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  service/proto/service.proto

var (
	httpPort = flag.String("http_port","8080","p1")
	grpcPort = flag.String("grpc_port","9092","p2")
)

func main()  {

	flag.Parse()

	c := make(chan struct{})
	go httpServer(c)

	go grpcServer(c)

	<-c
}

type articleServer struct {
	service.ArticleServer
}

func newServer() *articleServer {
	return &articleServer{}
}

var cc = map[string]string{"1":"test"}

func(s *articleServer) GetArticleByID(ctx context.Context, request *service.GetArticleRequest) (resp *service.ArticleResponce,err error) {


	fmt.Println("GetArticleByID: ID = ",request.GetId())

	resp = &service.ArticleResponce{}

	resp.Article = &service.Article{}

	resp.Article.Id = request.GetId()

	resp.Article.Title = cc[request.GetId()]

	return resp, err
}

func grpcServer(c chan struct{}) {

	lis, err := net.Listen("tcp", ":"+*grpcPort)

	if err != nil {
		panic("grpc Lisener")
	}

	s := grpc.NewServer()

	service.RegisterArticleServer(s, newServer())

	fmt.Println("grpc: port =", *grpcPort)


	err = s.Serve(lis)

	if err != nil {
		c <- struct{}{}
		panic("grpc liServesener")
	}

	c <- struct{}{}
}


func httpServer(c chan struct{}) {
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

	fmt.Println("http: port =", *httpPort)


	http.ListenAndServe(":"+*httpPort, r)

	c <- struct{}{}
}


func articleGet(w http.ResponseWriter, r *http.Request) {

	articleID := chi.URLParam(r, "articleID")

	resp := &service.ArticleResponce{}

	resp.Article = &service.Article{}

	resp.Article.Id = articleID

	resp.Article.Title = cc[articleID]

	repsonce, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	w.Write(repsonce)
}

func articlePost(w http.ResponseWriter, r *http.Request) {


}