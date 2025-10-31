package plg_plugin_automatic_login

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"
	"os"

	. "github.com/mickael-kerjean/filestash/server/common"
	"github.com/mickael-kerjean/filestash/server/model"
)

var isInitialized = false

func init() {
	Hooks.Register.Onload(func() {
		Log.Debug("Automatic Login Plugin loaded...")
	})

	Hooks.Register.Middleware(func(h HandlerFunc) HandlerFunc {
		return func(app *App, res http.ResponseWriter, req *http.Request) {
      // Auto-init backend?
			if shouldInitBackend(app, req) {
				if err := initBackend(app, res, req); err != nil {
					Log.Error("Failed to init backend: %v", err)
				}
			}

			h(app, res, req)
		}
	})
}

func shouldInitBackend(app *App, req *http.Request) bool {
  path := req.URL.Path

  if (path == "/" && !isInitialized) {
    return true
  }

  return false
}

func initBackend(ctx *App, res http.ResponseWriter, req *http.Request) error {
  // Check if S3 backend is available
  if Backend.Get("s3") == (Nothing {}) {
    Log.Error("S3 backend is not registered")
    // Should it be error or just continue to manual login?
    return ErrInternal
  }

  session := map[string]string{
    "timestamp":          time.Now().Format(time.RFC3339),
		"type":               "s3",
		"access_key_id":      os.Getenv("FILESTASH_S3_ACCESS_KEY"),
		"secret_access_key":  os.Getenv("FILESTASH_S3_SECRET_ACCESS_KEY"),
		"endpoint":           os.Getenv("FILESTASH_S3_ENDPOINT"),
		"region":             "",
		"path":               "/",
		"advanced":           "on",
	}

  session["path"] = EnforceDirectory(session["path"])

  backend, err := model.NewBackend(ctx, session)
	if err != nil {
		Log.Debug("[auth] action=authenticate::newBackend err=%s", ferror(err))
		Log.Stdout("AUDIT action[fail] backend[%s] user[%s] target[%s]", session["type"], backendID(session), ip(req))
		return err
	}

  s, err := json.Marshal(session)
	if err != nil {
		Log.Debug("[auth] action=authenticate::marshall err=%s", ferror(err))
		return NewError(err.Error(), 500)
	}

  obfuscate, err := EncryptString(SECRET_KEY_DERIVATE_FOR_USER, string(s))
	if err != nil {
		Log.Debug("[auth] action=authenticate::encrypt err=%s", ferror(err))
		return NewError(err.Error(), 500)
	}

	// split session cookie if greater than 3800 bytes
	value_limit := 3800
	index := 0
	end := 0
	for {
		if len(obfuscate) >= (index+1)*value_limit {
			end = (index + 1) * value_limit
		} else {
			end = len(obfuscate)
		}

		http.SetCookie(res, applyCookieRules(&http.Cookie{
			Name:   CookieName(index),
			Value:  obfuscate[index*value_limit : end],
			MaxAge: 60 * Config.Get("general.cookie_timeout").Int(),
			Path:   COOKIE_PATH,
		}, req))

    if end == len(obfuscate) {
			break
		} else {
			Log.Debug("[auth] action=authenticate::obfuscate index=%d length=%d total=%d", index, len(obfuscate[index*value_limit:end]), len(obfuscate))
			index++
		}
	}

  if Config.Get("features.protection.iframe").String() != "" {
		res.Header().Set("bearer", obfuscate)
	}

	ctx.Backend = backend
	ctx.Session = session

	return nil
}

func applyCookieRules(cookie *http.Cookie, req *http.Request) *http.Cookie {
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	if Config.Get("features.protection.iframe").String() != "" {
		if f := req.Header.Get("Referer"); strings.HasPrefix(f, "https://") {
			cookie.Secure = true
			cookie.SameSite = http.SameSiteNoneMode
			cookie.Partitioned = true
		} else {
			Log.Warning("you are trying to access Filestash from a non secure origin ('%s') and with iframe enabled. Either use SSL or disable iframe from the admin console.", f)
		}
	}
	return cookie
}

func backendID(session map[string]string) string {
	return Hash(GenerateID(session)+session["path"], 20)
}

func username(session map[string]string) string {
	if session["username"] != "" {
		return strings.ReplaceAll(session["username"], " ", "+")
	} else if session["user"] != "" {
		return strings.ReplaceAll(session["user"], " ", "+")
	}
	return GenerateID(session)
}

func ip(req *http.Request) string {
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		if parts := strings.Split(xff, ","); len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	if xrip := req.Header.Get("X-Real-Ip"); xrip != "" {
		return xrip
	}
	host, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return req.RemoteAddr
	}
	return host
}

func ferror(err error) string {
	return strings.ReplaceAll(err.Error(), " ", "+")
}
