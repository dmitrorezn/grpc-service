package main

import (
	"fmt"
	"io"
	service "github.com/dmitrorezn/grpc-service/gen/service/proto"
	"context"
)

type articleServer struct {
	service.ArticleServer
}

func newServer() *articleServer {
	return &articleServer{}
}

var cc = map[string]string{"1":"test"}

func(s *articleServer) GetArticleByID(ctx context.Context, request *service.GetArticleRequest) (resp *service.ArticleResponce, err error) {


	fmt.Println("GetArticleByID: ID = ",request.GetId())

	resp = &service.ArticleResponce{}

	resp.Article = &service.Article{}
	
	a := articles[fmt.Sprint(request.GetId())]

	resp.Article = &a

	return resp, err
}

var articles = map[string]service.Article{}
var ids int

func(s *articleServer) SetArticles(stream service.Article_SetArticlesServer) error {

	var response = &service.ArticlesFeature{}

	for {
		article, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Recv: err = %w",err)
		}

		ids++

		article.Id = fmt.Sprint(ids)

		articles[fmt.Sprint(article.GetId())] = *article

		response.Id = append(response.Id, article.GetId())
		fmt.Println("article:", article)
	}

	return stream.SendAndClose(response)
}
