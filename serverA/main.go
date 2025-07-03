package main

import (
	"fmt"
	pb "github.com/arnab-xyz/file-stream/protobuff"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"net/http"
)

func main() {

	server := gin.Default()
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	client := pb.NewFileStreamClient(conn)
	server.POST("/upload", func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		fmt.Printf("Uploaded file details: %s %d\n", fileHeader.Filename, fileHeader.Size)

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer file.Close()

		buffer := make([]byte, 32*1024)
		stream, err := client.Stream(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			chunk := &pb.File{
				Data: buffer[:n],
				Size: int32(n),
			}
			stream.Send(chunk)
		}
		_, err = stream.CloseAndRecv()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "File read successfully",
		})
	})

	server.Run("localhost:8080")

}
