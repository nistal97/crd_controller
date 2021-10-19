export GOROOT=E:/dev/tools/go1.17.windows-amd64/go
${GOROOT}/bin/go mod tidy
${GOROOT}/bin/go mod vendor
${GOROOT}/bin/go build -v -o ./output/crd_controller ./cmd/crd_controller
