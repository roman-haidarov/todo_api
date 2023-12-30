package handler

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/roman-haidarov/todo-app"
)

func (h *Handler) createList(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		var input todo.TodoList
		if err := c.BindJSON(&input); err != nil {
			newErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		id, err := h.services.TodoList.Create(userId, input)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
				"id": id,
		})
}

type getAllListsResponse struct {
		Data []todo.TodoList `json:"data"`
}

type getSearchListsResponse struct {
		Data []todo.TodoListSearch `json:"data"`
}

func (h *Handler) getAllLists(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		lists, err := h.services.TodoList.GetAll(userId)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, getAllListsResponse{
				Data: lists,
		})
}

func (h *Handler) getListById(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		list, err := h.services.TodoList.GetById(userId, id)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, list)
}

type SearchInputList struct {
		Search string `json:"search" binding:"required"`
}

func (h *Handler) getListsBySearch(c *gin.Context) {
		var search SearchInputList
		if err := c.BindJSON(&search); err != nil {
				newErrorResponse(c, http.StatusBadRequest, err.Error())
				return
		}

		lists, err := h.services.TodoList.GetListsBySearch(search.Search)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, getSearchListsResponse{
				Data: lists,
		})
}

func (h *Handler) updateList(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		var input todo.UpdateListInput
		if err := c.BindJSON(&input); err != nil {
				newErrorResponse(c, http.StatusBadRequest, err.Error())
				return
		}

		if err := h.services.UpdateList(userId, id, input); err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, statusResponse{
				Status: "ok",
		})
}

func (h *Handler) deleteList(c *gin.Context) {
		userId, err := getUserId(c)
		if err != nil {
				return
		}

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
				newErrorResponse(c, http.StatusBadRequest, "invalid id param")
				return
		}

		err = h.services.TodoList.DeleteList(userId, id)
		if err != nil {
				newErrorResponse(c, http.StatusInternalServerError, err.Error())
				return
		}

		c.JSON(http.StatusOK, statusResponse{
				Status: "ok",
		})
}
