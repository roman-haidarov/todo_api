package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/roman-haidarov/todo-app"
)

type TodoItemPostgres struct {
		db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
		return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(userId, listId int, item todo.TodoItem) (int, error) {
		tx, err := r.db.Begin()
		if err != nil {
				return 0, err
		}

		var itemId int
		createItemQuery := fmt.Sprintf("INSERT INTO %s (title, descriptions) VALUES ($1, $2) RETURNING id", todoItemsTable)
		row := tx.QueryRow(createItemQuery, item.Title, item.Descriptions)
		if err := row.Scan(&itemId); err != nil {
				tx.Rollback()
				return 0, err
		}

		createListsItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listsItemsTable)
		_, err = tx.Exec(createListsItemsQuery, listId, itemId)
		if err != nil {
				tx.Rollback()
				return 0, err
		}

		createUserItemsQuery := fmt.Sprintf("INSERT INTO %s (item_id, user_id) VALUES ($1, $2)", userItemsTable)
		_, err = tx.Exec(createUserItemsQuery, itemId, userId)
		if err != nil {
				tx.Rollback()
				return 0, err
		}

		if err := tx.Commit(); err != nil {
				return 0, err
		}

		return itemId, nil
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
		var items []todo.TodoItem
		query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.descriptions, ti.done FROM %s ti 
													INNER JOIN %s li on li.item_id = ti.id
													INNER JOIN %s ul on ul.list_id = li.list_id
													WHERE li.list_id = $1 AND ul.user_id = $2
													ORDER BY ti.id DESC`,
				todoItemsTable, listsItemsTable, userListsTable)
		err := r.db.Select(&items, query, listId, userId)

		return items, err
}
