package v1

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"rsc.io/qr"
)

type GeneratorCtrl struct {
}

var stat = Stat{Usage: 0, Links: make([]string, 0, 10000)}

func NewGeneratorCtrl() *GeneratorCtrl { return &GeneratorCtrl{} }

func (c *GeneratorCtrl) Register(e *echo.Group) {
	g := e.Group("/generator")
	g.GET("/create", c.generateQR)
	g.GET("/stat", c.getStat)
}

func (c *GeneratorCtrl) generateQR(e echo.Context) error {
	url := e.Request().URL.Query().Get("url")

	encode, err := qr.Encode(url, 2)
	if err != nil {
		return err
	}
	png := encode.PNG()

	updateStat(url)

	e.Response().Header().Set("Content-Type", "image/png")
	_, err = e.Response().Write(png)
	if err != nil {
		return err
	}
	return nil
}

func (c *GeneratorCtrl) getStat(e echo.Context) error {
	parsedJson, err := json.Marshal(stat)
	if err != nil {
		return e.HTML(http.StatusOK, "SOME ERROR")
	}

	return e.JSONBlob(http.StatusOK, parsedJson)
}

func updateStat(url string) {
	stat.Usage += 1

	if len(stat.Links) == 5 {
		stat.Links = stat.Links[:0]
	}

	stat.Links = append(stat.Links, url)

}

type Stat struct {
	Usage int64    `json:"usage"`
	Links []string `json:"links"`
}
