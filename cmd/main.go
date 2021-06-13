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
	e.Use(middleware.Logger())

	opts := tarantool.Opts{User: "guest"}
	conn, err := tarantool.Connect("host.docker.internal:3301", opts)
	// host.docker.internal or localhost or
	// (ip addr show | grep "\binet\b.*\bdocker0\b" | awk '{print $2}' | cut -d '/' -f 1)
	if err != nil {
		fmt.Println("tarantool connection failed ", err)
		return
	}
	repository := &repo.Repository{Conn: conn}
	h := &handler.Handler{Repo: repository}

	e.POST("/kv", h.InsertValue)
	e.GET("/kv/:key", h.GetByKey)
	e.DELETE("/kv/:key", h.DeleteValue)
	e.PUT("/kv/:key", h.ChangeValue)

	e.Logger.Fatal(e.Start(":80"))
}
