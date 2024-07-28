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

func ParseArtifactSpecFromEcho(c echo.Context) (ArtifactSpec, error) {
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

func ParseVersionSpecFromEcho(c echo.Context) (ArtifactVersionSpec, error) {
	artifactSpec, err := ParseArtifactSpecFromEcho(c)
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

func ParseArtifactSpec(value string) (ArtifactSpec, error) {
	parts := strings.Split(value, "/")

	namespace := parts[0]
	name := parts[1]

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

func ParseVersionSpec(value string) (ArtifactVersionSpec, error) {
	artifactSpec, err := ParseArtifactSpec(value)
	if err != nil {
		return ArtifactVersionSpec{}, err
	}

	parts := strings.Split(value, "/")

	if len(parts) != 3 {
		return ArtifactVersionSpec{}, errors.New("Invalid version " + value)
	}

	version := parts[2]

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
