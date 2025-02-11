# anytype-publish-renderer

- `make fmt` formats the code
- `make lints` does linter checks

this also wired to commit hooks

## to test with local publish:

1. set the env var to e.g. renderer/test_snapshots/test-me
`export ANYTYPE_LOCAL_PUBLISH_DIR=/home/zarkone/anytype/anytype-publish-renderer/test_snapshots/test-me`
Run the client with this env and publish a page. Make sure `test-me` was created with `index.json.gz` inside.

3. in renderer, adjust Makefile SNAPSHOT_PATH accordingly:
 `SNAPSHOT_PATH:=./test_snapshots/test-me`

4. `make render`
it will create `index.html` in renderer root.
(make sure your `~/go/bin` is in `$PATH`)

5. launch a static web server in renderer root:
```
cd /home/zarkone/anytype/anytype-publish-renderer/
python -m http.server 8011
```

6. Open http://localhost:8011 in browser

## to enable css debug:
```
export ANYTYPE_PUBLISH_CSS_DEBUG=y
```
