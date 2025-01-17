# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-cover
# SNAPSHOT_PATH:=./test_snapshots/test-solid-color-cover
# SNAPSHOT_PATH:=./test_snapshots/test-gradient-cover
SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-emoji-cover

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
	$(EXEC) $(SNAPSHOT_PATH) > index.html
