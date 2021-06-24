package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type arguments struct {
	commit  string
	uname   string
	time    string
	version string
}

type platformInfo struct {
	goos   string
	goarch string
}

func (p platformInfo) String() string {
	if p.goos == "darwin" {
		return "macos-" + p.goarch
	}
	return p.goos + "-" + p.goarch
}

func main() {
	log.SetFlags(0)

	var args arguments
	flag.StringVar(&args.commit, "build-commit", "", "build commit hash")
	flag.StringVar(&args.uname, "build-uname", "", "build uname information")
	flag.StringVar(&args.time, "build-time", "", "build time information string")
	flag.StringVar(&args.version, "version", "", "version information string")
	flag.Parse()

	platforms := []platformInfo{
		{"linux", "amd64"},
		{"linux", "arm64"},
		{"darwin", "amd64"},
		{"windows", "amd64"},
	}

	for _, platform := range platforms {
		if err := prepareArchive(args, platform); err != nil {
			log.Printf("error: build %s: %v", platform, err)
		}
	}
}

func prepareArchive(args arguments, platform platformInfo) error {
	log.Printf("building %s", platform)
	binaryExt := ""
	if platform.goos == "windows" {
		binaryExt = ".exe"
	}

	releaseDir := ".release-v" + args.version
	pack := "github.com/vkcom/nocolor/cmd"

	ldFlags := fmt.Sprintf(`-X '%[1]s.Version=%[2]s' -X '%[1]s.BuildTime=%[3]s' -X '%[1]s.BuildOSUname=%[4]s' -X '%[1]s.BuildCommit=%[5]s'`,
		pack, args.version, args.time, args.uname, args.commit)
	binaryName := filepath.Join(releaseDir, "nocolor"+binaryExt)
	buildCmd := exec.Command("go", "build",
		"-o", binaryName,
		"-ldflags", ldFlags,
		".")
	buildCmd.Env = append([]string{}, os.Environ()...) // Copy env slice
	buildCmd.Env = append(buildCmd.Env, "GOOS="+platform.goos)
	buildCmd.Env = append(buildCmd.Env, "GOARCH="+platform.goarch)
	out, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("run %s: %v: %s", buildCmd, err, out)
	}

	// Pack it into the archive.
	archiveName := "nocolor-v" + args.version + "-" + platform.String() + ".zip"
	zipCmd := exec.Command("zip", archiveName, "nocolor"+binaryExt)
	zipCmd.Dir = releaseDir
	log.Printf("creating %s archive", archiveName)
	if out, err := zipCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("make archive: %v: %s", err, out)
	}

	return nil
}
