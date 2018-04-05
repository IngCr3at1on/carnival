package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

func tempDir(ctx context.Context) (string, error) {
	_uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(err, "uuid.NewV4")
	}

	dir := fmt.Sprintf("%s/carnival_%s", os.TempDir(), _uuid.String())
	return dir, os.Mkdir(dir, os.ModeDir|0777)
}
