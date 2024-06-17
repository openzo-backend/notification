package utils

import (
	"context"
	"io"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func SaveFile(fileReader io.Reader, fileHeader *multipart.FileHeader) (string, error) {

	cld, _ := cloudinary.NewFromParams("dsvo76qzw", "963212776811178", "iV6eu41-Z993ncaBp76UjDWposs")
	var ctx = context.Background()
	resp, err := cld.Upload.Upload(ctx, fileReader, uploader.UploadParams{ // UploadParams is a struct that contains the parameters for the upload
		PublicID: fileHeader.Filename, // PublicID is the name of the file
	})
	if err != nil {
		log.Fatalf("Failed to upload: %v", err)
		return "", err
	}

	return resp.SecureURL, nil

}
