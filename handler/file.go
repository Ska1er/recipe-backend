package handler

import (
	"fmt"
	"net/http"
	"path/filepath"

	"recipe/internal"

	"github.com/gin-gonic/gin"
)

type FileHandler struct {
}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (handler *FileHandler) UploadHanlder(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation exception",
			"fields": gin.H{
				"file": "This field should be file",
			},
		})
		return
	}

	filename := internal.GenerateUniqueFilename(filepath.Base(file.Filename))
	fmt.Printf("Generated filename: %s\n", filename)
	path := fmt.Sprintf("%s/%s", internal.StaticFolder, filename)
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to save uploaded file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
	})
}
