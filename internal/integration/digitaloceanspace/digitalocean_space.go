package integration

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/infrastructure/config"
	"codebase-app/internal/integration/digitaloceanspace/entity"
	"codebase-app/pkg"
	"codebase-app/pkg/errmsg"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/rs/zerolog/log"
)

type DigitaloceanSpaceContract interface {
	UploadFile(ctx context.Context, req *entity.UploadFileRequest) (entity.UploadFileResponse, error)
	DeleteFile(ctx context.Context, req *entity.DeleteFileRequest) error
	ListFiles(ctx context.Context) ([]types.Object, error)
}

type dospace struct {
	storage *s3.Client
}

func NewDigitalOceanSpaceIntegration() DigitaloceanSpaceContract {
	return &dospace{
		storage: adapter.Adapters.ShopeefunStorage,
	}
}

func (d *dospace) UploadFile(ctx context.Context, req *entity.UploadFileRequest) (entity.UploadFileResponse, error) {
	var res = entity.UploadFileResponse{}

	if req.File == nil {
		return res, errmsg.NewCustomErrors(400, errmsg.WithErrors("file", "file is required."))
	}

	var (
		filename = pkg.SanitizeFilename(req.File.Filename, true)
		uploader = manager.NewUploader(d.storage)
	)

	f, err := req.File.Open()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("integration::dospace-UploadFile Error while opening file")
		return res, err
	}
	defer f.Close()

	result, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(config.Envs.ShopeefunStorage.Bucket),
		Key:    aws.String(filename),
		Body:   f,
		ACL:    types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("integration::dospace-UploadFile Error while uploading file")
		return res, err
	}

	res.FileName = filename
	res.Url = result.Location

	return res, nil
}

func (d *dospace) DeleteFile(ctx context.Context, req *entity.DeleteFileRequest) error {
	_, err := d.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(config.Envs.ShopeefunStorage.Bucket),
		Key:    aws.String(req.FileName),
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("integration::dospace-DeleteFile Error while deleting file")
		return err
	}

	return nil
}

func (d *dospace) ListFiles(ctx context.Context) ([]types.Object, error) {
	objects := []types.Object{}
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(config.Envs.ShopeefunStorage.Bucket),
	}

	paginator := s3.NewListObjectsV2Paginator(d.storage, input)

	// Iterate through the pages of results
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			log.Error().Err(err).Msg("integration::dospace-ListFiles failed to get page of results")
			return objects, errmsg.NewCustomErrors(500, errmsg.WithMessage("failed to get page of results"))
		}

		objects = append(objects, page.Contents...)
	}

	return objects, nil
}
