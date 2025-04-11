package tasks

import (
	"time"

	pb "Classroom/Gateway/pkg/api/tasks"
)

type Task struct {
	TaskID      string    `json:"task_id"`
	CourseID    string    `json:"course_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTaskRequest struct {
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func NewCreateTaskRequest(req CreateTaskRequest) *pb.CreateTaskRequest {
	return &pb.CreateTaskRequest{
		CourseId:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
	}
}

type CreateTaskResponse struct {
	TaskID string `json:"task_id"`
}

func NewCreateTaskResponse(resp *pb.CreateTaskResponse) CreateTaskResponse {
	return CreateTaskResponse{
		TaskID: resp.GetTaskId(),
	}
}

type GetTaskRequest struct {
	TaskID string `json:"task_id"`
}

func NewGetTaskRequest(req GetTaskRequest) *pb.GetTaskRequest {
	return &pb.GetTaskRequest{
		TaskId: req.TaskID,
	}
}

type GetTaskResponse struct {
	Task *pb.Task `json:"task"`
}

func NewGetTaskResponse(resp *pb.GetTaskResponse) GetTaskResponse {
	return GetTaskResponse{
		Task: resp.GetTask(),
	}
}

type GetTasksRequest struct {
	CourseID string `json:"course_id"`
}

func NewGetTasksRequest(req GetTasksRequest) *pb.GetTasksRequest {
	return &pb.GetTasksRequest{
		CourseId: req.CourseID,
	}
}

type GetTasksResponse struct {
	Tasks []*pb.Task `json:"tasks"`
}

func NewGetTasksResponse(resp *pb.GetTasksResponse) GetTasksResponse {
	return GetTasksResponse{
		Tasks: resp.GetTasks(),
	}
}

type UpdateTaskRequest struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

func NewUpdateTaskRequest(req UpdateTaskRequest) *pb.UpdateTaskRequest {
	return &pb.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}
}

type UpdateTaskResponse struct {
	
}

func NewUpdateTaskResponse(resp *pb.UpdateTaskResponse) UpdateTaskResponse {
	return UpdateTaskResponse{}
}

type ChangeStatusTaskRequest struct {
	TaskID string `json:"task_id"`
}

func NewChangeStatusTaskRequest(req ChangeStatusTaskRequest) *pb.ChangeStatusTaskRequest {
	return &pb.ChangeStatusTaskRequest{
		TaskId: req.TaskID,
	}
}

type ChangeStatusTaskResponse struct {

}

func NewChangeStatusTaskResponse(resp *pb.ChangeStatusTaskResponse) ChangeStatusTaskResponse {
	return ChangeStatusTaskResponse{}
}

type DeleteTaskRequest struct {
	TaskID string `json:"task_id"`
}

func NewDeleteTaskRequest(req DeleteTaskRequest) *pb.DeleteTaskRequest {
	return &pb.DeleteTaskRequest{
		TaskId: req.TaskID,
	}
}

type DeleteTaskResponse struct{}

func NewDeleteTaskResponse(resp *pb.DeleteTaskResponse) DeleteTaskResponse {
	return DeleteTaskResponse{}
}