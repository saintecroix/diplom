package main

import (
	"context"
	pb "github.com/saintecroix/diplom/cmd/inputConvert/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewInputConvertServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.ConvertExcelData(ctx, &pb.ConvertExcelDataRequest{FilePath: "114_09032025_1706_empty.xlsx"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Результат: %s", r.GetResults())
}
