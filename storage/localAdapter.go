package storage

import (
	"context"
	"io"
	"os"
	ospath "path"
	"path/filepath"

	"github.com/sevensolutions/tiny-repo/core"

	"github.com/Masterminds/semver/v3"
	"github.com/labstack/echo/v4"
)

type LocalDirectoryAdapter struct {
	rootDirectory string
}

func LocalDirectory() *LocalDirectoryAdapter {
	adapter := new(LocalDirectoryAdapter)
	adapter.rootDirectory = core.GetRequiredEnvVar("STORAGE_DIRECTORY")

	return adapter
}

func (a *LocalDirectoryAdapter) Upload(ctx context.Context, spec core.ArtifactVersionSpec, source io.Reader) error {
	fullPath := ospath.Join(a.rootDirectory, spec.Namespace, spec.Name, spec.Version.String())

	err := os.MkdirAll(fullPath, 0777)
	if err != nil {
		return err
	}

	fullPath = ospath.Join(fullPath, "blob")

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer f.Close()

	io.Copy(f, source)

	return nil
}

func (a *LocalDirectoryAdapter) Download(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context) error {
	fullPath := ospath.Join(a.rootDirectory, spec.Namespace, spec.Name, spec.Version.String(), "blob")

	fileName := target.QueryParam("name")

	if fileName == "" {
		fileName = "blob"
	}

	target.Attachment(fullPath, fileName)

	return nil
}

func (a *LocalDirectoryAdapter) GetVersions(artifactSpec core.ArtifactSpec) ([]*semver.Version, error) {
	fullPath := ospath.Join(a.rootDirectory, artifactSpec.Namespace, artifactSpec.Name)

	exists, err := folderExists(fullPath)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	versionFolders, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	result := make([]*semver.Version, len(versionFolders))

	for i, f := range versionFolders {
		v, err := semver.NewVersion(f.Name())
		if err == nil {
			result[i] = v
		}
	}

	return result, nil
}

func folderExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (a *LocalDirectoryAdapter) DeleteVersion(spec core.ArtifactVersionSpec) error {
	fullPath := ospath.Join(a.rootDirectory, spec.Namespace, spec.Name, spec.Version.String())

	return removeAll(fullPath)
}

func removeAll(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.RemoveAll(path)
	})
	if err != nil {
		return err
	}
	return nil
}
