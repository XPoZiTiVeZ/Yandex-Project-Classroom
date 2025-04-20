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

type StudentTask struct {
	TaskID      string `json:"task_id"`
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskStatus struct {
	TaskID    string `json:"task_id"`
	StudentID string `json:"student_id"`
	Completed bool   `json:"completed"`
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
	CourseID string `json:"course_id"`
	TaskID string `json:"task_id"`
}

func NewGetTaskRequest(req GetTaskRequest) *pb.GetTaskRequest {
	return &pb.GetTaskRequest{
		TaskId: req.TaskID,
	}
}

type GetTaskResponse struct {
	Task Task `json:"task"`
}

func NewGetTaskResponse(resp *pb.GetTaskResponse) GetTaskResponse {
	return GetTaskResponse{
		Task: Task{
			TaskID: resp.GetTask().GetTaskId(),
		},
	}
}

type GetStudentStatusesRequest struct {
	CourseID string `json:"course_id"`
	TaskID   string `json:"task_id"`
}

func NewGetStudentStatusesRequest(req GetStudentStatusesRequest) *pb.GetStudentStatusesRequest {
	return &pb.GetStudentStatusesRequest{
		TaskId: req.TaskID,
	}
}

type GetStudentStatusesResponse struct {
	Statuses []TaskStatus `json:"statuses"`
}

func NewGetStudentStatusesResponse(resp *pb.GetStudentStatusesResponse) GetStudentStatusesResponse {
	return GetStudentStatusesResponse{
		Statuses: func() []TaskStatus {
			var statuses []TaskStatus
			for _, status := range resp.GetStatuses() {
				statuses = append(statuses, TaskStatus{
					TaskID:    status.GetTaskId(),
					StudentID: status.GetStudentId(),
					Completed: status.GetCompleted(),
				})
			}
			return statuses
		}(),
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
	Tasks []Task `json:"tasks"`
}

func NewGetTasksResponse(resp *pb.GetTasksResponse) GetTasksResponse {
	return GetTasksResponse{
		Tasks: func() []Task {
			var tasks []Task
			for _, task := range resp.GetTasks() {		
				tasks = append(tasks, Task{
					TaskID:      task.GetTaskId(),
					CourseID:    task.GetCourseId(),
					Title:       task.GetTitle(),
					Description: task.GetContent(),
					CreatedAt:   task.GetCreatedAt().AsTime(),
				})
			}
			return tasks
		}(),
	}
}

type GetTasksForStudentRequest struct {
	CourseID  string `json:"course_id"`
	StudentID string `json:"student_id"`
}

func NewGetTasksForStudentRequest(req GetTasksForStudentRequest) *pb.GetTasksForStudentRequest {
	return &pb.GetTasksForStudentRequest{
		CourseId:  req.CourseID,
		StudentId: req.StudentID,
	}
}

type GetTasksForStudentResponse struct {
	Tasks []StudentTask `json:"tasks"`
}

func NewGetTasksForStudentResponse(resp *pb.GetTasksForStudentResponse) GetTasksForStudentResponse {
	return GetTasksForStudentResponse{
		Tasks: func() []StudentTask {
			tasks := make([]StudentTask, len(resp.GetTasks()))
			for _, task := range resp.GetTasks() {
				tasks = append(tasks, StudentTask{
					TaskID:      task.GetTaskId(),
					CourseID:    task.GetCourseId(),
					Title:       task.GetTitle(),
					Description: task.GetContent(),
					Completed:   task.GetCompleted(),
					CreatedAt:   task.GetCreatedAt().AsTime(),
				})
			}
			return tasks
		}(),
	}
}

type UpdateTaskRequest struct {
	CourseID string `json:"course_id"`
	TaskID   string `json:"task_id"`
	Title   *string `json:"title,omitempty"`
	Content *string `json:"description,omitempty"`
}

func NewUpdateTaskRequest(req UpdateTaskRequest) *pb.UpdateTaskRequest {
	return &pb.UpdateTaskRequest{
		TaskId:  req.TaskID,
		Title:   req.Title,
		Content: req.Content,
	}
}

type UpdateTaskResponse struct {
}

func NewUpdateTaskResponse(resp *pb.UpdateTaskResponse) UpdateTaskResponse {
	return UpdateTaskResponse{}
}

type ChangeStatusTaskRequest struct {
	CourseID string `json:"course_id"`
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
	CourseID string `json:"course_id"`
	TaskID string `json:"task_id"`
}

func NewDeleteTaskRequest(req DeleteTaskRequest) *pb.DeleteTaskRequest {
	return &pb.DeleteTaskRequest{
		TaskId: req.TaskID,
	}
}

type DeleteTaskResponse struct {
}

func NewDeleteTaskResponse(resp *pb.DeleteTaskResponse) DeleteTaskResponse {
	return DeleteTaskResponse{}
}
