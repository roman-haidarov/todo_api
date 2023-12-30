package service

import (
	"github.com/roman-haidarov/todo-app"
	"github.com/roman-haidarov/todo-app/pkg/repository"
)

type TodoListService struct {
		repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
		return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(userId int, list todo.TodoList) (int, error) {
		return s.repo.Create(userId, list)
}

func (s *TodoListService) GetAll(userId int) ([]todo.TodoList, error) {
		return s.repo.GetAll(userId)
}

func (s *TodoListService) GetListsBySearch(search string) ([]todo.TodoListSearch, error) {
		return s.repo.GetListsBySearch(search)
}

func (s *TodoListService) GetItemsBySearch(search string) ([]todo.TodoItemSearch, error) {
		return s.repo.GetItemsBySearch(search)
}

func (s *TodoListService) GetById(userId, listId int) (todo.TodoList, error) {
		return s.repo.GetById(userId, listId)
}

func (s *TodoListService) DeleteList(userId, listId int) error {
		return s.repo.DeleteList(userId, listId)
}

func (s *TodoListService) UpdateList(userId, listId int, input todo.UpdateListInput) error {
		if err := input.Validate(); err != nil {
				return err
		}

		return s.repo.UpdateList(userId, listId, input)
}
