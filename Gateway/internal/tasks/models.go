package tasks

import (
	"time"

	pb "Classroom/Gateway/pkg/api/tasks"
)

// Task - основная информация о задании
// @Description Полная информация о задании в курсе
type Task struct {
    // Уникальный идентификатор задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Идентификатор курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Название задания
    Title string `json:"title" example:"Домашнее задание 1" extensions:"x-order=2"`
    // Описание задания
    Description string `json:"description" example:"Решить задачи по алгоритмам" extensions:"x-order=3"`
    // Дата создания
    CreatedAt time.Time `json:"created_at" example:"2023-01-15T10:00:00Z" extensions:"x-order=4"`
} // @name Task

// StudentTask - информация о задании для студента
// @Description Расширенная информация о задании с указанием статуса выполнения
type StudentTask struct {
    // Уникальный идентификатор задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Идентификатор курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Название задания
    Title string `json:"title" example:"Домашнее задание 1" extensions:"x-order=2"`
    // Описание задания
    Description string `json:"description" example:"Решить задачи по алгоритмам" extensions:"x-order=3"`
    // Статус выполнения
    Completed bool `json:"completed" example:"false" extensions:"x-order=4"`
    // Дата создания
    CreatedAt time.Time `json:"created_at" example:"2023-01-15T10:00:00Z" extensions:"x-order=5"`
} // @name StudentTask

// TaskStatus - статус выполнения задания
// @Description Информация о выполнении задания конкретным студентом
type TaskStatus struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID студента
    StudentID string `json:"student_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
    // Статус выполнения
    Completed bool `json:"completed" example:"true" extensions:"x-order=2"`
} // @name TaskStatus

// CreateTaskRequest - запрос на создание задания
// @Description Параметры для создания нового задания в курсе
type CreateTaskRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Название задания
    Title string `json:"title" example:"Лабораторная работа 1" extensions:"x-order=1"`
    // Описание задания
    Description string `json:"description" example:"Реализовать алгоритм сортировки" extensions:"x-order=2"`
} // @name CreateTaskRequest

func NewCreateTaskRequest(req CreateTaskRequest) *pb.CreateTaskRequest {
	return &pb.CreateTaskRequest{
		CourseId:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
	}
}

// CreateTaskResponse - ответ после создания задания
// @Description Возвращает ID созданного задания
type CreateTaskResponse struct {
    // ID созданного задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name CreateTaskResponse

func NewCreateTaskResponse(resp *pb.CreateTaskResponse) CreateTaskResponse {
	return CreateTaskResponse{
		TaskID: resp.GetTaskId(),
	}
}

// GetTaskRequest - запрос информации о задании
// @Description Требует ID курса и задания для получения данных
type GetTaskRequest struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetTaskRequest

func NewGetTaskRequest(req GetTaskRequest) *pb.GetTaskRequest {
	return &pb.GetTaskRequest{
		TaskId: req.TaskID,
	}
}

// GetTaskResponse - информация о задании
// @Description Возвращает полные данные задания
type GetTaskResponse struct {
    // Объект задания
    Task Task `json:"task" extensions:"x-order=0"`
} // @name GetTaskResponse

func NewGetTaskResponse(resp *pb.GetTaskResponse) GetTaskResponse {
	return GetTaskResponse{
		Task: Task{
			TaskID: resp.GetTask().GetTaskId(),
		},
	}
}

// GetStudentStatusesRequest - запрос статусов студентов
// @Description Возвращает статусы выполнения задания для студентов курса
type GetStudentStatusesRequest struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetStudentStatusesRequest

func NewGetStudentStatusesRequest(req GetStudentStatusesRequest) *pb.GetStudentStatusesRequest {
	return &pb.GetStudentStatusesRequest{
		TaskId: req.TaskID,
	}
}

// GetStudentStatusesResponse - статусы студентов
// @Description Содержит статусы выполнения задания студентами
type GetStudentStatusesResponse struct {
    // Массив статусов
    Statuses []TaskStatus `json:"statuses" extensions:"x-order=0"`
} // @name GetStudentStatusesResponse

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

// GetTasksRequest - запрос списка заданий
// @Description Возвращает все задания в указанном курсе
type GetTasksRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name GetTasksRequest

func NewGetTasksRequest(req GetTasksRequest) *pb.GetTasksRequest {
	return &pb.GetTasksRequest{
		CourseId: req.CourseID,
	}
}

// GetTasksResponse - список заданий
// @Description Содержит массив заданий в курсе
type GetTasksResponse struct {
    // Массив заданий
    Tasks []Task `json:"tasks" extensions:"x-order=0"`
} // @name GetTasksResponse

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

// GetTasksForStudentRequest - запрос заданий для студента
// @Description Возвращает задания с информацией о статусе выполнения
type GetTasksForStudentRequest struct {
    // ID курса
    CourseID string `json:"course_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // ID студента
    StudentID string `json:"student_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=1"`
} // @name GetTasksForStudentRequest

