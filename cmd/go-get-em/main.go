package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Version = "dev"

func main() {
	var (
		review = true
		update = true
	)
	log.SetFlags(log.Lshortfile)
	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", filepath.Base(os.Args[0]), Version), flag.ExitOnError)
	flags.BoolVar(&review, "review", review, "When set, print 'open' commands with vcs diff URLs for outdated dependencies.")
	flags.BoolVar(&update, "update", update, "When set, print 'go get' commands to upgrade outdated dependencies.")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage of %s:\n", flags.Name())
		_, _ = fmt.Fprintf(flags.Output(), "  "+
			"Run the command in a directory w/ a go.mod file to emit a list of "+
			"outdated dependencies and commands to review and update them.",
		)
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])

	var output bytes.Buffer
	command := exec.Command("go", "list", "-m", "-u", "-json", "all")
	command.Stdout = &output
	command.Stderr = os.Stderr
	log.Printf("[INFO] Sit tight, executing command: %s", command.String())
	err := command.Run()
	if err != nil {
		log.Fatalf("[WARN] Failed to execute command [%s]: %v", command.String(), err)
	}

	var dependencies []Module
	for decoder := json.NewDecoder(&output); ; {
		var dependency Module
		err = decoder.Decode(&dependency)
		if dependency.Path != "" {
			dependencies = append(dependencies, dependency)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("[WARN] Failed to decode dependencies: %v", err)
		}
	}

	if len(dependencies) == 0 {
		log.Println("[INFO] No dependencies found.")
		return
	}

	var outdated []Module
	for _, dependency := range dependencies {
		if dependency.Deprecated != "" {
			log.Printf("[WARN] [%s] is deprecated: %s", dependency.Path, dependency.Deprecated)
		}
		if dependency.Update != nil && !dependency.Indirect {
			log.Printf("[INFO] %10s..%-10s  %s", dependency.Version, dependency.Update.Version, dependency.Path)
			outdated = append(outdated, dependency)
		}
	}

	if len(outdated) == 0 {
		log.Println("[INFO] No outdated dependencies found.")
		return
	}

	if !review && !update {
		log.Println("[INFO] No output requested")
		return
	}
	log.Println("[INFO] Execute what you will of the following output to review and/or update the outdated dependencies.")

	if review {
		fmt.Println()
		for _, dependency := range outdated {
			if strings.HasPrefix(dependency.Path, "bitbucket.org") {
				fmt.Println("open", bitbucketDiffURL(dependency))
			} else if strings.HasPrefix(dependency.Path, "github.com") {
				fmt.Println("open", githubDiffURL(dependency))
			} else {
				fmt.Println(unknownDiffURL(dependency))
			}
		}
	}

	if update {
		fmt.Println()
		for _, dependency := range outdated {
			fmt.Printf("go get -u %s@%s\n", dependency.Path, dependency.Update.Version)
		}
	}
}

func bitbucketDiffURL(dependency Module) string {
	return fmt.Sprintf("https://%s/branches/compare/%s"+"%%0D"+"%s#commits",
		vcsRepoPath(dependency),
		dependency.Update.Version,
		dependency.Version,
	)
}
func githubDiffURL(dependency Module) string {
	return fmt.Sprintf("https://%s/compare/%s...%s",
		vcsRepoPath(dependency),
		dependency.Version,
		dependency.Update.Version,
	)
}
func unknownDiffURL(dependency Module) string {
	return fmt.Sprintln("# Not sure how to render a diff URL for this module:",
		dependency.Path,
		dependency.Version,
		dependency.Update.Version,
	)
}
func vcsRepoPath(module Module) string {
	base := path.Base(module.Path)
	if !strings.HasPrefix(base, "v") {
		return module.Path
	}
	vLess := base[1:]
	n, err := strconv.Atoi(vLess)
	if err == nil && n > 1 {
		return path.Dir(module.Path)
	}
	return module.Path
}

// Module source: https://go.dev/ref/mod#go-list-m
type Module struct {
	Path       string       // module path
	Version    string       // module version
	Versions   []string     // available module versions
	Replace    *Module      // replaced by this module
	Time       *time.Time   // time version was created
	Update     *Module      // available update (with -u)
	Main       bool         // is this the main module?
	Indirect   bool         // module is only indirectly needed by main module
	Dir        string       // directory holding local copy of files, if any
	GoMod      string       // path to go.mod file describing module, if any
	GoVersion  string       // go version used in module
	Retracted  []string     // retraction information, if any (with -retracted or -u)
	Deprecated string       // deprecation message, if any (with -u)
	Error      *ModuleError // error loading module
}

// ModuleError source: https://go.dev/ref/mod#go-list-m
type ModuleError struct {
	Err string // the error itself
}
