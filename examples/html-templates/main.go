package main

import (
	"log"

	"github.com/DevanshuTripathi/vodka"
)

func main() {
	app := vodka.DefaultRouter()

	err := app.LoadHTMLGlob("templates/*.html")
	if err != nil {
		log.Fatal("Failed to load templates:", err)
	}

	app.GET("/hello", func(c *vodka.Context) {
		c.HTML(200, "index.html", vodka.M{
			"Title":   "Vodka Templates",
			"Name":    "Developer",
			"Message": "Welcome to Vodka framework!",
		})
	})

	app.GET("/user/:name", func(c *vodka.Context) {
		name := c.Param("name")
		c.HTML(200, "user.html", vodka.M{
			"Title": "User Profile",
			"Name":  name,
		})
	})

	log.Println("Server running on http://localhost:8080")
	log.Println("Try: http://localhost:8080/hello")

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}