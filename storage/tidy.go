package storage

import (
	"log"

	"github.com/sevensolutions/tiny-repo/core"
)

func Tidy(storage StorageAdapter, artifactSpec core.ArtifactSpec, keep int, belowVersion *core.ArtifactVersionSpec) error {
	versions, err := GetSortedVersions(storage, artifactSpec)
	if err != nil {
		return err
	}

	i := 0

	for _, v := range versions {
		if belowVersion != nil && v.Compare(belowVersion.Version) > 0 {
			continue
		}

		if i >= keep {
			log.Println("Deleting version", v)

			spec := core.ArtifactVersionSpec{
				ArtifactSpec: artifactSpec,
				Version:      v,
			}

			storage.DeleteVersion(spec)
		}

		i++
	}

	return nil
}
