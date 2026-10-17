package task

type TaskDTO struct {
	Title       string `json:"title"       validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
}

type Task struct {
	ID int `json:"id"`
	TaskDTO
	Done bool `json:"done"`
}
