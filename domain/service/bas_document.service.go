package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_repo"
	"GoRestify/domain/base/base_term"
	"GoRestify/domain/base/enum/document_type"
	"GoRestify/internal/core"

	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"
	"GoRestify/pkg/utils"
	"GoRestify/pkg/validator"
)

// BaseDocumentServ for injecting auth base_repo
type BaseDocumentServ struct {
	Repo   base_repo.DocumentRepo
	Engine *core.Engine
}

// ProvideBaseDocumentService for document is used in wire
func ProvideBaseDocumentService(documentRepo base_repo.DocumentRepo) BaseDocumentServ {
	return BaseDocumentServ{
		Repo:   documentRepo,
		Engine: documentRepo.Engine,
	}
}

// GetByRefrenceIdType get documents by refrence_id and type
func (s *BaseDocumentServ) GetByRefrenceIdType(tx tx.Tx, referenceID uint, docType pkg_types.Enum) (documents []base_model.Document, err error) {

	key := fmt.Sprintf("%v-type-%v-%v", base_term.Document, referenceID, docType)
	if ok := s.Engine.RedisCacheAPI.GetCache(tx, key, &documents); ok {
		return
	}

	if documents, err = s.Repo.GetByRefrenceIdType(tx, referenceID, docType); err != nil {
		pkg_err.Log(err, "E1144538", "can't fetch the documents", referenceID)
		return
	}

	err = s.Engine.RedisCacheAPI.Set(key, documents)

	return
}

// List of documents, it supports pagination and search and return count
func (s *BaseDocumentServ) List(params param.Param) (documents []base_model.Document,
	count int64, err error) {

	if documents, err = s.Repo.List(params); err != nil {
		pkg_log.CheckError(err, "error in documents list")
		return
	}

	if count, err = s.Repo.Count(params); err != nil {
		pkg_log.CheckError(err, "error in documents count")
	}

	return
}

// Create a document
func (s *BaseDocumentServ) Create(tx tx.Tx, document base_model.Document) (createdDocument base_model.Document, err error) {

	if err = validator.ValidateModel(document, base_term.Document, validator.Create); err != nil {
		pkg_err.TickValidate(err, "E1159448", pkg_err.ValidationFailed, document)
		return
	}

	if createdDocument, err = s.Repo.Create(tx, document); err != nil {
		pkg_err.Log(err, "E1175150", "document not saved")
		return
	}

	key := fmt.Sprintf("%v-type-%v-%v", base_term.Document, createdDocument.ReferenceID, createdDocument.Type)
	s.Engine.RedisCacheAPI.Delete(key)

	return
}

