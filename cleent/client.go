package main

import (
	service "github.com/dmitrorezn/grpc-service/gen/service/proto"
	"google.golang.org/grpc"

	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var grpcPort = flag.String("grpc_port", "9094", "p2")

func main() {

	flag.Parse()

	conn, err := grpc.Dial("localhost:"+*grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Dial: err =", err)

		return
	}

	fmt.Println("Dial: connectd port", *grpcPort)

	defer conn.Close()

	client = service.NewArticleClient(conn)

	fmt.Println("NewArticleClient: connectd ")

	go ping()

	save()
}

var client service.ArticleClient

var ids []string

func get() []string {
	return ids
}

func ping() {
	fmt.Println("ping: start ")

	r := &service.GetArticleRequest{}

	for {
		i := 0
		for Ids := get(); i <= len(Ids) ; i++ {
			id := Ids[i]

			r.Id = id

			resp, err := client.GetArticleByID(context.Background(), r)
			if err != nil {
				fmt.Println("GetArticleByID: err =", err)
				continue
			}
			fmt.Println("GetArticleByID", resp)

			time.Sleep(time.Second * 2)
		}
	}
}

func save() {
	fmt.Println("save: start ")

	a := &service.Article{}

	a.Type = "new"

	setStream, err := client.SetArticles(context.Background())
	if err != nil {
		fmt.Println("SetArticles: err =", err)
		return
	}

	count := 0

	for {
		a.Text = "Test" + fmt.Sprint(count)

		err = setStream.Send(a)
		if err != nil {
			fmt.Println("setStream.Send: err =", err)
			continue
		}
		count++
		if count > 10 {
			count = 0
			resp, err := setStream.CloseAndRecv()
			if err != nil {
				fmt.Println("setStream.Send: err =", err)
			}

			ids = append(ids, resp.Id...)

			fmt.Println("SAVE.IDS: =", resp.Id, ids)

			time.Sleep(time.Second * 10)

			setStream, err = client.SetArticles(context.Background())
		}
	}
}
