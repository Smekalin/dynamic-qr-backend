package v1

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"rsc.io/qr"
)

type GeneratorCtrl struct {
}

func NewGeneratorCtrl() *GeneratorCtrl { return &GeneratorCtrl{} }

var a = make(map[string][]byte)

func (c *GeneratorCtrl) Register(e *echo.Group) {
	g := e.Group("/generator")
	g.GET("", c.GetList)
	g.GET("/create", c.GenerateQR)
}

func (c *GeneratorCtrl) GenerateQR(e echo.Context) error {
	url := e.Request().URL.Query().Get("url")

	encode, err := qr.Encode(url, 2)
	if err != nil {
		return err
	}
	png := encode.PNG()

	// TODO: RECORD to DB
	a[url] = png

	e.Response().Header().Set("Content-Type", "image/png")
	_, err = e.Response().Write(png)
	if err != nil {
		return err
	}
	return nil
}

func (c *GeneratorCtrl) GetList(e echo.Context) error {
	fmt.Println(e.Request().Context())

	// TODO: READ from DB
	return e.JSON(http.StatusOK, a)

}
