package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/roman-haidarov/todo-app"
)

func (h *Handler) createItem(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		listId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
				return
		}

		var input todo.TodoItem
		if err := c.BindJSON(&input); err != nil {
				newErrorResponse(c, http.StatusBadRequest, err.Error())
				return
		}

		id, err := h.services.TodoItem.Create(userId, listId, input)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
				"id": id,
		})
}

type getAllItemsResponse struct {
		Data []todo.TodoItem `json:"data"`
}

type getSearchItemsResponse struct {
		count int
		Data []todo.TodoItemSearch `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		listId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		items, err := h.services.TodoItem.GetAll(userId, listId)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, getAllItemsResponse{
				Data: items,
		})
}

type SearchInputItem struct {
		Search string `json:"search" binding:"required"`
}

func (h *Handler) getItemsBySearch(c *gin.Context) {
		var search SearchInputItem
		if err := c.BindJSON(&search); err != nil {
				newErrorResponse(c, http.StatusBadRequest, err.Error())
				return
		}

		items, err := h.services.TodoList.GetItemsBySearch(search.Search)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, getSearchItemsResponse{
				count: len(items),
				Data: items,
		})
}

func (h *Handler) getItemById(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		item, err := h.services.TodoItem.GetItemById(userId, itemId)

		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItem(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		var input todo.UpdateItemInput
		if err := c.BindJSON(&input); err != nil {
				newErrorResponse(c, http.StatusBadRequest, err.Error())
				return
		}

		if err := h.services.UpdateItem(userId, id, input); err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, statusResponse{
				Status: "ok",
		})
}

func (h *Handler) deleteItem(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		itemId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		err = h.services.TodoItem.DeleteItem(userId, itemId)

		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, statusResponse{
				Status: "ok",
		})
}
