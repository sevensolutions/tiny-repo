package storage

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/sevensolutions/tiny-repo/core"

	"github.com/Masterminds/semver/v3"
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioAdapter struct {
	client     *minio.Client
	bucketName string
}

func MinIO() *MinioAdapter {
	adapter := new(MinioAdapter)

	endpoint := core.GetRequiredEnvVar("S3_ENDPOINT")
	accessKeyID := core.GetRequiredEnvVar("S3_ACCESSKEY")
	secretAccessKey := core.GetRequiredEnvVar("S3_SECRETKEY")
	useSSL := core.GetRequiredEnvVarBool("S3_USESSL")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		panic(err)
	}

	adapter.client = minioClient
	adapter.bucketName = core.GetRequiredEnvVar("S3_BUCKETNAME")

	return adapter
}

func (a *MinioAdapter) Upload(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context, source io.Reader) error {
	objectName := spec.Namespace + "/" + spec.Name + spec.Version.String() + "/blob"

	_, err := a.client.PutObject(ctx, a.bucketName, objectName, source, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return err
	}

	return nil
}

func (a *MinioAdapter) Download(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context) error {
	objectName := spec.Namespace + "/" + spec.Name + spec.Version.String() + "/blob"

	signedUrl, err := a.client.Presign(ctx, "GET", a.bucketName, objectName, time.Duration(5)*time.Minute, url.Values{})

	if err != nil {
		log.Println(err)

		return target.String(http.StatusNotFound, "Not found")
	}

	return target.Redirect(http.StatusTemporaryRedirect, signedUrl.String())
}

func (a *MinioAdapter) GetVersions(artifactSpec core.ArtifactSpec) ([]*semver.Version, error) {
	panic("Not implemented")
}

func (a *MinioAdapter) DeleteVersion(spec core.ArtifactVersionSpec) error {
	panic("Not implemented")
}
