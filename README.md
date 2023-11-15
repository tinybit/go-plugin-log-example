# Modified KV Example

This is an improved example of usage for [hashicorp/plugin](https://github.com/hashicorp/go-plugin)

Modifications include:
- injection of [zerolog](https://github.com/rs/zerolog) logger into plugin mechanism.
- removal of unnesessary net/rpc
- capturing native stderr logs from Plugins in Main process
- proper logger calls from Plugins into main Process using bidirectional communication

This example builds a simple key/value store CLI where the mechanism for storing and retrieving keys is pluggable.

To build this example:
```sh
$ make
```

To run this example:
```sh
$ make run_put
$ make run_get
```
