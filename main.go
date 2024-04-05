package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
)

func main() {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")

	e := echo.New()
	ctx := context.Background()
	authConfig := registry.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		log.Fatalf("error when encoding authConfig. err: %v", err)
	}

	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	log.Printf("authStr created: %s", authStr)
	cli, err := client.NewClientWithOpts(
		client.WithVersion("1.43"),
	)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v", err)
	}
	log.Printf("client created: %s", cli.ClientVersion())

	ping, err := cli.Ping(ctx)
	if err != nil {
		log.Fatalf("Error pinging Docker client: %v", err)
	}
	log.Printf("client pinged: %s", ping.APIVersion)

	res, err := cli.RegistryLogin(ctx, authConfig)
	if err != nil {
		log.Fatalf("Error RegistryLogin: %v", err)
	}
	log.Printf("RegistryLogin res: %s", res.Status)
	log.Printf("RegistryLogin res.IdentityToken: %s", res.IdentityToken)

	e.POST("/pull", func(c echo.Context) error {
		var payload Payload

		if err := c.Bind(&payload); err != nil {
			log.Print("Error parsing JSON:", err)
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		b, err := json.Marshal(payload)
		if err != nil {
			log.Print("marshal err: ", err.Error())
		}
		log.Print("Webhook received json:", string(b))
		log.Print(
			"----------------------------------------------------------------",
		)
		log.Print("Webhook received PushedAt:", payload.PushData.PushedAt)
		log.Print(
			"----------------------------------------------------------------",
		)

		if payload.PushData.PushedAt != 0 {
			log.Print("Image push event detected")

			out, err := cli.ImagePull(
				ctx,
				payload.Repository.RepoName,
				image.PullOptions{
					RegistryAuth: authStr,
				},
			)
			if err != nil {
				log.Fatalf("Error pulling image: %v", err)
			}
			defer out.Close()
			body, err := io.ReadAll(out)
			if err != nil {
				log.Print("body ReadAll error: ", err.Error())
			}
			log.Println("out body:", string(body))

			if err := cli.ContainerRemove(
				ctx,
				payload.Repository.Name,
				container.RemoveOptions{
					RemoveVolumes: false,
					RemoveLinks:   false,
					Force:         true,
				},
			); err != nil {
				log.Print("container remove error: ", err.Error())
			}

			creatRes, err := cli.ContainerCreate(
				ctx,
				&container.Config{
					Image: payload.Repository.Name,
				},
				nil,
				&network.NetworkingConfig{
					EndpointsConfig: map[string]*network.EndpointSettings{
						"sadeem": {
							NetworkID: "sadeem",
						},
					},
				},
				nil,
				"",
			)
			if err != nil {
				log.Print("container create error: ", err.Error())
			}
			log.Println("container created id: ", creatRes.ID)

			return c.String(
				http.StatusOK,
				"Webhook received and processed successfully",
			)
		}

		return c.String(http.StatusBadRequest, "Invalid Webhook Event")
	})

	e.Logger.Fatal(e.Start(":8000"))
}
