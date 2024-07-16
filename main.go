package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/sevensolutions/tiny-repo/core"
	myMiddleware "github.com/sevensolutions/tiny-repo/middleware"
	"github.com/sevensolutions/tiny-repo/storage"

	"github.com/Masterminds/semver/v3"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type App struct {
	Storage storage.StorageAdapter
}

func (app *App) upload(c echo.Context) error {
	spec, err := core.ParseVersionSpec(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if spec.Latest {
		return echo.NewHTTPError(http.StatusBadRequest, "uploading to latest version is not allowed")
	}

	tidyKeep := 0
	keepParam := c.QueryParam("keep")
	if keepParam != "" {
		keep, err := strconv.ParseInt(keepParam, 10, 32)
		if err != nil || keep < 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid keep parameter")
		}
		tidyKeep = int(keep)
	}

	if tidyKeep > 0 {
		err = app.tidyVersions(spec.ArtifactSpec, tidyKeep)
		if err != nil {
			return err
		}
	}

	ctx := c.Request().Context()

	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing file")
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	err = app.Storage.Upload(ctx, spec, src)

	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (app *App) download(c echo.Context) error {
	ctx := c.Request().Context()

	spec, err := core.ParseVersionSpec(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if spec.Latest {
		versions, err := app.getSortedVersions(spec.ArtifactSpec)
		if err != nil {
			return err
		}

		if len(versions) > 0 {
			spec = core.ArtifactVersionSpec{
				ArtifactSpec: spec.ArtifactSpec,
				Version:      versions[0],
				Latest:       false,
			}
		} else {
			return c.NoContent(http.StatusNotFound)
		}
	}

	app.Storage.Download(ctx, spec, c)

	return c.NoContent(http.StatusNotFound)
}

type GetVersionsResponse struct {
	Count    int      `json:"count"`
	Latest   string   `json:"latest"`
	Versions []string `json:"versions"`
}

func (app *App) getVersions(c echo.Context) error {
	spec, err := core.ParseArtifactSpec(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	versions, err := app.getSortedVersions(spec)
	if err != nil {
		return err
	}

	if len(versions) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	response := &GetVersionsResponse{
		Count:  len(versions),
		Latest: versions[0].String(),
		Versions: core.MapArray(versions, func(v *semver.Version) string {
			return v.String()
		}),
	}

	return c.JSON(http.StatusOK, response)
}

func (app *App) getSortedVersions(artifactSpec core.ArtifactSpec) ([]*semver.Version, error) {
	versions, err := app.Storage.GetVersions(artifactSpec)
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

func (app *App) tidyVersions(artifactSpec core.ArtifactSpec, keep int) error {
	versions, err := app.getSortedVersions(artifactSpec)
	if err != nil {
		return err
	}

	for _, v := range versions {

		if keep <= 1 {
			log.Println("Deleting version", v)

			spec := core.ArtifactVersionSpec{
				ArtifactSpec: artifactSpec,
				Version:      v,
			}

			app.Storage.DeleteVersion(spec)
		} else {
			keep--
		}
	}

	return nil
}

func createStorageAdapter() storage.StorageAdapter {
	adapterType := core.GetRequiredEnvVar("STORAGE_TYPE")

	switch adapterType {
	case "Local":
		return storage.LocalDirectory()
	case "S3":
		return storage.MinIO()
	default:
		panic("Invalid STORAGE_TYPE. Only Local or S3 are supported.")
	}
}

func (app *App) Run() {

	app.Storage = createStorageAdapter()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.JWT([]byte(core.GetRequiredEnvVar("JWT_SECRET"))))
	e.Use(myMiddleware.ValidateAuth)

	e.GET("/:namespace/:name/:version", app.download)
	e.PUT("/:namespace/:name/:version", app.upload)
	e.GET("/:namespace/:name", app.getVersions)

	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	godotenv.Load()

	app := new(App)
	app.Run()
}
