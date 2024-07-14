package main

import (
	"context"
	"fmt"

	"go-service/internal/app"
)

func main() {
	var cfg app.Config
	cfg.Mongo.Uri = "mongodb+srv://dbUser:Demoaccount1@projectdemo.g0lah.mongodb.net"
	cfg.Mongo.Database = "masterdata"

	ctx := context.Background()
	fmt.Println("Import file")
	app, err := app.NewApp(ctx, cfg)
	if err != nil {
		fmt.Println("Error when initialize: ", err.Error())
		panic(err)
	}

	total, success, err := app.Import(ctx)
	if err != nil {
		fmt.Println("Error when import: ", err.Error())
		panic(err)
	}
	fmt.Println(fmt.Sprintf("total: %d, success: %d", total, success))
	fmt.Println(ctx, "Imported file")
}
