package v1

//go:generate oapi-codegen -package v1 -include-tags=v1 -generate spec,server -o server.gen.go -old-config-style ../../api/v1.yaml
//go:generate oapi-codegen -package v1 -include-tags=v1 -generate types -o types.gen.go -old-config-style ../../api/v1.yaml
