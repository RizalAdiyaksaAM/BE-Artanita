package cloudinary

import (
	"context"
	"errors"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadImage(ctx context.Context, input multipart.File, folder string) (string, error)
	UploadVideo(ctx context.Context, input multipart.File, folder string) (string, error)
	DeleteImage(ctx context.Context, imageURL string) error
}

type cloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService(cloudinary *cloudinary.Cloudinary) CloudinaryService {
	return &cloudinaryService{
		cloudinary: cloudinary,
	}
}

func (cs *cloudinaryService) UploadImage(ctx context.Context, input multipart.File, folder string) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}
	result, err := cs.cloudinary.Upload.Upload(ctx, input, uploadParams)
	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}

func (cs *cloudinaryService) UploadVideo(ctx context.Context, input multipart.File, folder string) (string, error) {
	uploadParams := uploader.UploadParams{
		Folder: folder,
	}
	result, err := cs.cloudinary.Upload.Upload(ctx, input, uploadParams)
	if err != nil {
		return "", err
	}
	return result.SecureURL, nil
}

func (cs *cloudinaryService) DeleteImage(ctx context.Context, imageURL string) error {
	publicID, err := extractPublicID(imageURL)
	if err != nil {
		return err
	}

	_, err = cs.cloudinary.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicID})
	return err
}

func extractPublicID(imageURL string) (string, error) {
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}
	pathSegments := strings.Split(parsedURL.Path, "/")
	if len(pathSegments) < 2 {
		return "", errors.New("invalid URL format")
	}

	// Extract the public ID from the URL path
	publicIDWithExt := pathSegments[len(pathSegments)-1]
	publicID := strings.TrimSuffix(publicIDWithExt, filepath.Ext(publicIDWithExt))
	folderPath := strings.Join(pathSegments[len(pathSegments)-3:len(pathSegments)-1], "/")
	return folderPath + "/" + publicID, nil
}
