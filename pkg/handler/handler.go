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
		ctx.Logger().Printf("key not found: %s", key)
		return ctx.NoContent(404)
	}
	if err != nil {
		ctx.Logger().Errorf("db error %s", err)
		return ctx.NoContent(500)
	}
	jsonResp := &map[string]string{}
	err = json.Unmarshal([]byte(resp), jsonResp)
	if err != nil {
		ctx.Logger().Errorf("error in unmarshall %s", jsonResp)
		return ctx.NoContent(500)
	}
	ctx.Logger().Printf("responsed {%s: %v}", key, resp)
	return ctx.JSON(200, jsonResp)
}

func (h *Handler) InsertValue(ctx echo.Context) error {
	var body map[string]interface{}
	err := ctx.Bind(&body)
	if err != nil {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	key, ok := body["key"]
	if !ok {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	if reflect.TypeOf(key).Kind() != reflect.String {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	value, ok := body["value"]
	if !ok {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	err = h.Repo.InsertValue(reflect.ValueOf(key).String(), string(valueJson))
	if errors.Is(err, repo.ErrKeyFound) {
		ctx.Logger().Printf("key already exists: %s", key)
		return ctx.NoContent(409)
	}
	if err != nil {
		ctx.Logger().Errorf("db error %s", err)
		return ctx.NoContent(500)
	}
	ctx.Logger().Printf("inserted {%s: %s}", key, string(valueJson))
	return ctx.JSON(200, value)
}

func (h *Handler) DeleteValue(ctx echo.Context) error {
	key := ctx.Param("key")
	err := h.Repo.DeleteValue(key)
	if errors.Is(err, repo.ErrNoContent) {
		ctx.Logger().Printf("key not found: %s", key)
		return ctx.NoContent(404)
	}
	if err != nil {
		ctx.Logger().Errorf("db error %s", err)
		return ctx.NoContent(500)
	}
	ctx.Logger().Printf("deleted by key: %s", key)
	return ctx.NoContent(200)
}

func (h *Handler) ChangeValue(ctx echo.Context) error {
	key := ctx.Param("key")
	var body map[string]interface{}
	err := ctx.Bind(&body)
	if err != nil {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	value, ok := body["value"]
	if !ok {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	valueJson, err := json.Marshal(value)
	if err != nil {
		ctx.Logger().Printf("incorrect body: %s", body)
		return ctx.NoContent(400)
	}
	err = h.Repo.ChangeValue(key, string(valueJson))
	if errors.Is(err, repo.ErrNoContent) {
		ctx.Logger().Printf("key not found: %s", key)
		return ctx.NoContent(404)
	}
	if err != nil {
		ctx.Logger().Errorf("db error %s", err)
		return ctx.NoContent(500)
	}
	ctx.Logger().Printf("updated {%s, %s}", key, string(valueJson))
	return ctx.NoContent(200)
}

func (h Handler) Hello(ctx echo.Context) error {
	hellomes := map[string]string{"message": "Hello!"}
	return ctx.JSON(200, hellomes)
}
