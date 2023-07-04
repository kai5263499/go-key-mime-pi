// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package v1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// GUI landing page
	// (GET /)
	GetHome(ctx echo.Context) error

	// (GET /socket.io)
	GetWebSocketConnect(ctx echo.Context) error

	// (POST /v1/paste)
	PostPaste(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetHome converts echo context to params.
func (w *ServerInterfaceWrapper) GetHome(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetHome(ctx)
	return err
}

// GetWebSocketConnect converts echo context to params.
func (w *ServerInterfaceWrapper) GetWebSocketConnect(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetWebSocketConnect(ctx)
	return err
}

// PostPaste converts echo context to params.
func (w *ServerInterfaceWrapper) PostPaste(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostPaste(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/", wrapper.GetHome)
	router.GET(baseURL+"/socket.io", wrapper.GetWebSocketConnect)
	router.POST(baseURL+"/v1/paste", wrapper.PostPaste)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/6xSTW8TQQz9K5bP23yAuOyxCJXcKkrVA+Lg3XE30+yOhxkndFXlvyPPtqQlK7iQQ6Jx",
	"nv3e8/MT+nAvWD+h49wmH9VLwBq/8CDK/QitBE3SA0GiHBtOaYToYcdjI5QcdOQ6VqxQvfaMNXZysePx",
	"YvADX0SPFR445WnoerFarPBYoUQOFD3W+L6UKoyk22wylvZlE88k3XEDDWV2J/Kr2w30FJwPHUTqGMvo",
	"RNaxcVjjFetnGayeOEcJmQvJu9XKfswbh0Kl/KjLrQ69PXK75YFKeYzmKWvyocOjfSrM+2GgNNr4c36l",
	"LmP9DQ9r/G7gZZZ2x7rw8srYmcg7bm4K7qOEwK3+KXi9Ws8uZGqCfewSObbVPlt7i3ye6iVA3rcts2Nn",
	"4A9z4E1QToF6+JSSJDQTZ6YO62WkrGzdUfJMWtf2NxBMqwMV0C3/Tu4sqGvJWlqK8x97znopbpxLKfbk",
	"g734kYZYbu7r1mfwGQiaXtodyD0YFDbwk4IaeebgXkTc3lyezvYvWU9KfGKHtaY9H/95RRRj79viavmQ",
	"Jbw9ppjMs/qp+xTEayvGU71IaUR6pjBFMJWkebDzKOrmFv4f4j0efwUAAP//A7xjuRUEAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}