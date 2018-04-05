package internal

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/IngCr3at1on/x/build"
	"github.com/IngCr3at1on/x/carnival/config"
	pb "github.com/IngCr3at1on/x/carnival/proto"
	"github.com/labstack/echo"
)

const dockerfile = `FROM progrium/busybox
COPY app /bin/app`

// SetupRoutes sets up the the API routes.
func SetupRoutes(cfg *config.Config, g *echo.Group) {
	if cfg == nil {
		panic("cfg cannot be nil")
	}

	g.GET("/test", makeTestHandler(cfg))
	g.POST("", makeNewJobHandler(cfg))
}

func makeTestHandler(cfg *config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusTeapot)
	}
}

func makeNewJobHandler(cfg *config.Config) echo.HandlerFunc {
	// TODO: implement timeout.
	return func(c echo.Context) error {
		var job pb.Job
		if err := c.Bind(&job); err != nil {
			msg := "malformed request body"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}

		ctx := c.Request().Context()

		dir, err := tempDir(ctx)
		if err != nil {
			msg := "couldn't create temporary environment"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusBadRequest, msg)
		}
		defer func() {
			if err := os.RemoveAll(dir); err != nil {
				fmt.Fprint(os.Stderr, err.Error())
			}
		}()

		if _, err := os.Create(fmt.Sprintf("%s/go.mod", dir)); err != nil {
			msg := "error prepping build environment"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		_file := fmt.Sprintf("%s/main.go", dir)
		if err := ioutil.WriteFile(_file, []byte(job.Body), 0700); err != nil {
			msg := "error writing main.go to build environment"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		var _out bytes.Buffer
		if err := build.Exec(&_out, "go", "fmt", _file); err != nil {
			msg := "error prepping build environment"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		_out = bytes.Buffer{}
		if err := build.Exec(&_out, "go", "build", "-o", fmt.Sprintf("%s/app", dir), _file); err != nil {
			msg := "error building application"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusUnprocessableEntity, msg)
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/Dockerfile", dir), []byte(dockerfile), 0700); err != nil {
			msg := "error writing dockerfile"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		_out = bytes.Buffer{}
		if err := build.Exec(&_out, "docker", "build", "-q", "-t", "carnival-app-latest", dir); err != nil {
			msg := "error building docker image"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		out, err := runDockerContainer(ctx, "carnival-app-latest")
		if err != nil {
			msg := "error running docker container"
			fmt.Fprintf(os.Stderr, "%s : %s", msg, err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		job.Body = out

		return c.JSON(http.StatusOK, &job)
	}
}
