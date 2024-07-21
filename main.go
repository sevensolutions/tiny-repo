package main

import (
	"log"
	"net/http"
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

	err = app.Storage.Upload(ctx, spec, c, src)

	if err != nil {
		return err
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
		go func() {
			err = storage.Tidy(app.Storage, spec.ArtifactSpec, tidyKeep)
			if err != nil {
				log.Println(err)
			}
		}()
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
		versions, err := storage.GetSortedVersions(app.Storage, spec.ArtifactSpec)
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

	err = app.Storage.Download(ctx, spec, c)
	if err != nil {
		return nil
	}

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

	versions, err := storage.GetSortedVersions(app.Storage, spec)
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

func (app *App) deleteVersion(c echo.Context) error {
	spec, err := core.ParseVersionSpec(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return app.Storage.DeleteVersion(spec)
}

func (app *App) deleteArtifact(c echo.Context) error {
	spec, err := core.ParseArtifactSpec(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = storage.Tidy(app.Storage, spec, 0)
	if err != nil {
		return err
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

	e.GET("/:namespace/:name", app.getVersions)
	e.GET("/:namespace/:name/:version/:filename", app.download)
	e.GET("/:namespace/:name/:version", app.download)

	e.PUT("/:namespace/:name/:version/:filename", app.upload)
	e.PUT("/:namespace/:name/:version", app.upload)

	e.DELETE("/:namespace/:name", app.deleteArtifact)
	e.DELETE("/:namespace/:name/:version/:filename", app.deleteVersion)
	e.DELETE("/:namespace/:name/:version", app.deleteVersion)

	e.Logger.Fatal(e.Start(":8080"))
}

func main() {
	godotenv.Load()

	app := new(App)
	app.Run()
}
