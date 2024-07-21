package storage

import (
	"log"

	"github.com/sevensolutions/tiny-repo/core"
)

func Tidy(storage StorageAdapter, artifactSpec core.ArtifactSpec, keep int) error {
	versions, err := GetSortedVersions(storage, artifactSpec)
	if err != nil {
		return err
	}

	for i, v := range versions {
		if i >= keep {
			log.Println("Deleting version", v)

			spec := core.ArtifactVersionSpec{
				ArtifactSpec: artifactSpec,
				Version:      v,
			}

			storage.DeleteVersion(spec)
		}
	}

	return nil
}
