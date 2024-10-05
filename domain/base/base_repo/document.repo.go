package base_repo

import (
	"reflect"

	"GoRestify/domain/base/base_model"
	"GoRestify/domain/base/base_term"
	"GoRestify/internal/core"

	"GoRestify/pkg/db_error"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_sql"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"
	"GoRestify/pkg/validator"
)

// DocumentRepo for injecting engine
type DocumentRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideDocumentRepo is used in wire and initiate the Cols
func ProvideDocumentRepo(engine *core.Engine) DocumentRepo {
	return DocumentRepo{
		Engine: engine,
		Cols:   pkg_sql.ColumnExtractor(reflect.TypeOf(base_model.Document{}), base_model.DocumentTable),
	}
}

// FindByID finds the document via its id
func (r *DocumentRepo) FindByID(tx tx.Tx, id uint) (document base_model.Document, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.DocumentTable).
		Where("id = ?", id).
		First(&document).Error

	err = db_error.Parse(err, base_term.Documents, validator.Find)
	return
}

// GetByRefrenceIdType get documents by refrence_id and type
func (r *DocumentRepo) GetByRefrenceIdType(tx tx.Tx, referenceID uint, docType pkg_types.Enum) (documents []base_model.Document, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.DocumentTable).
		Where("refrence_id = ? AND type = ?", referenceID, docType).
		Find(&documents).Error

	err = db_error.Parse(err, base_term.Documents, validator.Find)
	return
}

// CountByRefrenceIdType count the documents by refrence_id and type
func (r *DocumentRepo) CountByRefrenceIdType(tx tx.Tx, id uint, docType pkg_types.Enum) (count int64, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.DocumentTable).
		Where("refrence_id = ? AND type = ?", id, docType).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Documents, validator.Find)
	return
}

// List returns an array of documents
func (r *DocumentRepo) List(params param.Param) (documents []base_model.Document, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(r.Cols, params.Select); err != nil {
		err = pkg_err.Take(err, "E7138163").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E7120704").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.DocumentTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&documents).Error

	err = db_error.Parse(err, base_term.Documents, validator.List)
	return
}

// Count of documents, mainly calls with List
func (r *DocumentRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(r.Cols); err != nil {
		err = pkg_err.Take(err, "E7157340").Custom(pkg_err.ValidationFailedErr).Build()
		return
	}

	err = r.Engine.DB.Table(base_model.DocumentTable).
		Where(whereStr).
		Count(&count).Error

	err = db_error.Parse(err, base_term.Documents, validator.List)
	return
}

// Create a document
func (r *DocumentRepo) Create(tx tx.Tx, document base_model.Document) (u base_model.Document, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.DocumentTable).Create(&document).Scan(&u).Error

	err = db_error.Parse(err, base_term.Documents, validator.Create)
	return
}

// Save the document, in case it is not exist create it
func (r *DocumentRepo) Save(tx tx.Tx, document base_model.Document) (u base_model.Document, err error) {
	err = tx.GetDB(r.Engine.DB).Table(base_model.DocumentTable).Where("id = ?", document.ID).Save(&document).Find(&u).Error

	err = db_error.Parse(err, base_term.Documents, validator.Find)
	return
}

// Delete the document
func (r *DocumentRepo) Delete(tx tx.Tx, document base_model.Document) (err error) {
	err = tx.GetDB(r.Engine.DB).Unscoped().Table(base_model.DocumentTable).Delete(&document).Error

	err = db_error.Parse(err, base_term.Documents, validator.Delete)
	return
}

// DeleteByType the delete document by type
func (r *DocumentRepo) DeleteByType(tx tx.Tx, docType pkg_types.Enum) (err error) {
	err = tx.GetDB(r.Engine.DB).Exec("DELETE FROM base_documents where type = ?", docType).Error

	err = db_error.Parse(err, base_term.Documents, validator.Delete)
	return
}
