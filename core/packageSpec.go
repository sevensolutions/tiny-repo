package core

import (
	"errors"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/labstack/echo/v4"
)

type ArtifactSpec struct {
	Namespace string
	Name      string
}

type ArtifactVersionSpec struct {
	ArtifactSpec
	Version *semver.Version
	Latest  bool
}

func ParseArtifactSpec(c echo.Context) (ArtifactSpec, error) {
	namespace := c.Param("namespace")
	name := c.Param("name")

	// TODOPEI: Better sanitizing

	if namespace == "" {
		return ArtifactSpec{}, errors.New("namespace must not be empty")
	}
	if strings.Contains(namespace, "/") {
		return ArtifactSpec{}, errors.New("namespace must not contain /")
	}
	if name == "" {
		return ArtifactSpec{}, errors.New("name must not be empty")
	}
	if strings.Contains(name, "/") {
		return ArtifactSpec{}, errors.New("name must not contain /")
	}

	return ArtifactSpec{
		Namespace: namespace,
		Name:      name,
	}, nil
}

func ParseVersionSpec(c echo.Context) (ArtifactVersionSpec, error) {
	artifactSpec, err := ParseArtifactSpec(c)
	if err != nil {
		return ArtifactVersionSpec{}, err
	}

	version := c.Param("version")

	if version == "latest" {
		return ArtifactVersionSpec{
			ArtifactSpec: artifactSpec,
			Version:      nil,
			Latest:       true,
		}, nil
	} else {
		parsedVersion, err := semver.NewVersion(version)
		if err != nil {
			return ArtifactVersionSpec{}, err
		}

		return ArtifactVersionSpec{
			ArtifactSpec: artifactSpec,
			Version:      parsedVersion,
			Latest:       false,
		}, nil
	}
}
