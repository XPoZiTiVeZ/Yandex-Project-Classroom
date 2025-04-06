package courses

import (
	pb "Classroom/Gateway/pkg/api/course"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
)

type CourseServiceClient struct {
	Client *pb.CourseServiceClient
}

func NewCoursesServiceClient(address string, port int) (*CourseServiceClient, error) {
	var opts []grpc.DialOption

	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", address, port), opts...)
	if err != nil {
		slog.Error("fail to dial: %v", slog.Any("error", err))
		return nil, err
	}
	defer conn.Close()

	client := pb.NewCourseServiceClient(conn)

	return &CourseServiceClient{
		Client: &client,
	}, nil
}
