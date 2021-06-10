package handler

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/FluFFka/tarantool-kv-storage/pkg/repo"
	"github.com/labstack/echo"
)

type Handler struct {
	Repo *repo.Repository
}

func (h *Handler) GetByKey(ctx echo.Context) error {
	key := ctx.Param("key")
	resp, err := h.Repo.GetByKey(key)
	if errors.Is(err, repo.ErrNoContent) {
		return ctx.NoContent(404)
	}
	jsonResp := &map[string]string{}
	err = json.Unmarshal([]byte(resp), jsonResp)
	if err != nil {
		return ctx.NoContent(500)
	}
	return ctx.JSON(200, jsonResp)
}

func (h *Handler) InsertValue(ctx echo.Context) error {
	var body map[string]interface{}
	err := ctx.Bind(&body)
	if err != nil {
		return err
	}
	key, ok := body["key"]
	if !ok {
		return ctx.NoContent(400)
	}
	if reflect.TypeOf(key).Kind() != reflect.String {
		return ctx.NoContent(400)
	}
	value, ok := body["value"]
	if !ok {
		return ctx.NoContent(400)
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		return ctx.NoContent(400)
	}
	err = h.Repo.InsertValue(reflect.ValueOf(key).String(), string(valueJson))
	if errors.Is(err, repo.ErrKeyFound) {
		return ctx.NoContent(409)
	}
	if err != nil {
		return ctx.NoContent(500)
	}
	return ctx.JSON(200, value)
}

func (h *Handler) DeleteValue(ctx echo.Context) error {
	key := ctx.Param("key")
	err := h.Repo.DeleteValue(key)
	if errors.Is(err, repo.ErrNoContent) {
		return ctx.NoContent(404)
	}
	if err != nil {
		return ctx.NoContent(500)
	}
	return ctx.NoContent(200)
}

func (h *Handler) ChangeValue(ctx echo.Context) error {
	key := ctx.Param("key")
	var body map[string]interface{}
	err := ctx.Bind(&body)
	if err != nil {
		return ctx.NoContent(400)
	}
	value, ok := body["value"]
	if !ok {
		return ctx.NoContent(400)
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		return ctx.NoContent(400)
	}
	err = h.Repo.ChangeValue(key, string(valueJson))
	if errors.Is(err, repo.ErrNoContent) {
		return ctx.NoContent(404)
	}
	if err != nil {
		return ctx.NoContent(500)
	}
	return ctx.NoContent(200)
}

func (h Handler) Hello(ctx echo.Context) error {
	hellomes := map[string]string{"message": "Hello!"}
	return ctx.JSON(200, hellomes)
}
