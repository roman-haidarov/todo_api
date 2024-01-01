package repository

import (
	"fmt"
	"strings"
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

func (r *TodoItemPostgres) GetItemById(userId, itemId int) (todo.TodoItem, error) {
		var item todo.TodoItem
		query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.descriptions FROM %s ti
													INNER JOIN %s ui on ti.id = ui.item_id WHERE ui.user_id = $1 AND ui.item_id = $2`,
				todoItemsTable, userItemsTable)
		err := r.db.Get(&item, query, userId, itemId)

		return item, err
}

func (r *TodoItemPostgres) DeleteItem(userId, itemId int) error {
		query := fmt.Sprintf(`DELETE FROM %s ti USING %s li, %s ul, %s ui
													WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1
													AND ui.item_id = ti.id AND ui.user_id = $1 AND ti.id = $2`,
				todoItemsTable, listsItemsTable, userListsTable, userItemsTable)

		_, err := r.db.Exec(query, userId, itemId)

		return err
}

func (r *TodoItemPostgres) UpdateItem(userId, itemId int, input todo.UpdateItemInput) error {
		setValues := make([]string, 0)
		args 			:= make([]interface{}, 0)
		argId 		:= 1

		if input.Title != nil {
				setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
				args 			= append(args, *input.Title)
				argId++
		}

		if input.Descriptions != nil {
				setValues = append(setValues, fmt.Sprintf("descriptions=$%d", argId))
				args 			= append(args, *input.Descriptions)
				argId++
		}


		if input.Done != nil {
				setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
				args 			= append(args, *input.Done)
				argId++
		}

		setQuery := strings.Join(setValues, ", ")
		query    := fmt.Sprintf(`UPDATE %s ti SET %s FROM %s li, %s ul
														 WHERE ti.id = li.item_id AND li.list_id=ul.list_id AND ul.user_id=$%d AND ti.id=$%d`,
														todoItemsTable, setQuery, listsItemsTable, userListsTable, argId, argId+1)
		args = append(args, userId, itemId)
		_, err := r.db.Exec(query, args...)
		return err
}
