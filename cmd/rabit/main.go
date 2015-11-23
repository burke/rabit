package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/flynn/go-docopt"
)

const usageTpl = `usage: %s [-h|--help] <command> [<args>...]

Environment Variables:
  RABIT_DIR     Path on disk to the rabit repository
  RABIT_REMOTE  URL of remote rabit repository

Options:
  -h, --help

Commands:
  help       Show usage for a specific command
  init       Initialize a new rabit repository
  add        Add a file to the rabit repository
  ls-files   List files in a rabit repository
  cat-file   Print the contents of a file in the repository
  rm         Remove a file from the rabit repository
  gc         Remove any blocks belonging only to removed manifests
  push       Upload to the rabit server
  fetch      Download from the rabit server
  ls-remote  List files available for download from the rabit server

See "%s help <command>" for more information on a specific command.
`

func main() {
	log.SetFlags(0)

	usage := fmt.Sprintf(usageTpl, os.Args[0], os.Args[0])
	args, _ := docopt.Parse(usage, nil, true, "", true)

	cmd := args.String["<command>"]
	cmdArgs := args.All["<args>"].([]string)

	if cmd == "help" {
		if len(cmdArgs) == 0 { // `rabit help`
			fmt.Println(usage)
			return
		} else { // `rabit help <command>`
			cmd = cmdArgs[0]
			cmdArgs = []string{"--help"}
		}
	}

	if err := runCommand(cmd, cmdArgs); err != nil {
		log.Fatalln("ERROR:", err)
	}
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "usage: %s <command>\n", os.Args[0])
}

type cmdFunc func(*docopt.Args, string, string) error

type command struct {
	usage        string
	f            cmdFunc
	verifyDir    bool
	verifyRemote bool
}

var commands = make(map[string]*command)

func register(name string, f cmdFunc, verifyDir, verifyRemote bool, usage string) {
	commands[name] = &command{
		usage:        fmt.Sprintf(usage, os.Args[0]),
		f:            f,
		verifyDir:    verifyDir,
		verifyRemote: verifyRemote,
	}
}

func runCommand(name string, args []string) error {
	argv := make([]string, 1, 1+len(args))
	argv[0] = name
	argv = append(argv, args...)

	cmd, ok := commands[name]
	if !ok {
		return fmt.Errorf("%s is not a rabit command. See 'rabit help'", name)
	}

	parsedArgs, err := docopt.Parse(cmd.usage, argv, true, "", true)
	if err != nil {
		return err
	}

	rabitDir := os.Getenv("RABIT_DIR")
	rabitRemote := os.Getenv("RABIT_REMOTE")

	if cmd.verifyDir {
		if rabitDir == "" {
			return fmt.Errorf("RABIT_DIR must specify a path to a valid rabit repository")
		}
		stat, err := os.Stat(rabitDir)
		if err != nil || !stat.IsDir() {
			return fmt.Errorf("RABIT_DIR must specify a path to a valid rabit repository")
		}
	}

	if cmd.verifyRemote {
		if rabitRemote == "" {
			return fmt.Errorf("RABIT_REMOTE must be set")
		}
	}

	return cmd.f(parsedArgs, rabitDir, rabitRemote)
}
