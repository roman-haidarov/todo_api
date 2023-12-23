package todo

type TodoList struct {
		Id int 							`json:"-"`
		Title string 				`json:"title"`
		Description string	`json:"body"`
}

type UserList struct {
		Id int
		UserId int
		ListId int
}

type TodoItem struct {
		Id int 							`json:"-"`
		Title string 				`json:"title"`
		Description string	`json:"body"`
		Done bool						`json:"done"`
}

type ListsItem struct {
		Id int
		ListId int
		ItemId string
}
