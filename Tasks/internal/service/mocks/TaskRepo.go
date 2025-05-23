// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package service

import (
	"Classroom/Tasks/internal/domain"
	"Classroom/Tasks/internal/dto"
	"context"

	mock "github.com/stretchr/testify/mock"
)

// NewMockTaskRepo creates a new instance of MockTaskRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockTaskRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockTaskRepo {
	mock := &MockTaskRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// MockTaskRepo is an autogenerated mock type for the TaskRepo type
type MockTaskRepo struct {
	mock.Mock
}

type MockTaskRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *MockTaskRepo) EXPECT() *MockTaskRepo_Expecter {
	return &MockTaskRepo_Expecter{mock: &_m.Mock}
}

// CourseExists provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) CourseExists(ctx context.Context, courseID string) (bool, error) {
	ret := _mock.Called(ctx, courseID)

	if len(ret) == 0 {
		panic("no return value specified for CourseExists")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return returnFunc(ctx, courseID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = returnFunc(ctx, courseID)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, courseID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTaskRepo_CourseExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CourseExists'
type MockTaskRepo_CourseExists_Call struct {
	*mock.Call
}

// CourseExists is a helper method to define mock.On call
//   - ctx
//   - courseID
func (_e *MockTaskRepo_Expecter) CourseExists(ctx interface{}, courseID interface{}) *MockTaskRepo_CourseExists_Call {
	return &MockTaskRepo_CourseExists_Call{Call: _e.mock.On("CourseExists", ctx, courseID)}
}

func (_c *MockTaskRepo_CourseExists_Call) Run(run func(ctx context.Context, courseID string)) *MockTaskRepo_CourseExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTaskRepo_CourseExists_Call) Return(b bool, err error) *MockTaskRepo_CourseExists_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *MockTaskRepo_CourseExists_Call) RunAndReturn(run func(ctx context.Context, courseID string) (bool, error)) *MockTaskRepo_CourseExists_Call {
	_c.Call.Return(run)
	return _c
}

// Create provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) Create(ctx context.Context, payload dto.CreateTaskDTO) (domain.Task, error) {
	ret := _mock.Called(ctx, payload)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 domain.Task
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, dto.CreateTaskDTO) (domain.Task, error)); ok {
		return returnFunc(ctx, payload)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, dto.CreateTaskDTO) domain.Task); ok {
		r0 = returnFunc(ctx, payload)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, dto.CreateTaskDTO) error); ok {
		r1 = returnFunc(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTaskRepo_Create_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Create'
type MockTaskRepo_Create_Call struct {
	*mock.Call
}

// Create is a helper method to define mock.On call
//   - ctx
//   - payload
func (_e *MockTaskRepo_Expecter) Create(ctx interface{}, payload interface{}) *MockTaskRepo_Create_Call {
	return &MockTaskRepo_Create_Call{Call: _e.mock.On("Create", ctx, payload)}
}

func (_c *MockTaskRepo_Create_Call) Run(run func(ctx context.Context, payload dto.CreateTaskDTO)) *MockTaskRepo_Create_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(dto.CreateTaskDTO))
	})
	return _c
}

func (_c *MockTaskRepo_Create_Call) Return(task domain.Task, err error) *MockTaskRepo_Create_Call {
	_c.Call.Return(task, err)
	return _c
}

func (_c *MockTaskRepo_Create_Call) RunAndReturn(run func(ctx context.Context, payload dto.CreateTaskDTO) (domain.Task, error)) *MockTaskRepo_Create_Call {
	_c.Call.Return(run)
	return _c
}

// Delete provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) Delete(ctx context.Context, id string) error {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = returnFunc(ctx, id)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockTaskRepo_Delete_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Delete'
type MockTaskRepo_Delete_Call struct {
	*mock.Call
}

// Delete is a helper method to define mock.On call
//   - ctx
//   - id
func (_e *MockTaskRepo_Expecter) Delete(ctx interface{}, id interface{}) *MockTaskRepo_Delete_Call {
	return &MockTaskRepo_Delete_Call{Call: _e.mock.On("Delete", ctx, id)}
}

func (_c *MockTaskRepo_Delete_Call) Run(run func(ctx context.Context, id string)) *MockTaskRepo_Delete_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTaskRepo_Delete_Call) Return(err error) *MockTaskRepo_Delete_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockTaskRepo_Delete_Call) RunAndReturn(run func(ctx context.Context, id string) error) *MockTaskRepo_Delete_Call {
	_c.Call.Return(run)
	return _c
}

// GetByID provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) GetByID(ctx context.Context, id string) (domain.Task, error) {
	ret := _mock.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 domain.Task
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (domain.Task, error)); ok {
		return returnFunc(ctx, id)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) domain.Task); ok {
		r0 = returnFunc(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, id)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTaskRepo_GetByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetByID'
type MockTaskRepo_GetByID_Call struct {
	*mock.Call
}

// GetByID is a helper method to define mock.On call
//   - ctx
//   - id
func (_e *MockTaskRepo_Expecter) GetByID(ctx interface{}, id interface{}) *MockTaskRepo_GetByID_Call {
	return &MockTaskRepo_GetByID_Call{Call: _e.mock.On("GetByID", ctx, id)}
}

func (_c *MockTaskRepo_GetByID_Call) Run(run func(ctx context.Context, id string)) *MockTaskRepo_GetByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTaskRepo_GetByID_Call) Return(task domain.Task, err error) *MockTaskRepo_GetByID_Call {
	_c.Call.Return(task, err)
	return _c
}

