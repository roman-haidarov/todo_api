package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/roman-haidarov/todo-app"
)

type TodoListPostgres struct {
		db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
		return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todo.TodoList) (int, error) {
		tx, err := r.db.Begin()
		if err != nil {
			return 0, err
		}

		var id int
		createListQuery := fmt.Sprintf("INSERT INTO %s (title, descriptions) VALUES ($1, $2) RETURNING id", todoListsTable)
		row := tx.QueryRow(createListQuery, list.Title, list.Descriptions)
		if err := row.Scan(&id); err != nil {
			tx.Rollback()
			return 0, err
		}

		createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", userListsTable)
		_, err = tx.Exec(createUsersListQuery, userId, id)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
		var lists []todo.TodoList
		query := fmt.Sprintf("SELECT tl.id, tl.title, tl.descriptions FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1",
				todoListsTable, userListsTable)
		err := r.db.Select(&lists, query, userId)

		return lists, err
}

func (r *TodoListPostgres) GetListsBySearch(serach string) ([]todo.TodoList, error) {
		var lists []todo.TodoList
		query := fmt.Sprintf("SELECT id, title, descriptions FROM %s WHERE tsv @@ to_tsquery($1)", todoListsTable)
		err := r.db.Select(&lists, query, serach)
		
		return lists, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
		var list todo.TodoList
		query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.descriptions FROM %s tl
													INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
				todoListsTable, userListsTable)
		err := r.db.Get(&list, query, userId, listId)

		return list, err
}
