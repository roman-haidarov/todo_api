package todo

import "errors"

type TodoList struct {
		Id 					 int 		`json:"-" db:"id"`
		Title 			 string `json:"title" db:"title" binding:"required"`
		Descriptions string	`json:"descriptions" db:"descriptions"`
}

type TodoListSearch struct {
		Username string `json:"username"`
		Title 	 string `json:"title"`
}

type UserList struct {
		Id 		 int
		UserId int
		ListId int
}

type TodoItem struct {
		ID           int    `json:"-" db:"id"`
		Title        string `json:"title" db:"title" binding:"required"`
		Descriptions string `json:"descriptions" db:"descriptions"`
		Done         bool   `json:"done" db:"done"`
}

type TodoItemSearch struct {
    Username string    `db:"username"`
    Done     bool 		 `db:"done"`
}

type ListsItem struct {
		Id 		 int
		ListId int
		ItemId string
}

type UpdateListInput struct {
		Title				 *string `json:"title"`
		Descriptions *string `json:"descriptions"`
}

func (i UpdateListInput) Validate() error {
		if i.Title == nil && i.Descriptions == nil {
				return errors.New("update structure has no values")
		}
		return nil
}