// UploadDocs to upload documents
func (s *BaseDocumentServ) UploadDocs(params param.Param, docs base_model.Document) (documents []base_model.Document, err error) {

	var DirectoryPath, AcceptedFiles string
	if DirectoryPath, AcceptedFiles, err = s.DocumentPathFolder(docs.Type); err != nil {
		pkg_err.Log(err, "E1184055", "document type not found")
		return
	}

	var docPayload base_model.Document
	var documentReader []base_model.DocumentReader

	switch {

	case docs.Attachments != nil:
		for _, v := range docs.Attachments {
			// read file like io file
			var out multipart.File
			if out, err = v.Open(); err != nil {
				err = pkg_err.Take(err, "E1053409").
					Message(pkg_err.BadRequest).
					Custom(pkg_err.BadRequestErr).Build()
				return
			}

			documentReader = append(documentReader, base_model.DocumentReader{File: out, Attachment: v})
		}

	case docs.AttachmentsBase64 != nil:
		for _, v := range docs.AttachmentsBase64 {
			// decode base64 file
			idx := strings.Index(v, ";base64,")
			out := base64.NewDecoder(base64.StdEncoding, strings.NewReader(v[idx+8:]))

			documentReader = append(documentReader, base_model.DocumentReader{File: out})
		}

	default:
		err = pkg_err.Take(err, "E1114486").
			Message(pkg_err.BadRequest).
			Custom(pkg_err.BadRequestErr).Build()
		return
	}

	for i, v := range documentReader {

		// read file in buffer
		buff := bytes.Buffer{}
		if _, err = buff.ReadFrom(v.File); err != nil {
			pkg_err.Log(err, "E1112830", "can't read file", buff)
		}

		filetype := http.DetectContentType(buff.Bytes())
		fileExt := strings.Split(filetype, "/")[1]

		FileName := utils.RandomString(20)
		FileName = fmt.Sprintf(`%v-%v-%v.%v`, docs.ReferenceID, i+1, FileName, fileExt)
		FilePath := fmt.Sprintf(`%v/%v`, DirectoryPath, FileName)

		if !strings.Contains(AcceptedFiles, fileExt) {
			pkg_err.Log(err, "E1175309", pkg_err.DocumentTypeNotAccepted, AcceptedFiles,
				fmt.Sprintf(" >>> fileExt:%v", fileExt), fmt.Sprintf(" >>> fileType:%v", filetype))
			err = pkg_err.Take(err, "E1086937").
				Message(pkg_err.DocumentTypeNotAccepted).
				Custom(pkg_err.ValidationFailedErr).Build()
			return
		}

		buffByte := buff.Bytes()

		// check to know this document is single or its allow to multiple
		if !document_type.IsSingleDocs(docs.Type) {
			count, _ := s.Repo.CountByRefrenceIdType(params.Tx, docs.ReferenceID, docs.Type)
			if count != 0 {
				err = pkg_err.Take(err, "E1024310").
					Message(pkg_err.CanNotUploadMultipleDocument).
					Custom(pkg_err.ValidationFailedErr).Build()
				return
			}
		}

		// write file
		if err = ioutil.WriteFile(FilePath, buffByte, 0644); err != nil {
			err = pkg_err.Take(err, "E1116850").
				Message(pkg_err.CouldNotUploadDocument).
				Custom(pkg_err.ValidationFailedErr).Build()
			return
		}

		docPayload.ReferenceID = docs.ReferenceID
		docPayload.Type = docs.Type
		docPayload.Name = FileName

		_, err = s.Create(params.Tx, docPayload)
		if err != nil {
			err = pkg_err.Take(err, "E1029482").
				Message("document did not saved").
				Custom(pkg_err.BadRequestErr).Build()
			return
		}
	}

	if documents, err = s.GetByRefrenceIdType(params.Tx, docPayload.ReferenceID, docPayload.Type); err != nil {
		return
	}

	return
}

// DeleteAllDocsById to delete documents
func (s *BaseDocumentServ) DeleteAllDocs(params param.Param, ReferenceID uint, docType pkg_types.Enum) (err error) {

	var DirectoryPath string
	if DirectoryPath, _, err = s.DocumentPathFolder(docType); err != nil {
		pkg_err.Log(err, "E1150525", "document type not found")
		return
	}

	var documents []base_model.Document
	if documents, err = s.GetByRefrenceIdType(params.Tx, ReferenceID, docType); err != nil {
		pkg_err.Log(err, "E1148824", "document get by type and refrence id not found")
		return
	}

	// delete in database
	for _, v := range documents {
		if err = s.Repo.Delete(params.Tx, v); err != nil {
			err = pkg_err.Take(err, "E1000389").Message(pkg_err.CouldNotDeleteDocuments).Build()
			return
		}
	}

	// delete in system
	for _, v := range documents {
		oldFile := filepath.Join(DirectoryPath, v.Name)
		os.Remove(oldFile)
	}

	return
}

// DocumentPathFolder find folder path by checking type of doc
func (s *BaseDocumentServ) DocumentPathFolder(docType pkg_types.Enum) (path, AcceptedFiles string, err error) {

	switch docType {
	case document_type.Gift:
		path = document_type.GiftPath
		AcceptedFiles = document_type.AcceptedImage
		return
	default:
		err = pkg_err.New(pkg_err.InternalServerError, "E1162835").
			Message(pkg_err.InternalServerError).Build()
		return
	}
}
