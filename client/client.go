package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	pb "chat-app/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var serverUri string = "localhost:8080"
	const exitCommand string = "$$exit"

	var options []grpc.DialOption
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(serverUri, options...)

	if err != nil {
		fmt.Printf("Failed to connect to server: %v", err)
	}

	defer conn.Close()

	var clientId = os.Args[1]

	client := pb.NewChatServiceClient(conn)

	go func() {
		stream, err := client.Connect(context.Background(), &pb.Client{Id: clientId})
		if err != nil {
			fmt.Printf("Failed to connect to server: %v", err)
		}
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				fmt.Printf("Failed to receive message: %v", err)
			}
			fmt.Printf("%s: %s\n", message.Sender, message.Content)
		}
	}()

	// create a loop that asks for user input and sends it to the server as a message until the user types "exit"

	for {
		reader := bufio.NewReader(os.Stdin)
		content, _ := reader.ReadString('\n')
		content = content[:len(content)-1] // remove the newline character
		if content == exitCommand {
			break
		}
		client.SendMessage(context.Background(), &pb.Message{Sender: clientId, Content: content})
	}
}