func (_c *MockTaskRepo_GetByID_Call) RunAndReturn(run func(ctx context.Context, id string) (domain.Task, error)) *MockTaskRepo_GetByID_Call {
	_c.Call.Return(run)
	return _c
}

// ListByCourseID provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) ListByCourseID(ctx context.Context, courseID string) ([]domain.Task, error) {
	ret := _mock.Called(ctx, courseID)

	if len(ret) == 0 {
		panic("no return value specified for ListByCourseID")
	}

	var r0 []domain.Task
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) ([]domain.Task, error)); ok {
		return returnFunc(ctx, courseID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) []domain.Task); ok {
		r0 = returnFunc(ctx, courseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Task)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, courseID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTaskRepo_ListByCourseID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByCourseID'
type MockTaskRepo_ListByCourseID_Call struct {
	*mock.Call
}

// ListByCourseID is a helper method to define mock.On call
//   - ctx
//   - courseID
func (_e *MockTaskRepo_Expecter) ListByCourseID(ctx interface{}, courseID interface{}) *MockTaskRepo_ListByCourseID_Call {
	return &MockTaskRepo_ListByCourseID_Call{Call: _e.mock.On("ListByCourseID", ctx, courseID)}
}

func (_c *MockTaskRepo_ListByCourseID_Call) Run(run func(ctx context.Context, courseID string)) *MockTaskRepo_ListByCourseID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *MockTaskRepo_ListByCourseID_Call) Return(tasks []domain.Task, err error) *MockTaskRepo_ListByCourseID_Call {
	_c.Call.Return(tasks, err)
	return _c
}

func (_c *MockTaskRepo_ListByCourseID_Call) RunAndReturn(run func(ctx context.Context, courseID string) ([]domain.Task, error)) *MockTaskRepo_ListByCourseID_Call {
	_c.Call.Return(run)
	return _c
}

// ListByStudentID provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) ListByStudentID(ctx context.Context, studentID string, courseID string) ([]domain.StudentTask, error) {
	ret := _mock.Called(ctx, studentID, courseID)

	if len(ret) == 0 {
		panic("no return value specified for ListByStudentID")
	}

	var r0 []domain.StudentTask
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) ([]domain.StudentTask, error)); ok {
		return returnFunc(ctx, studentID, courseID)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string, string) []domain.StudentTask); ok {
		r0 = returnFunc(ctx, studentID, courseID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.StudentTask)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = returnFunc(ctx, studentID, courseID)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// MockTaskRepo_ListByStudentID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListByStudentID'
type MockTaskRepo_ListByStudentID_Call struct {
	*mock.Call
}

// ListByStudentID is a helper method to define mock.On call
//   - ctx
//   - studentID
//   - courseID
func (_e *MockTaskRepo_Expecter) ListByStudentID(ctx interface{}, studentID interface{}, courseID interface{}) *MockTaskRepo_ListByStudentID_Call {
	return &MockTaskRepo_ListByStudentID_Call{Call: _e.mock.On("ListByStudentID", ctx, studentID, courseID)}
}

func (_c *MockTaskRepo_ListByStudentID_Call) Run(run func(ctx context.Context, studentID string, courseID string)) *MockTaskRepo_ListByStudentID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *MockTaskRepo_ListByStudentID_Call) Return(studentTasks []domain.StudentTask, err error) *MockTaskRepo_ListByStudentID_Call {
	_c.Call.Return(studentTasks, err)
	return _c
}

func (_c *MockTaskRepo_ListByStudentID_Call) RunAndReturn(run func(ctx context.Context, studentID string, courseID string) ([]domain.StudentTask, error)) *MockTaskRepo_ListByStudentID_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function for the type MockTaskRepo
func (_mock *MockTaskRepo) Update(ctx context.Context, task domain.Task) error {
	ret := _mock.Called(ctx, task)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, domain.Task) error); ok {
		r0 = returnFunc(ctx, task)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// MockTaskRepo_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type MockTaskRepo_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx
//   - task
func (_e *MockTaskRepo_Expecter) Update(ctx interface{}, task interface{}) *MockTaskRepo_Update_Call {
	return &MockTaskRepo_Update_Call{Call: _e.mock.On("Update", ctx, task)}
}

func (_c *MockTaskRepo_Update_Call) Run(run func(ctx context.Context, task domain.Task)) *MockTaskRepo_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(domain.Task))
	})
	return _c
}

func (_c *MockTaskRepo_Update_Call) Return(err error) *MockTaskRepo_Update_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *MockTaskRepo_Update_Call) RunAndReturn(run func(ctx context.Context, task domain.Task) error) *MockTaskRepo_Update_Call {
	_c.Call.Return(run)
	return _c
}
