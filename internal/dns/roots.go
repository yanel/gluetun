package dns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/files"
)

func (c *configurator) DownloadRootHints(ctx context.Context, uid, gid int) error {
	c.logger.Info("downloading root hints from %s", constants.NamedRootURL)
	content, status, err := c.client.Get(ctx, string(constants.NamedRootURL))
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP status code is %d for %s", status, constants.NamedRootURL)
	}
	return c.fileManager.WriteToFile(
		string(constants.RootHints),
		content,
		files.Ownership(uid, gid),
		files.Permissions(constants.UserReadPermission))
}

func (c *configurator) DownloadRootKey(ctx context.Context, uid, gid int) error {
	c.logger.Info("downloading root key from %s", constants.RootKeyURL)
	content, status, err := c.client.Get(ctx, string(constants.RootKeyURL))
	if err != nil {
		return err
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP status code is %d for %s", status, constants.RootKeyURL)
	}
	return c.fileManager.WriteToFile(
		string(constants.RootKey),
		content,
		files.Ownership(uid, gid),
		files.Permissions(constants.UserReadPermission))
}
