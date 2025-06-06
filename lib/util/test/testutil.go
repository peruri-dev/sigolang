package testutil

import (
	"mime/multipart"
)

func AddMultipartFields(writer *multipart.Writer, fields map[string]string) error {
	for key, value := range fields {
		if err := MultipartAddFormField(writer, key, value); err != nil {
			return err
		}
	}
	return nil
}

func MultipartAddFormField(writer *multipart.Writer, fieldname, fieldvalue string) error {
	part, err := writer.CreateFormField(fieldname)
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(fieldvalue))
	return err
}
