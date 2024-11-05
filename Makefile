SNAPSHOT_PB:=./snapshot_pb/objects/bafyreif2t3jzrbcn6gvm37n3b3vcslwcgrnxmamehriu265oboprfgrvbe.pb
EXEC:=./bin/anytype-publish-renderer
render:
	templ generate -lazy
	go build -o $(EXEC) . && $(EXEC) $(SNAPSHOT_PB)  > index.html
