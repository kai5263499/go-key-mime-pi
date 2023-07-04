package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/felixge/fgprof"
	socketio "github.com/googollee/go-socket.io"
	gokeymimepi "github.com/kai5263499/go-key-mime-pi"
	"github.com/kai5263499/go-key-mime-pi/internal/domain"
	v1 "github.com/kai5263499/go-key-mime-pi/internal/v1"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag"
)

type Api interface {
	Start() error
	Stop() error
	GetHome(ctx echo.Context) error
	PostPaste(ctx echo.Context) error
	GetWebSocketConnect(ctx echo.Context) error
}

var _ v1.ServerInterface = (*api)(nil)
var _ v1.ServerInterface = (Api)(nil)
var _ Api = (*api)(nil)

type api struct {
	shutdownContext    context.Context
	shutdownCancelFunc context.CancelFunc
	cfg                *domain.Config
	echo               *echo.Echo
	socketIOServer     *socketio.Server
	hid                gokeymimepi.Hid
}

func New(shutdownContext context.Context,
	shutdownCancelFunc context.CancelFunc,
	cfg *domain.Config,
	h gokeymimepi.Hid) (*api, error) {

	e := echo.New()
	e.Debug = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	a := &api{
		shutdownContext:    shutdownContext,
		shutdownCancelFunc: shutdownCancelFunc,
		cfg:                cfg,
		echo:               e,
		hid:                h,
	}

	all := v1.ServerInterfaceWrapper{
		Handler: a,
	}

	e.GET("/", all.GetHome)

	e.Static("/assets", "static")

	if cfg.EnableProfiling {
		e.GET("/debug/fgprof", func(c echo.Context) error {
			fgprof.Handler().ServeHTTP(c.Response().Writer, c.Request())
			return nil
		})
	}

	if cfg.EnableSwagger {
		if err := a.setupSwagger(); err != nil {
			return nil, err
		}
	}

	e.POST("/v1/paste", all.PostPaste)

	if err := a.setupSocketIO(); err != nil {
		logrus.WithError(err).Fatal("error setting up socketio server")
	}

	return a, nil
}

func (a *api) setupSwagger() error {
	swagger, err := v1.GetSwagger()
	if err != nil {
		return err
	}

	swaggerJson, err := swagger.MarshalJSON()
	if err != nil {
		return err
	}

	var SwaggerInfo = &swag.Spec{
		Version:          "",
		Host:             "",
		BasePath:         "",
		Schemes:          []string{},
		Title:            "",
		Description:      "",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(swaggerJson),
	}

	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)

	a.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	return nil
}

func (a *api) setupSocketIO() error {

	a.socketIOServer = socketio.NewServer(nil)

	a.socketIOServer.OnConnect("/", func(s socketio.Conn) error {
		logrus.Debugf("connected: %s", s.ID())
		return nil
	})
	a.socketIOServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		logrus.Debugf("disconnected: %s reason: %s", s.ID(), reason)
		return
	})
	a.socketIOServer.OnEvent("/", "keystroke", func(s socketio.Conn, keyEvent gokeymimepi.JSKeyEvent) {
		logrus.Debugf("received keystroke: %+#v", keyEvent)

		controlKeys, keycode, err := gokeymimepi.ConvertJSKeycodes(keyEvent)
		if err != nil {
			logrus.WithError(err).Error("error converting key event")
			s.Emit("keystroke-received", gokeymimepi.JSKeyEventResponse{Success: false})
			return
		}

		a.hid.Send(a.cfg.HidPath, byte(controlKeys), []byte{byte(keycode)})

		s.Emit("keystroke-received", gokeymimepi.JSKeyEventResponse{Success: true})
	})

	a.socketIOServer.OnError("/", func(s socketio.Conn, e error) {
		logrus.Debugf("socket error: %s", e.Error())
	})

	go func() {
		if err := a.socketIOServer.Serve(); err != nil {
			logrus.WithError(err).Fatal("socketio server error")
		}
	}()

	a.echo.Any("/socket.io/", echo.WrapHandler(a.socketIOServer))

	return nil
}

func (a *api) Start() (err error) {
	go a.httpListenerLoop()
	return
}

func (a *api) Stop() (err error) {
	a.shutdownCancelFunc()
	a.socketIOServer.Close()
	return
}

func (a *api) httpListenerLoop() {
	a.echo.Logger.Fatal(a.echo.Start(fmt.Sprintf(":%d", a.cfg.HttpPort)))
}

func (a *api) GetHome(ctx echo.Context) (err error) {
	file, err := ioutil.ReadFile("static/index.html")
	if err != nil {
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	return ctx.HTML(http.StatusOK, string(file))
}

// (POST /v1/paste)
func (a *api) PostPaste(ctx echo.Context) error {

	reqBody := ctx.Request().Body
	defer reqBody.Close()

	body, err := ioutil.ReadAll(reqBody)
	if err != nil {
		logrus.WithError(err).Error("error reading request body")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}

	logrus.Debugf("sending keystrokes: %s", string(body))
	a.hid.SendString(a.cfg.HidPath, string(body))

	return nil
}

// (GET /v1/ws)
func (a *api) GetWebSocketConnect(ctx echo.Context) error {
	return nil
}
