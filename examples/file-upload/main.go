package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/DevanshuTripathi/vodka"
)

func main() {
	app := vodka.DefaultRouter()

	// Health check endpoint
	app.GET("/health", func(c *vodka.Context) {
		c.JSON(200, vodka.M{
			"status":  "ok",
			"service": "file-upload-example",
		})
	})

	// POST /upload — accepts a multipart file upload, saves it to ./uploads,
	// and returns metadata about the saved file.
	app.POST("/upload", func(c *vodka.Context) {
		// Ensure the uploads directory exists
		uploadsDir := "./uploads"
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Printf("[upload] failed to create uploads directory: %v", err)
			c.JSON(500, vodka.M{
				"error": "could not prepare upload directory",
			})
			return
		}

		// Retrieve the file from the multipart form (field name: "file")
		fileHeader, err := c.FormFile("file")
		if err != nil {
			log.Printf("[upload] failed to parse multipart file: %v", err)
			c.JSON(400, vodka.M{
				"error": "missing or invalid file field — use form key \"file\"",
			})
			return
		}

		// Build a unique destination path to avoid name collisions
		timestamp := time.Now().UnixNano()
		safeName := filepath.Base(fileHeader.Filename)
		dstPath := filepath.Join(uploadsDir, fmt.Sprintf("%d_%s", timestamp, safeName))

		// Stream the uploaded file to disk using the built-in Vodka helper
		if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
			log.Printf("[upload] failed to save file %q: %v", safeName, err)
			c.JSON(500, vodka.M{
				"error": "failed to save uploaded file",
			})
			return
		}

		log.Printf("[upload] saved %q (%d bytes) → %s", safeName, fileHeader.Size, dstPath)

		// Respond with success and file metadata
		c.JSON(200, vodka.M{
			"message":  "file uploaded successfully",
			"filename": safeName,
			"size":     fileHeader.Size,
			"path":     dstPath,
		})
	})

	log.Println("File-upload example starting on :8080")
	if err := app.Run(":8080"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
