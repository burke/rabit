# rabit

rabit is an implementation of Rabin fingerprinting for large binary blobs to
enable differential updates.

Currently, the `push`, `fetch`, and `ls-remote` commands are not implemented,
but they will enable integration with a remote [TUF
server](https://github.com/flynn/go-tuf) to manage signed differential updates.

Use the CLI as documented below or see [`the API
docs`](https://godoc.org/github.com/burke/rabit/pkg/repo).

```
usage: rabit [-h|--help] <command> [<args>...]

Environment Variables:
  RABIT_DIR     Path on disk to the rabit repository
  RABIT_REMOTE  URL of remote rabit repository

Options:
  -h, --help

Commands:
  help       Show usage for a specific command
  init       Initialize a new rabit repository
  add        Add a file to the rabit repository
  ls         List files in a rabit repository
  cat        Print the contents of a file in the repository
  rm         Remove a file from the rabit repository
  gc         Remove any blocks belonging only to removed manifests
  push       Upload to the rabit server
  fetch      Download from the rabit server
  ls-remote  List files available for download from the rabit server

See "rabit help <command>" for more information on a specific command.
```
