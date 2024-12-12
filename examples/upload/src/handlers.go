package src

import (
	"fmt"
	"github.com/aiechoic/admin/core/gins"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (hs *Handlers) UploadFile(fileServeUrl, dir string) gins.Handler {
	type params struct {
		File *multipart.FileHeader `form:"file" binding:"required" description:"The file"`
	}
	type response struct {
		Url string `json:"url" description:"The url of the file"`
	}
	return gins.Handler{
		Request: gins.Request{
			Form: params{},
		},
		Response: gins.Response{
			Json: response{},
		},
		Handle: func(c *gin.Context) {
			var p params
			if err := c.ShouldBind(&p); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			filename := p.File.Filename
			err := c.SaveUploadedFile(p.File, filepath.Join(dir, filename))
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, response{Url: fmt.Sprintf("%s/%s", fileServeUrl, filename)})
		},
	}
}

func (hs *Handlers) ServeFile(pathParamName, dir string) gins.Handler {
	return gins.Handler{
		Response: gins.Response{
			Contents: gins.ContentsOctetStream,
		},
		Handle: func(c *gin.Context) {
			file := c.Param(pathParamName)
			c.File(filepath.Join(dir, file))
		},
	}
}
