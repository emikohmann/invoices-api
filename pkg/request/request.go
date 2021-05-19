package request

import (
	"bytes"
	"github.com/emikohmann/invoices-api/pkg/apierrors"
	"io"
	"net/http"
	"strings"
)

func GetFile(request *http.Request, filename string, maxSize int64) ([]string, apierrors.ApiError) {
	if err := request.ParseMultipartForm(maxSize); err != nil {
		apiErr := apierrors.NewBadRequest(err.Error())
		return nil, apiErr
	}

	file, _, err := request.FormFile(filename)
	if err != nil {
		apiErr := apierrors.NewBadRequest(err.Error())
		return nil, apiErr
	}

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		apiErr := apierrors.NewInternalServer(err.Error())
		return nil, apiErr
	}

	return strings.Split(buffer.String(), "\n"), nil
}
