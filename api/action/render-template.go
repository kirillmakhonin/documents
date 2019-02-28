package action

import (
	"encoding/json"
	"net/http"

	"github.com/opsway/documents/cmd/document"
	"github.com/opsway/documents/cmd/template"
)

// RenderTemplateRequest refers request of rendering pdf
type RenderTemplateRequest struct {
	TemplateName    string            `json:"templateName"`
	Data            template.Context  `json:"data"`
	DocumentOptions document.Document `json:"documentOptions"`
}

// RenderTemplate render pdf by template from request
func RenderTemplate(w http.ResponseWriter, r *http.Request) {
	request := RenderTemplateRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pdf, err := document.NewPdf()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pdf.SetOptions(request.DocumentOptions)
	err = pdf.RenderByTemplate(w, request.TemplateName, request.Data)

	if err != nil {
		panic(err)
	}
}
