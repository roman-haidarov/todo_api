package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/roman-haidarov/todo-app"
)

type Authorization interface {
		CreateUser(user todo.User) (int, error)
		GetUser(username, password string) (todo.User, error)
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
		Create(userId, listId int, item todo.TodoItem) (int, error)
		GetAll(userId, listId int) ([]todo.TodoItem, error)
		GetItemById(userId, itemId int) (todo.TodoItem, error)
		DeleteItem(userId, itemId int) error
		UpdateItem(userId, itemId int, input todo.UpdateItemInput) error
}

type Repository struct {
		Authorization
		TodoList
		TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
		return &Repository{
				Authorization: NewAuthPostgres(db),
				TodoList: NewTodoListPostgres(db),
				TodoItem: NewTodoItemPostgres(db),
		}
}
