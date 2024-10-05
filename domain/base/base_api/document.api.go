package base_api

import (
	"net/http"
	"os"
	"path/filepath"

	"GoRestify/domain/base"
	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/service"
	"GoRestify/internal/core"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_terms"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// DocumentAPI for injecting document service
type DocumentAPI struct {
	Service service.BaseDocumentServ
	Engine  *core.Engine
}

// ProvideDocumentAPI for document is used in wire
func ProvideDocumentAPI(c service.BaseDocumentServ) DocumentAPI {
	return DocumentAPI{Service: c, Engine: c.Engine}
}

// DownloadDocs finds the document via its doc_id and doc_type
func (a *DocumentAPI) DownloadDocs(c *gin.Context) {
	var err error
	var DirectoryPath string
	docName := c.Param("docName")
	docType := c.Param("docType")

	if DirectoryPath, _, err = a.Service.DocumentPathFolder(pkg_types.Enum(docType)); err != nil {
		return
	}

	fileFullPath := filepath.Join(DirectoryPath, docName)
	if _, err = os.Stat(fileFullPath); os.IsNotExist(err) {
		pkg_err.Log(err, "can't download document", fileFullPath)
	}

	c.FileAttachment(fileFullPath, docName)
}

// List of documents
func (a *DocumentAPI) List(c *gin.Context) {
	resp, params := response.NewParam(c, base_model.DocumentTable)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = a.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(base.ListDocument)
	resp.Status(http.StatusOK).
		Message(pkg_terms.ListOfV, base_term.Documents).
		JSON(data)
}
