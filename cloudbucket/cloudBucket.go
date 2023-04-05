package cloudbucket

import (
	"io"
	"net/http"
	"net/url"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

// HandleFileUploadToBucket uploads file to bucket
func HandleFileUploadToBucket(c *gin.Context) {
	bucket := "cybergame-dev" //your bucket name

	var err error

	ctx := appengine.NewContext(c.Request)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsFile("google_service_account.json"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	f, uploadedFile, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)

	if _, err := io.Copy(sw, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err := sw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	u, err := url.Parse("/" + bucket + "/" + "image/setting/" + sw.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}
