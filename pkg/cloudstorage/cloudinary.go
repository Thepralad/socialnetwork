package cloudstorage

import (
	"context"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

var cld *cloudinary.Cloudinary

func init() {
	var err error
	cld, err = cloudinary.NewFromParams("dps1mbifg", "887592117532741", "s8V7UCbfYepYky-H-SGDtTpMKhA")
	if err != nil {
		log.Fatalf("Cloudinary initialization failed: %v", err)
	}
}

func UploadImg(file multipart.File) (string, error) {
	resp, err := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Transformation: "w_500/f_auto,q_auto",
	})
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}