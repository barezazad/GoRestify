package document_type

import "GoRestify/pkg/pkg_types"

// account status enum
const (
	Gift pkg_types.Enum = "gift"
)

// List account status list
var List = []pkg_types.Enum{
	Gift,
}

// SingleDocs all SingleDocs type
var SingleDocs = []pkg_types.Enum{
	Gift,
}

// accepted type of documents
const (
	AcceptedImage      string = ".png,.jpeg,.jpg,.svg"
	AcceptedVideo      string = ".mp4,.avi,.mwv,.flv,.mkv"
	AcceptedCoverage   string = ".kml"
	AcceptedVideoImage string = ".png,.jpeg,.jpg,.mp4,.avi,.mwv,.flv,.mkv"
	AcceptedPdfFile    string = ".pdf"
	AcceptedDocs       string = ".pdf,.png,.jpeg,.jpg,.svg,.gif,.tiff,.eps,.ai,.indd,.raw"
)

// document path folder
const (
	GiftPath string = "public/gift"
)

// IsSingleDocs to check is document type is single
func IsSingleDocs(docType pkg_types.Enum) (result bool) {
	for _, v := range SingleDocs {
		if v == docType {
			result = true
			return
		}
	}
	return
}
