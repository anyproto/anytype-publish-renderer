SNAPSHOTS_DIR:=./test_snapshots
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-cover
# SNAPSHOT_PATH:=./test_snapshots/test-solid-color-cover
# SNAPSHOT_PATH:=./test_snapshots/test-gradient-cover
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-icon
# SNAPSHOT_PATH:=./test_snapshots/test-uploaded-image-icon
# SNAPSHOT_PATH:=./test_snapshots/test-table-rows
# SNAPSHOT_PATH:=./test_snapshots/Anytype.WebPublish.20241217.112212.67
# SNAPSHOT_PATH:=./test_snapshots/test-three-column
# SNAPSHOT_PATH:=./test_snapshots/test-angle-brackets
# SNAPSHOT_PATH:=./test_snapshots/test-me
# SNAPSHOT_PATH:=https://anytype-prod-publishserver.s3.eu-central-2.amazonaws.com/67aa2b841b625f4939f38013
# SNAPSHOT_PATH:=https://anytype-prod-publishserver.s3.eu-central-2.amazonaws.com/67ab78a0c755e143dabb9891
SNAPSHOT_PATH := https://anytype-prod-publishserver.s3.eu-central-2.amazonaws.com/67accb61888bdb1bed39e762

EXEC:=./bin/anytype-publish-renderer
TEMPL_VER:=$(shell cat go.mod | grep templ | cut -d' ' -f2)
GO_FILES_CMD = find . -type f -name '*.go' -not -name '*_templ.go'

.PHONY :

setup-go:
	@echo 'Setting up go modules...'
	@go mod download

deps-goimport:
	go install golang.org/x/tools/cmd/goimports@latest

deps-linter:
	go install github.com/daixiang0/gci@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4

deps: deps-goimport deps-linter
	echo $(TEMPL_VER)
	go install github.com/a-h/templ/cmd/templ@$(TEMPL_VER)

check-fmt:
	@GO_FILES=$$($(GO_FILES_CMD)); \
	GOIMPORTS_OUTPUT=$$(goimports -d -l $$GO_FILES); \
	if [ -n "$$GOIMPORTS_OUTPUT" ]; then \
		echo "The following files have improperly ordered imports. Please run 'make fmt' to fix them:"; \
		echo "$$GOIMPORTS_OUTPUT"; \
		exit 1; \
	fi

fmt:
	@GO_FILES=$$($(GO_FILES_CMD)); \
	goimports -w -l $$GO_FILES

lint:
	golangci-lint run -v ./... --new-from-rev=origin/main --timeout 15m --verbose

build-templ:
	templ generate -lazy

build-js-css:
	npm run build

build-go:
	go build -o $(EXEC) .

build: setup-go deps build-js-css build-templ build-go

test: setup-go
	go test -v ./...

render-no-js-css: build-templ build-go
	$(EXEC) $(SNAPSHOT_PATH) > index.html

render: build
	$(EXEC) $(SNAPSHOT_PATH) > index.html

clean-html:
	rm *.html

render-all: build
	templ generate -lazy
	for p in $(shell ls $(SNAPSHOTS_DIR)); do \
		$(EXEC) $(SNAPSHOTS_DIR)/$$p > $$p.html; \
	done
