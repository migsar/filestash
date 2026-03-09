package plugin

import (
	. "github.com/mickael-kerjean/filestash/server/common"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_local"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_nop"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_backend_s3"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_handler_mcp"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_handler_site"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_image_ascii"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_image_c"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_license"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_metadata_sqlite"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_search_stateless"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_security_scanner"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_security_svg"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_starter_http"
	_ "github.com/mickael-kerjean/filestash/server/plugin/plg_plugin_automatic_login"
)

func init() {
	Hooks.Register.Onload(func() { Log.Debug("plugins loaded") })
}
