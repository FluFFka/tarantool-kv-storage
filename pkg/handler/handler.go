package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/FluFFka/tarantool-kv-storage/pkg/repo"
	"github.com/labstack/echo"
)

type Handler struct {
	Repo *repo.Repository
}

func (h Handler) GetByKey(ctx echo.Context) error {
	key := ctx.Param("key")
	resp, err := h.Repo.GetByKey(key)
	if errors.Is(err, repo.ErrNoContent) {
		return ctx.NoContent(404)
	}
	fmt.Println(resp)
	jsonResp := &map[string]string{}
	err = json.Unmarshal([]byte(resp), jsonResp)
	if err != nil {
		return ctx.NoContent(500)
	}
	return ctx.JSON(200, jsonResp)
}

func (h Handler) InsertValue(ctx echo.Context) error {
	var body map[string]interface{}
	err := ctx.Bind(&body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	key := body["key"]
	if reflect.TypeOf(key).Kind() != reflect.String {
		return ctx.NoContent(400)
	}
	value := body["value"]
	valueJson, err := json.Marshal(value)
	if err != nil {
		return ctx.NoContent(400)
	}
	err = h.Repo.InsertValue(reflect.ValueOf(key).String(), string(valueJson))
	if errors.Is(err, repo.ErrKeyFound) {
		return ctx.NoContent(409)
	}
	if err != nil {
		fmt.Println(err)
		return ctx.NoContent(500)
	}
	return ctx.JSON(200, value)
}

func (h Handler) Hello(ctx echo.Context) error {
	hellomes := map[string]string{"message": "Hello!"}
	return ctx.JSON(200, hellomes)
}
