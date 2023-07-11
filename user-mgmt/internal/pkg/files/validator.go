package files

import (
	"bytes"
	"github.com/unidoc/unipdf/v3/model"
)

func CheckPDFCorrupted(content []byte) error {
	reader := bytes.NewReader(content)
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return err
	}

	_, err = pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	return nil
}
