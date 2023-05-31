package node_client

//go:generate go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
// nolint:lll
//go:generate oapi-codegen -old-config-style -package node_client -generate types,skip-prune -o types.gen.go  ../../../../../api/node/v1/openapi.yaml
// nolint:lll
//go:generate oapi-codegen -old-config-style -package node_client -generate client           -o client.gen.go ../../../../../api/node/v1/openapi.yaml
