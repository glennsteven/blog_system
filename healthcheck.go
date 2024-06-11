package blog_system

import (
	"blog-system/internal/helper"
	"blog-system/internal/resources"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	var response resources.Response
	response.Code = http.StatusOK
	response.Message = "ok"
	helper.ResponseJSON(w, http.StatusOK, response)
	return
}
