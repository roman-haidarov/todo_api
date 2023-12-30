package service

import (
	"github.com/roman-haidarov/todo-app"
	"github.com/roman-haidarov/todo-app/pkg/repository"
)

type TodoItemService struct {
		repo repository.TodoItem
		listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *TodoItemService {
		return &TodoItemService{repo: repo, listRepo: listRepo}
}

func (s *TodoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
		_, err := s.listRepo.GetById(userId, listId)

		if err != nil {
				return 0, err
		}

		return s.repo.Create(userId, listId, item)
}
