package app

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"crypto/sha1"
	"encoding/hex"
	"lab1/internal/app/role"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

func (app *Application) uploadImage(c *gin.Context, image *multipart.FileHeader, UUID string) (*string, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	extension := filepath.Ext(image.Filename)
	if extension != ".jpg" && extension != ".jpeg" {
		return nil, fmt.Errorf("разрешены только jpeg изображения")
	}
	imageName := UUID + extension

	_, err = app.minioClient.PutObject(c, app.config.Minio.BucketName, imageName, src, image.Size, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})
	if err != nil {
		return nil, err
	}
	imageURL := fmt.Sprintf("%s/%s/%s", app.config.Minio.Endpoint, app.config.Minio.BucketName, imageName)
	return &imageURL, nil
}

func (app *Application) deleteImage(c *gin.Context, UUID string) error {
	imageName := UUID + ".jpg"
	fmt.Println(imageName)
	err := app.minioClient.RemoveObject(c, app.config.Minio.BucketName, imageName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func getUserId(c *gin.Context) string {
	userId, _ := c.Get("userId")
	return userId.(string)
}

func getUserRole(c *gin.Context) role.Role {
	userRole, _ := c.Get("userRole")
	return userRole.(role.Role)
}

func generateHashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func fundingRequest(mission_id string) error {
	url := "http://127.0.0.1:8082/"
	payload := fmt.Sprintf(`{"mission_id": "%s"}`, mission_id)

	resp, err := http.Post(url, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return fmt.Errorf(`funding failed with status: {%s}`, resp.Status)
	}
	return nil
}
