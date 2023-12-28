package todo

type TodoList struct {
		Id int 							`json:"-" db:"id"`
		Title string 				`json:"title" db:"title" binding:"required"`
		Descriptions string	`json:"descriptions" db:"descriptions"`
}

type UserList struct {
		Id int
		UserId int
		ListId int
}

type TodoItem struct {
		Id int 							`json:"-"`
		Title string 				`json:"title"`
		Descriptions string	`json:"descriptions"`
		Done bool						`json:"done"`
}

type ListsItem struct {
		Id int
		ListId int
		ItemId string
}
