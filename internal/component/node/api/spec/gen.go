package spec

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
// nolint:lll
//go:generate oapi-codegen -old-config-style -package spec -generate types,skip-prune -o types.gen.go  ../../../../../api/node/v1/openapi.yaml
// nolint:lll
//go:generate oapi-codegen -old-config-style -package spec -generate server           -o server.gen.go ../../../../../api/node/v1/openapi.yaml
