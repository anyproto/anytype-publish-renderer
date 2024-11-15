SNAPSHOT_DIR:=./snapshot_pb/
ROOT_ID:=bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui
EXEC:=./bin/anytype-publish-renderer
TEMPL_VER:=$(shell cat go.mod | grep templ | cut -d' ' -f2)

.PHONY :

setup-go:
	@echo 'Setting up go modules...'
	@go mod download

build: setup-go
	templ generate -lazy
	go build -o $(EXEC) .

deps:
	echo $(TEMPL_VER)
	go install github.com/a-h/templ/cmd/templ@$(TEMPL_VER)

test: setup-go
	go test -v ./...

render: build
	templ generate -lazy
	$(EXEC) $(SNAPSHOT_DIR) $(ROOT_ID) > index.html
