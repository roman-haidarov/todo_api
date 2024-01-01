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

func (s *TodoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
		return s.repo.GetAll(userId, listId)
}

func (s *TodoItemService) GetItemById(userId, itemId int) (todo.TodoItem, error) {
		return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
		return s.repo.DeleteItem(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
		if err := input.Validate(); err != nil {
				return err
		}

		return s.repo.UpdateItem(userId, itemId, input)
}
