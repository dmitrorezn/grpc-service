package main

import (
	service "github.com/dmitrorezn/grpc-service/gen/service/proto"
	"google.golang.org/grpc"

	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	metadata "google.golang.org/grpc/metadata"
	status "google.golang.org/grpc/status"
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
		Ids := get()
		for _, id := range Ids {

			r.Id = id

			header := metadata.New(map[string]string{})

			header.Set("client", "GetArticleByID")

			ctx := metadata.NewOutgoingContext(context.Background(), header)

			resp, err := client.GetArticleByID(ctx, r)
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					fmt.Println("FromError: Message =", err)
				}

				time.Sleep(time.Second)

				if msg := st.Message(); 	msg == "not creted" {
					fmt.Println("GetArticleByID: Message =", msg)
					continue
				}

				fmt.Println("GetArticleByID: err =", "m:",st.Message(),"c:",st.Code(),"err:",st.Err(),"dt:", st.Details() )
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

	header := metadata.New(map[string]string{})

	header.Set("client", "save")

	ctx := metadata.NewOutgoingContext(context.Background(), header)

	setStream, err := client.SetArticles(ctx)
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

			setStream, err = client.SetArticles(ctx)
		}
	}
}
