SNAPSHOT_DIR:=./snapshot_pb/
ROOT_ID:=bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui
EXEC:=./bin/anytype-publish-renderer

.PHONY :

setup-go:
	@echo 'Setting up go modules...'
	@go mod download

build: setup-go
	templ generate -lazy
	go build -o $(EXEC) .

test: setup-go
	go test -v ./...

render: build
	templ generate -lazy
	$(EXEC) $(SNAPSHOT_DIR) $(ROOT_ID) > index.html
