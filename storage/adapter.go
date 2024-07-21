package storage

import (
	"context"
	"io"
	"sort"

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

func GetSortedVersions(storage StorageAdapter, artifactSpec core.ArtifactSpec) ([]*semver.Version, error) {
	versions, err := storage.GetVersions(artifactSpec)
	if err != nil {
		return nil, err
	}

	if versions == nil {
		return []*semver.Version{}, nil
	}

	versions = core.FilterArray(versions, func(x *semver.Version) bool {
		return x != nil
	})

	sort.Slice(versions, func(d1, d2 int) bool {
		return versions[d1].Compare(versions[d2]) > 0
	})

	return versions, nil
}
