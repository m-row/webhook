package main

import (
	"encoding/json"
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
		log.Println("Webhook received:", payload)
		b, err := json.Marshal(payload)
		if err != nil {
			log.Println("marshal err: ", err.Error())
		}
		log.Println(
			"----------------------------------------------------------------",
		)
		log.Println("Webhook received json:", string(b))
		log.Println(
			"----------------------------------------------------------------",
		)

		if payload.PushData.PushedAt != 0 {
			log.Println("Image push event detected")

			cmd := exec.Command("docker", "pull", payload.Repository.RepoName)
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
