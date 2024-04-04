package main

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.POST("/pull", func(c echo.Context) error {
		var payload map[string]any

		if err := c.Bind(&payload); err != nil {
			log.Println("Error parsing JSON:", err)
			return c.String(http.StatusBadRequest, "Bad Request")
		}
		log.Println("Webhook received:", payload)

		if event, ok := payload["event"].(string); ok && event == "push" {
			log.Println("Image push event detected")

			// TODO:

			// cmd := exec.Command("docker", "pull", "your_image_name")
			// _, err := cmd.Output()
			// if err != nil {
			// 	log.Println("Error pulling image:", err)
			// 	return c.String(http.StatusInternalServerError, "Internal Server Error")
			// }

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
