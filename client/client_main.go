package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/EgMeln/CRUDentity/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	conn, err := grpc.Dial(":50005", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("can't connect: %s", err)
	}
	client := protocol.NewUserServiceClient(conn)
	client2 := protocol.NewParkingServiceClient(conn)
	signInUserRequest := protocol.SignInUserRequest{Username: "egor10", Password: "1234"}
	log.Infof(signInUserRequest.String())

	tokens, err := client.SignIn(context.Background(), &signInUserRequest)
	if err != nil {
		log.Warnf("can't sign in %v", err)
	}
	log.Infof(tokens.String())
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"authorization": "Bearer " + tokens.Access}))

	user := protocol.GetUserRequest{Username: "egor10"}
	users, err := client.GetUser(ctx, &user)
	if err != nil {
		log.Warnf("%v", err)
	}
	log.Info(users)
	parkingLot := protocol.GetParkingLotRequest{Num: 1234}
	parkingLots, err := client2.GetParkingLot(ctx, &parkingLot)
	if err != nil {
		log.Warnf("%v", err)
	}
	log.Info(parkingLots)

	UploadImage(client, "img2.png")
}

func UploadImage(client protocol.UserServiceClient, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	stream, err := client.UploadImage(ctx)
	if err != nil {
		cancel()
		log.Fatal("cannot upload image: ", err)
	}

	req := &protocol.UploadImageRequest{
		Data: &protocol.UploadImageRequest_Info{
			Info: &protocol.ImageInfo{
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, ok := reader.Read(buffer)
		if ok == io.EOF {
			break
		}
		if ok != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &protocol.UploadImageRequest{
			Data: &protocol.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		ok = stream.Send(req)
		if ok != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded size: %d", res.GetSize())
	err = file.Close()
	if err != nil {
		log.Fatal("cannot file close: ", err)

	}
}
