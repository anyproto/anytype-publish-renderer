SNAPSHOT_PATH:=./test_snapshots/Anytype.WebPublish.20241217.112212.67
TEST_PATH:=./test_snapshots/test
# SNAPSHOT_PATH:=http://localhost:8017
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
	$(EXEC) $(TEST_PATH) > index1.html
