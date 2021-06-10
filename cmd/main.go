package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/tarantool/go-tarantool"

	"github.com/FluFFka/tarantool-kv-storage/pkg/handler"
	"github.com/FluFFka/tarantool-kv-storage/pkg/repo"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect("localhost:3301", opts) //host.docker.internal
	if err != nil {
		fmt.Println("tarantool connection failed ", err)
		return
	}
	repository := &repo.Repository{Conn: conn}
	/*resp, err := conn.Insert("storage", []interface{}{"message", `{"message":"Hello"}`})
	if err != nil {
		fmt.Println("db error", err)
		fmt.Println(resp.Code)
		return
	}
	resp, err := conn.Select("storage", "primary", 0, 100, tarantool.IterAll, []interface{}{"message"})
	if err != nil {
		fmt.Println("db error", err)
		fmt.Println(resp.Code)
		return
	}
	fmt.Println(resp.Data)
	for _, item := range resp.Data {	//[[message {"message":"Hello"}]]
		fmt.Println(item)
	}
	*/
	h := &handler.Handler{
		Repo: repository,
	}

	e.GET("/", h.Hello)
	e.POST("/kv", h.InsertValue)
	e.GET("/kv/:key", h.GetByKey)

	e.Logger.Fatal(e.Start(":80"))
}
