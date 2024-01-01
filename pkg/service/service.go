package service

import (
	"github.com/roman-haidarov/todo-app"
	"github.com/roman-haidarov/todo-app/pkg/repository"
)

type Authorization interface {
		CreateUser(user todo.User) (int, error)
		GenerateToken(username, password string) (string, error)
		ParseToken(token string) (int, error)
}

type TodoList interface {
		Create(userId int, list todo.TodoList) (int, error)
		GetAll(userId int) ([]todo.TodoList, error)
		GetListsBySearch(search string) ([]todo.TodoListSearch, error)
		GetItemsBySearch(search string) ([]todo.TodoItemSearch, error)
		GetById(userId, listId int) (todo.TodoList, error)
		DeleteList(userId, listId int) error
		UpdateList(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface {
		Create(userId, listId int, input todo.TodoItem) (int, error)
		GetAll(userId, listId int) ([]todo.TodoItem, error)
		GetItemById(userId, itemId int) (todo.TodoItem, error)
		DeleteItem(userId, itemId int) error
		UpdateItem(userId, itemId int, input todo.UpdateItemInput) error
}

type Service struct {
		Authorization
		TodoList
		TodoItem
}

func NewService(repos *repository.Repository) *Service {
		return &Service{
				Authorization: NewAuthService(repos.Authorization),
				TodoList: NewTodoListService(repos.TodoList),
				TodoItem: NewTodoItemService(repos.TodoItem, repos.TodoList),
		}
}
