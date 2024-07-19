package storage

import (
	"context"
	"encoding/json"
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

func (a *LocalDirectoryAdapter) Upload(ctx context.Context, spec core.ArtifactVersionSpec, echo echo.Context, source io.Reader) error {
	fullPath := ospath.Join(a.rootDirectory, spec.Namespace, spec.Name, spec.Version.String())

	err := os.MkdirAll(fullPath, 0777)
	if err != nil {
		return err
	}

	blobPath := ospath.Join(fullPath, "blob")

	f, err := os.Create(blobPath)
	if err != nil {
		return err
	}

	defer f.Close()

	io.Copy(f, source)

	metaPath := ospath.Join(fullPath, "meta.json")

	err = saveMeta(metaPath, core.BlobMeta{
		OriginalFilename: echo.Param("filename"),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *LocalDirectoryAdapter) Download(ctx context.Context, spec core.ArtifactVersionSpec, target echo.Context) error {
	fullPath := ospath.Join(a.rootDirectory, spec.Namespace, spec.Name, spec.Version.String())
	blobPath := ospath.Join(fullPath, "blob")
	metaPath := ospath.Join(fullPath, "meta.json")

	meta, err := readMeta(metaPath)
	if err != nil {
		return nil
	}

	filename := meta.OriginalFilename

	requestedFilename := target.Param("filename")

	if requestedFilename != "" {
		filename = requestedFilename
	}

	if filename == "" {
		filename = "blob"
	}

	target.Attachment(blobPath, filename)

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

func saveMeta(metaPath string, meta core.BlobMeta) error {
	jsonBytes, _ := json.MarshalIndent(meta, "", "  ")

	err := os.WriteFile(metaPath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
func readMeta(metaPath string) (core.BlobMeta, error) {
	meta := core.BlobMeta{}

	jsonBytes, err := os.ReadFile(metaPath)
	if err != nil {
		return meta, nil
	}

	err = json.Unmarshal(jsonBytes, &meta)
	if err != nil {
		return meta, err
	}

	return meta, nil
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
