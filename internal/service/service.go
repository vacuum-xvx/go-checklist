package service

import "go-checklist-pet/internal/task"

type TaskRepository interface {
	SaveTask(task task.Task) task.Task
	List() []task.Task
	MarkDone(id int) (task.Task, bool)
	Delete(id int) bool
}

type TaskService struct {
	repo TaskRepository
}

func NewTaskService(repo TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(TaskDTO task.TaskDTO) task.Task {
	newTask := task.Task{
		TaskDTO: TaskDTO,
		Done:    false,
	}

	return s.repo.SaveTask(newTask)
}

func (s *TaskService) ListTasks() []task.Task {
	return s.repo.List()
}

func (s *TaskService) MarkTaskDone(id int) (task.Task, bool) {
	return s.repo.MarkDone(id)
}

func (s *TaskService) DeleteTask(id int) bool {
	return s.repo.Delete(id)
}
