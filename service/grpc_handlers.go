package main

import (
	"fmt"
	"io"
	service "github.com/dmitrorezn/grpc-service/gen/service/proto"
	"context"
	metadata "google.golang.org/grpc/metadata"
)

type articleServer struct {
	service.ArticleServer
}

func newServer() *articleServer {
	return &articleServer{}
}

var cc = map[string]string{"1":"test"}

func(s *articleServer) GetArticleByID(ctx context.Context, request *service.GetArticleRequest) (resp *service.ArticleResponce, err error) {
	h, _ := metadata.FromIncomingContext(ctx)

	fmt.Println("GetArticleByID: Metadata = ", h, " ID ",request.GetId())
	
	resp = &service.ArticleResponce{}

	resp.Article = &service.Article{}
	
	a := articles[fmt.Sprint(request.GetId())]

	resp.Article = &a

	err = fmt.Errorf("not creted")

	header := metadata.New(map[string]string{})

	header.Set("error", err.Error())

	ctx = metadata.NewOutgoingContext(ctx, header)

	return resp, err
}

var articles = map[string]service.Article{}
var ids int

func(s *articleServer) SetArticles(stream service.Article_SetArticlesServer) error {

	var response = &service.ArticlesFeature{}

	h, _ := metadata.FromIncomingContext(stream.Context())

	fmt.Println("SetArticles: Metadata = ", h)

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
