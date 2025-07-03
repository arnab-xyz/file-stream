package main

import (
	"fmt"
	pb "github.com/arnab-xyz/file-stream/protobuff"
	"io"
	"net"

	"google.golang.org/grpc"
)

type FileStreamServer struct {
	pb.UnimplementedFileStreamServer
}

func (f *FileStreamServer) Stream(stream pb.FileStream_StreamServer) error {
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return stream.SendAndClose(&pb.Response{
				Message: "Upload Failed",
				Success: true,
			})
		}
		fmt.Printf("Received data, size: %d, content: %s\n", chunk.GetSize(), string(chunk.GetData()[:chunk.GetSize()]))
	}
	return stream.SendAndClose(&pb.Response{
		Message: "Upload Successfully",
		Success: true,
	})
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err.Error())
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFileStreamServer(grpcServer, &FileStreamServer{})
	fmt.Println("Grpc server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		panic(err.Error())
	}
}
