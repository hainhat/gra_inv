package utils

import (
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

func RenderHTMLWithPartials(c *gin.Context, htmlPath string) {
	content, err := ioutil.ReadFile(htmlPath)
	if err != nil {
		c.String(500, "Error reading file")
		return
	}

	html := string(content)

	if strings.Contains(html, "{{header}}") {
		header, _ := ioutil.ReadFile("./frontend/partials/header.html")
		html = strings.Replace(html, "{{header}}", string(header), 1)
	}
	if strings.Contains(html, "{{footer}}") {
		footer, _ := ioutil.ReadFile("./frontend/partials/footer.html")
		html = strings.Replace(html, "{{footer}}", string(footer), 1)
	}
	if strings.Contains(html, "{{head}}") {
		head, _ := ioutil.ReadFile("./frontend/partials/head.html")
		html = strings.Replace(html, "{{head}}", string(head), 1)
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(200, html)
}
