package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//
// Minimal struct for holding our parsed command arguments
//
type Command struct {
	Args []string
}

func (c *Command) Run() error {
	cmd := exec.Command(c.Args[0], c.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func HookHandler(w http.ResponseWriter, r *http.Request) {
	h, p, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Errorf("Failed: %s", err)
		return
	}

	fmt.Printf("Received request from %s:%s...\n", h, p)

	for _, cmd := range Commands {
		fmt.Printf("Executing %s...\n", cmd.Args)
		err := cmd.Run()
		if err != nil {
			fmt.Errorf("Failed: %s", err)
			return
		}
	}
}

func ParseArgs(args []string, runnerTmpl string) ([]*Command, error) {
	cmds := []*Command{}

	//loop each argument 'echo done' 'docker run -p -u' etc
	for _, arg := range args {
		cmdparts := []string{}

		//so we compile each part of the runner seperately so
		//argument seperation is maintained correctly when input is complex
		for _, p := range strings.Split(runnerTmpl, " ") {

			rt, err := template.New("runner").Parse(p)
			if err != nil {
				return cmds, fmt.Errorf("Error while parsing runner template (%s), part: (%s): %s", runnerTmpl, p, err)
			}

			b := bytes.NewBuffer(nil)
			err = rt.Execute(b, arg)
			if err != nil {
				return cmds, fmt.Errorf("Error while executing runner template (%s, part: %s): %s", runnerTmpl, p, err)
			}

			cmdparts = append(cmdparts, b.String())
		}

		cmds = append(cmds, &Command{cmdparts})
	}

	return cmds, nil
}

//configuration options
var bindAddr = flag.String("bind", ":30000", "Address on which the updater will listen for incoming webhooks.")
var runner = flag.String("runner", "sh -c {{.}}", "Shell wrapper that runs each command.")

// global commands list
var Commands = []*Command{}

func main() {
	var err error

	flag.Parse()

	Commands, err = ParseArgs(flag.Args(), *runner)
	if err != nil {
		log.Print(err)
	}

	fmt.Printf("\nListening on '%s', when requests arrive the following %d command(s) will run: \n\n", *bindAddr, len(flag.Args()))
	for i, cmd := range Commands {
		fmt.Printf("\t%d. %s \n", i, cmd.Args)
	}

	http.HandleFunc("/", HookHandler)
	http.ListenAndServe(*bindAddr, nil)
}
