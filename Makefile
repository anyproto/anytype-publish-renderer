SNAPSHOT_PB:=./snapshot_pb/objects/bafyreiftcdiken5kayp3x4ix7tm4okmizhyshev3jjl5r2jjenz2d5uwui.pb
EXEC:=./bin/anytype-publish-renderer
render:
	templ generate -lazy
	go build -o $(EXEC) . && $(EXEC) $(SNAPSHOT_PB)  > index.html
