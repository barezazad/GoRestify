package base_api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	var directoryPath string
	docName := filepath.Base(c.Param("docName"))
	docType := c.Param("docType")

	if docName == "" || docName == "." || docName == ".." {
		err = pkg_err.New("invalid document name", "E1159001").
			Custom(pkg_err.BadRequestErr).Message(pkg_err.BadRequest).Build()
		response.New(c).Error(err).JSON()
		return
	}

	if directoryPath, _, err = a.Service.DocumentPathFolder(pkg_types.Enum(docType)); err != nil {
		response.New(c).Error(err).JSON()
		return
	}

	fileFullPath := filepath.Join(directoryPath, docName)
	absRoot, err := filepath.Abs(directoryPath)
	if err != nil {
		response.New(c).Error(err).JSON()
		return
	}
	absFile, err := filepath.Abs(fileFullPath)
	if err != nil {
		response.New(c).Error(err).JSON()
		return
	}
	rootPrefix := absRoot + string(os.PathSeparator)
	if absFile != absRoot && !strings.HasPrefix(absFile, rootPrefix) {
		err = pkg_err.New("invalid document path", "E1159002").
			Custom(pkg_err.BadRequestErr).Message(pkg_err.BadRequest).Build()
		response.New(c).Error(err).JSON()
		return
	}

	if _, err = os.Stat(absFile); os.IsNotExist(err) {
		err = pkg_err.New("document not found", "E1159003").
			Custom(pkg_err.NotFoundErr).Message(pkg_err.RecordNotFound).Build()
		response.New(c).Error(err).JSON()
		return
	}

	c.FileAttachment(absFile, docName)
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
