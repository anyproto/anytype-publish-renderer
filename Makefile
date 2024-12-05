SNAPSHOT_DIR:=./test_snapshots/snapshot_pb6_with_embeds/
ROOT_ID:=bafyreiecs2ivic6ne2lvkohrf3ojqizngmsk5ywilexo6xhtzlqfgngp64
EXEC:=./bin/anytype-publish-renderer
TEMPL_VER:=$(shell cat go.mod | grep templ | cut -d' ' -f2)

.PHONY :

setup-go:
	@echo 'Setting up go modules...'
	@go mod download

build: setup-go deps
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
