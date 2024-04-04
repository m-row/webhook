package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/pull", func(c echo.Context) error {
		var payload Payload

		if err := c.Bind(&payload); err != nil {
			log.Println("Error parsing JSON:", err)
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("marshal err: ", err.Error())
		}
		log.Println("Webhook received json:", string(b))
		log.Println(
			"----------------------------------------------------------------",
		)
		log.Println("Webhook received PushedAt:", payload.PushData.PushedAt)
		log.Println(
			"----------------------------------------------------------------",
		)

		if payload.PushData.PushedAt != 0 {
			log.Println("Image push event detected")

			cli, err := client.NewClientWithOpts(client.WithVersion("1.43"))
			if err != nil {
				log.Fatalf("Error creating Docker client: %v", err)
			}
			if _, err = cli.ImagePull(
				context.Background(),
				payload.Repository.RepoName,
				image.PullOptions{},
			); err != nil {
				log.Fatalf("Error pulling image: %v", err)
			}
			// TODO: restarting containers, etc.

			return c.String(
				http.StatusOK,
				"Webhook received and processed successfully",
			)
		}

		return c.String(http.StatusBadRequest, "Invalid Webhook Event")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
