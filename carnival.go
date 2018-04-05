package carnival

import (
	"context"
	"errors"

	"github.com/IngCr3at1on/x/carnival/config"
	"github.com/IngCr3at1on/x/carnival/internal"
	"github.com/labstack/echo"
)

// Start blocking until ctx is cancelled or until a fatal error occurs.
func Start(ctx context.Context, cfg config.Config) error {
	if cfg.Address == "" {
		return errors.New("cfg.Address cannot be empty")
	}

	e := echo.New()
	g := e.Group("/api")

	internal.SetupRoutes(&cfg, g)

	return e.Start(cfg.Address)
}