func NewGetTasksForStudentRequest(req GetTasksForStudentRequest) *pb.GetTasksForStudentRequest {
	return &pb.GetTasksForStudentRequest{
		CourseId:  req.CourseID,
		StudentId: req.StudentID,
	}
}

// GetTasksForStudentResponse - ответ с заданиями для студента
// @Description Содержит задания с информацией о статусе выполнения
type GetTasksForStudentResponse struct {
    // Массив заданий со статусами
    Tasks []StudentTask `json:"tasks" extensions:"x-order=0"`
} // @name GetTasksForStudentResponse

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

// UpdateTaskRequest - запрос обновления задания
// @Description Позволяет обновить данные задания
type UpdateTaskRequest struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
    // Новое название (опционально)
    Title *string `json:"title,omitempty" example:"Обновленное название" extensions:"x-order=1"`
    // Новое описание (опционально)
    Content *string `json:"description,omitempty" example:"Новые требования" extensions:"x-order=2"`
} // @name UpdateTaskRequest

func NewUpdateTaskRequest(req UpdateTaskRequest) *pb.UpdateTaskRequest {
	return &pb.UpdateTaskRequest{
		TaskId:  req.TaskID,
		Title:   req.Title,
		Content: req.Content,
	}
}

// UpdateTaskResponse - результат обновления
// @Description Пустой ответ при успешном обновлении
type UpdateTaskResponse struct {
} // @name UpdateTaskResponse

func NewUpdateTaskResponse(resp *pb.UpdateTaskResponse) UpdateTaskResponse {
	return UpdateTaskResponse{}
}

// ChangeStatusTaskRequest - запрос изменения статуса
// @Description Позволяет изменить статус выполнения задания
type ChangeStatusTaskRequest struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name ChangeStatusTaskRequest

func NewChangeStatusTaskRequest(req ChangeStatusTaskRequest) *pb.ChangeStatusTaskRequest {
	return &pb.ChangeStatusTaskRequest{
		TaskId: req.TaskID,
	}
}

// ChangeStatusTaskResponse - результат изменения статуса
// @Description Пустой ответ при успешном изменении
type ChangeStatusTaskResponse struct {
} // @name ChangeStatusTaskResponse

func NewChangeStatusTaskResponse(resp *pb.ChangeStatusTaskResponse) ChangeStatusTaskResponse {
	return ChangeStatusTaskResponse{}
}

// DeleteTaskRequest - запрос удаления задания
// @Description Требует ID курса и задания для удаления
type DeleteTaskRequest struct {
    // ID задания
    TaskID string `json:"task_id" example:"d277084b-e1f6-4670-825b-53951d20b5d3" extensions:"x-order=0"`
} // @name DeleteTaskRequest

func NewDeleteTaskRequest(req DeleteTaskRequest) *pb.DeleteTaskRequest {
	return &pb.DeleteTaskRequest{
		TaskId: req.TaskID,
	}
}

// DeleteTaskResponse - результат удаления
// @Description Пустой ответ при успешном удалении
type DeleteTaskResponse struct {
} // @name DeleteTaskResponse

func NewDeleteTaskResponse(resp *pb.DeleteTaskResponse) DeleteTaskResponse {
	return DeleteTaskResponse{}
}
