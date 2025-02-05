SNAPSHOTS_DIR:=./test_snapshots
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-cover
# SNAPSHOT_PATH:=./test_snapshots/test-solid-color-cover
# SNAPSHOT_PATH:=./test_snapshots/test-gradient-cover
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-icon
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-icon
SNAPSHOT_PATH:=./test_snapshots/test-table-rows
# SNAPSHOT_PATH:=./test_snapshots/Anytype.WebPublish.20241217.112212.67
# SNAPSHOT_PATH:=./test_snapshots/test-three-column
# SNAPSHOT_PATH:=./test_snapshots/test-angle-brackets
# SNAPSHOT_PATH:=./test_snapshots/test-me

EXEC:=./bin/anytype-publish-renderer
TEMPL_VER:=$(shell cat go.mod | grep templ | cut -d' ' -f2)

.PHONY :

setup-go:
	@echo 'Setting up go modules...'
	@go mod download

build: setup-go deps
	npm run build
	templ generate -lazy
	go build -o $(EXEC) .

deps:
	echo $(TEMPL_VER)
	go install github.com/a-h/templ/cmd/templ@$(TEMPL_VER)

test: setup-go
	go test -v ./...

render: build
	templ generate -lazy
	ANYTYPE_PUBLISH_CSS_DEBUG=yesplease $(EXEC) $(SNAPSHOT_PATH) > index.html

render-all: build
	rm -f index.html
	templ generate -lazy
	for p in $(shell ls $(SNAPSHOTS_DIR)); do \
		ANYTYPE_PUBLISH_CSS_DEBUG=yesplease $(EXEC) $(SNAPSHOTS_DIR)/$$p > $$p.html; \
	done
