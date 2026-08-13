package response

import (
	"strconv"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/pkg_terms"
)

// NotBind use special custom_error for reduced it
func (r *Response) NotBind(err error, code, part string) {
	err = pkg_err.Take(err, code).
		Message(pkg_err.ErrorInBindingV, part).
		Custom(pkg_err.BindingErr).Build()

	r.Error(err).JSON()
}

// Bind is used to make it easier for binding items
func (r *Response) Bind(st interface{}, code, part string) (err error) {
	if err = r.Context.ShouldBindJSON(&st); err != nil {
		r.NotBind(err, code, part)
		return
	}

	return
}

// GetID returns the ID
func (r *Response) GetID(idIn, code, part string) (id uint, err error) {
	tmpID, err := strconv.ParseUint(idIn, 10, 64)
	if err != nil {
		err = pkg_err.Take(err, code).
			Message(pkg_err.InvalidXForX, pkg_terms.ID, part).
			Custom(pkg_err.ValidationFailedErr).Build()
		r.Error(err).JSON()
	}
	id = uint(tmpID)
	return
}
