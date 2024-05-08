package nero

// OpenAPI server generation has template overrides at ./server/api/templates

// Changes in the templates:
// - expose the http.Request for strict response visit

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 --config ./server/api/nekos/v2/models.cfg.yaml -o ./server/api/nekos/v2/models.gen.go ./server/api/schema/nekos/v2.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 --config ./server/api/nekos/v2/server.cfg.yaml --templates ./server/api/templates -o ./server/api/nekos/v2/server.gen.go ./server/api/schema/nekos/v2.yaml

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 --config ./server/api/v1/models.cfg.yaml -o ./server/api/v1/models.gen.go ./server/api/schema/v1.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 --config ./server/api/v1/client.cfg.yaml -o ./server/api/v1/client.gen.go ./server/api/schema/v1.yaml
//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0 --config ./server/api/v1/server.cfg.yaml --templates ./server/api/templates -o ./server/api/v1/server.gen.go ./server/api/schema/v1.yaml
