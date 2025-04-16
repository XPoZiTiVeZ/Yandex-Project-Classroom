package main

import (
	"Classroom/Courses/internal/config"
	repo "Classroom/Courses/internal/repo/postgres"
	service "Classroom/Courses/internal/service"
	pb "Classroom/Courses/pkg/api/courses"
	"Classroom/Courses/pkg/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	conf := config.MustNew()

	postgres, err := postgres.MustNew(conf.PostgresURL)
	defer postgres.Close()
	if err != nil {
		log.Fatalf("Connect database error: %v", err)
		return
	}

	courseRepo := repo.NewCourseRepo(postgres)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port))
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}

	grpcServer := grpc.NewServer()

	courseService := service.NewCoursesService(courseRepo)
	pb.RegisterCoursesServiceServer(grpcServer, courseService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
