package lessons

import (
	pb "Classroom/Gateway/pkg/api/lesson"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

type LessonServiceClient struct {
	Client *pb.LessonServiceClient
}

func NewLessonsServiceClient(address string, port int) (*LessonServiceClient, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		slog.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}
	defer conn.Close()

	client := pb.NewLessonServiceClient(conn)

	return &LessonServiceClient{
		Client: &client,
	}, nil
}
