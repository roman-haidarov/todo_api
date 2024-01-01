package repository

import (
	"fmt"
	"strings"
	"github.com/jmoiron/sqlx"
	"github.com/roman-haidarov/todo-app"
	"github.com/sirupsen/logrus"
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
		query := fmt.Sprintf("SELECT tl.id, tl.title, tl.descriptions FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 ORDER BY tl.id DESC",
				todoListsTable, userListsTable)
		err := r.db.Select(&lists, query, userId)

		return lists, err
}

func (r *TodoListPostgres) GetListsBySearch(search string) ([]todo.TodoListSearch, error) {
		var lists []todo.TodoListSearch

		query := fmt.Sprintf(`SELECT u.username, tl.title FROM %s tl
													INNER JOIN %s ul on tl.id = ul.list_id
													INNER JOIN %s u on ul.user_id = u.id
													WHERE u.tsv @@ to_tsquery($1) OR tl.tsv @@ to_tsquery($1)
													ORDER BY tl.id DESC`,
		todoListsTable, userListsTable, usersTable)
		err := r.db.Select(&lists, query, search)

		return lists, err
}

func (r *TodoListPostgres) GetItemsBySearch(search string) ([]todo.TodoItemSearch, error) {
		var items []todo.TodoItemSearch

    query := fmt.Sprintf(`SELECT u.username, ti.done FROM %s ti
													INNER JOIN %s ui ON ti.id = ui.item_id
													INNER JOIN %s u ON ui.user_id = u.id
													WHERE u.tsv @@ to_tsquery($1) OR ti.tsv @@ to_tsquery($1)
													ORDER BY u.username DESC`,
													todoItemsTable, userItemsTable, usersTable)

		err := r.db.Select(&items, query, search)

		return items, err
}

func (r *TodoListPostgres) GetById(userId, listId int) (todo.TodoList, error) {
		var list todo.TodoList
		query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.descriptions FROM %s tl
													INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2`,
				todoListsTable, userListsTable)
		err := r.db.Get(&list, query, userId, listId)

		return list, err
}

func (r *TodoListPostgres) DeleteList(userId, listId int) error {
		query := fmt.Sprintf(`DELETE FROM %s tl USING %s ul
													WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2`,
				todoListsTable, userListsTable)
		_, err := r.db.Exec(query, userId, listId)

		return err
}

func (r *TodoListPostgres) UpdateList(userId, listId int, input todo.UpdateListInput) error {
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

		setQuery := strings.Join(setValues, ", ")
		query    := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
														todoListsTable, setQuery, userListsTable, argId, argId+1)
		args = append(args, listId, userId)

		logrus.Debugf("udateQuery: %s", query)
		logrus.Debugf("args: %s", args)

		_, err := r.db.Exec(query, args...)
		return err
}
