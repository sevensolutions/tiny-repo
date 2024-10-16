package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Masterminds/semver/v3"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sevensolutions/tiny-repo/core"
	myMiddleware "github.com/sevensolutions/tiny-repo/middleware"
	"github.com/sevensolutions/tiny-repo/storage"
)

type Server struct {
	Storage storage.StorageAdapter
}

func (srv *Server) upload(c echo.Context) error {
	spec, err := core.ParseVersionSpecFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if spec.Latest {
		return echo.NewHTTPError(http.StatusBadRequest, "uploading to latest version is not allowed")
	}

	ctx := c.Request().Context()

	err = srv.Storage.Upload(ctx, spec, c, c.Request().Body)

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
			err = storage.Tidy(srv.Storage, spec.ArtifactSpec, tidyKeep, &spec)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	return c.NoContent(http.StatusOK)
}

func (srv *Server) download(c echo.Context) error {
	ctx := c.Request().Context()

	spec, err := core.ParseVersionSpecFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if spec.Latest {
		versions, err := storage.GetSortedVersions(srv.Storage, spec.ArtifactSpec)
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

	err = srv.Storage.Download(ctx, spec, c)
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

func (srv *Server) getVersions(c echo.Context) error {
	spec, err := core.ParseArtifactSpecFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	versions, err := storage.GetSortedVersions(srv.Storage, spec)
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

func (srv *Server) deleteVersion(c echo.Context) error {
	spec, err := core.ParseVersionSpecFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return srv.Storage.DeleteVersion(spec)
}

func (srv *Server) deleteArtifact(c echo.Context) error {
	spec, err := core.ParseArtifactSpecFromEcho(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = storage.Tidy(srv.Storage, spec, 0, nil)
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

func (srv *Server) Run() {
	printBanner()

	srv.Storage = createStorageAdapter()

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(echojwt.JWT([]byte(core.GetRequiredEnvVar("JWT_SECRET"))))
	e.Use(myMiddleware.ValidateAuth)

	e.GET("/:namespace/:name", srv.getVersions)
	e.GET("/:namespace/:name/:version/:filename", srv.download)
	e.GET("/:namespace/:name/:version", srv.download)

	e.PUT("/:namespace/:name/:version/:filename", srv.upload)
	e.PUT("/:namespace/:name/:version", srv.upload)

	e.DELETE("/:namespace/:name", srv.deleteArtifact)
	e.DELETE("/:namespace/:name/:version/:filename", srv.deleteVersion)
	e.DELETE("/:namespace/:name/:version", srv.deleteVersion)

	e.Logger.Fatal(e.Start(":8080"))
}

func printBanner() {
	println(`
  _____ _             ____                  
 |_   _(_)_ __  _   _|  _ \ ___ _ __   ___  
   | | | | '_ \| | | | |_) / _ \ '_ \ / _ \ 
   | | | | | | | |_| |  _ <  __/ |_) | (_) |
   |_| |_|_| |_|\__, |_| \_\___| .__/ \___/ 
                |___/          |_|          

  sevensolutions - TinyRepo
	`)
}
