package plg_plugin_example

import (
  "fmt"
  "net/http"
  "time"
  "strconv"
	"context"

  . "github.com/mickael-kerjean/filestash/server/common"
  "github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
)

type S3Backend struct {
	client     *s3.S3
	config     *aws.Config
	params     map[string]string
	Context    context.Context
	threadSize int
	timeout    time.Duration
}

func init() {
  Hooks.Register.Onload(func() {
    Log.Debug("Woopi! ")
  })

  Hooks.Register.HttpEndpoint(func(r *mux.Router, app *App) error {
    // Not working since I copied everything from plg_backend_s3 here
    // if Backend.Get("s3") == nil {
    //   Log.Error("S3 backend is not registered, plg_plugin_example won't work properly")
    //   return ErrInternal
    // }

    params := map[string]string{
      "timestamp": "2025-10-16T09:39:37Z",
      "access_key_id": "X01W77VXAMEPUWXVEP0Z",
      "advanced": "on",
      "session_token": "",
      "secret_access_key": "EL4gG5Yc5ibCeiVBvFse9hIZJ7msy1viYDi3Tm7W",
      "endpoint": "us-southeast-1.linodeobjects.com",
      "role_arn": "",
      "path": "/",
      "type": "s3",
      "region": "",
      "encryption_key": "",
      "number_thread": "",
      "timeout": "",
    }

    creds := []credentials.Provider{}
    if params["access_key_id"] != "" || params["secret_access_key"] != "" {
      creds = append(creds, &credentials.StaticProvider{Value: credentials.Value{
        AccessKeyID:     params["access_key_id"],
        SecretAccessKey: params["secret_access_key"],
      }})
    }

    config := &aws.Config{
      Credentials:                   credentials.NewChainCredentials(creds),
      CredentialsChainVerboseErrors: aws.Bool(true),
      S3ForcePathStyle:              aws.Bool(true),
      Region:                        aws.String("auto"),
    }
    if params["endpoint"] != "" {
      config.Endpoint = aws.String(params["endpoint"])
    }

    var timeout time.Duration
    if secs, err := strconv.Atoi(params["timeout"]); err == nil {
      timeout = time.Duration(secs) * time.Second
    }

    threadSize, err := strconv.Atoi(params["number_thread"])
    if err != nil {
      threadSize = 50
    } else if threadSize > 5000 || threadSize < 1 {
      threadSize = 2
    }


    backend := &S3Backend{
      config:     config,
      params:     params,
      client:     s3.New(session.New(config)),
      Context:    app.Context,
      threadSize: threadSize,
      timeout:    timeout,
    }

    fmt.Printf("%v", backend)
    Log.Debug("Wooa!")
    return nil
  })


  Hooks.Register.Middleware(func(h HandlerFunc) HandlerFunc {
    return func(app *App, res http.ResponseWriter, req *http.Request) {
      Log.Debug("This should run on each request!")

      http.SetCookie(res, &http.Cookie{
        Name:   "Patito",
        Value:  "Cuac",
        MaxAge: 3600,
        Path:   "/",
      })

      h(app, res, req)
    }
  })
}
