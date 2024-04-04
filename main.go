package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"

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

			dockerpull := fmt.Sprintf(
				"docker pull %s",
				payload.Repository.RepoName,
			)
			//nolint: gosec
			cmd := exec.Command("/bin/bash", "-c", dockerpull)
			_, err := cmd.Output()
			if err != nil {
				log.Println("Error pulling image:", err)
				return c.String(
					http.StatusInternalServerError,
					"Internal Server Error",
				)
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
