package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

func main() {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")

	e := echo.New()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.43"))
	if err != nil {
		e.Logger.Fatalf("Error creating Docker client: %v", err)
	}
	authConfig := registry.AuthConfig{
		Username: username,
		Password: password,
	}
	res, err := cli.RegistryLogin(context.Background(), authConfig)
	if err != nil {
		e.Logger.Fatalf("Error RegistryLogin: %v", err)
	}
	e.Logger.Info("RegistryLogin res:", res.Status)

	e.POST("/pull", func(c echo.Context) error {
		var payload Payload

		if err := c.Bind(&payload); err != nil {
			e.Logger.Info("Error parsing JSON:", err)
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		b, err := json.Marshal(payload)
		if err != nil {
			e.Logger.Info("marshal err: ", err.Error())
		}
		e.Logger.Info("Webhook received json:", string(b))
		e.Logger.Info(
			"----------------------------------------------------------------",
		)
		e.Logger.Info("Webhook received PushedAt:", payload.PushData.PushedAt)
		e.Logger.Info(
			"----------------------------------------------------------------",
		)

		if payload.PushData.PushedAt != 0 {
			e.Logger.Info("Image push event detected")

			if _, err = cli.ImagePull(
				context.Background(),
				payload.Repository.RepoName,
				image.PullOptions{},
			); err != nil {
				e.Logger.Fatalf("Error pulling image: %v", err)
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
