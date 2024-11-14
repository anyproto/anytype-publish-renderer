SNAPSHOT_DIR:=./snapshot_pb/
ROOT_ID:=bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui
EXEC:=./bin/anytype-publish-renderer
render:
	templ generate -lazy
	go build -o $(EXEC) . && $(EXEC) $(SNAPSHOT_DIR) $(ROOT_ID) > index.html

test:
	go test -v ./...
