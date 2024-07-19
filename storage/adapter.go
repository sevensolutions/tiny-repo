package storage

import (
	"context"
	"io"

	"github.com/sevensolutions/tiny-repo/core"

	"github.com/Masterminds/semver/v3"
	"github.com/labstack/echo/v4"
)

type StorageAdapter interface {
	Upload(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context, source io.Reader) error
	Download(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context) error
	GetVersions(artifactSpec core.ArtifactSpec) ([]*semver.Version, error)
	DeleteVersion(spec core.ArtifactVersionSpec) error
}
