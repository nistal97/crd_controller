go mod tidy
go mod vendor
go build -v -o ./output/crd_controller ./cmd/crd_controller