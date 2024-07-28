package storage

import (
	"context"
	"io"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/labstack/echo/v4"
	"github.com/sevensolutions/tiny-repo/core"
)

type testStorage struct {
}

func TestTidy(t *testing.T) {
	storage := new(testStorage)

	Tidy(storage, core.ArtifactSpec{
		Namespace: "",
		Name:      "",
	}, 3, nil)
}

func (s *testStorage) Upload(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context, source io.Reader) error {
	return nil
}
func (s *testStorage) Download(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context) error {
	return nil
}
func (s *testStorage) GetVersions(artifactSpec core.ArtifactSpec) ([]*semver.Version, error) {
	versions := []*semver.Version{
		semver.New(1, 0, 0, "", ""),
	}

	return versions, nil
}
func (s *testStorage) DeleteVersion(spec core.ArtifactVersionSpec) error {
	return nil
}
