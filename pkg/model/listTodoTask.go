package model

type ListTodoTask struct {
	Tasks []Task `json:"tasks"`
}

func NewListTodoTask(task []Task) *ListTodoTask {
	return &ListTodoTask{Tasks: task}
}
