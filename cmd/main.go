package cmd

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/VKCOM/noverify/src/cmd"
	"github.com/VKCOM/noverify/src/linter"

	"github.com/vkcom/nocolor/internal/checkers"
	"github.com/vkcom/nocolor/internal/walkers"
)

// Build* are initialized during the build via -ldflags
var (
	BuildTime    string
	BuildOSUname string
	BuildCommit  string
)

func printVersion(mainConfig *cmd.MainConfig) {
	if BuildCommit == "" {
		log.Printf("Version %s: built without version info", mainConfig.LinterVersion)
	} else {
		log.Printf("Version %s: built on: %s OS: %s Commit: %s\n", mainConfig.LinterVersion, BuildTime, BuildOSUname, BuildCommit)
	}
}

func disableAllDefaultReports(r *linter.Report) bool {
	return checkers.Contains(r.CheckName)
}

// DefaultCacheDir function returns the default, depending on the system, directory for storing the cache.
func DefaultCacheDir() string {
	defaultCacheDir, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return filepath.Join(defaultCacheDir, "nocolor-cache")
}

// Main is the function that launches the program.
func Main() {
	config := linter.NewConfig()
	context := walkers.NewGlobalContext(nil)

	status, err := cmd.Run(&cmd.MainConfig{
		BeforeReport:             disableAllDefaultReports,
		LinterVersion:            "0.1.0",
		RegisterCheckers:         checkers.List,
		LinterConfig:             config,
		DisableCriticalIssuesLog: true,
		AfterFlagParse: func(env cmd.InitEnvironment) {
			context.Info = env.MetaInfo
		},
		ModifyApp: func(app *cmd.App) {
			app.Name = "nocolor"
			app.Description = "Architecture checking tool based on the concept of colors"

			// Clear all the default commands.
			app.Commands = nil

			app.Commands = append(app.Commands, &cmd.Command{
				Name:        "version",
				Description: "The command outputs the version",
				Action: func(ctx *cmd.AppContext) (int, error) {
					printVersion(ctx.MainConfig)
					return 0, nil
				},
			})

			app.Commands = append(app.Commands, &cmd.Command{
				Name:        "check",
				Description: "The command to start checking files",
				RegisterFlags: func(ctx *cmd.AppContext) *flag.FlagSet {
					flags := &extraCheckFlags{}

					fs := flag.NewFlagSet("check", flag.ContinueOnError)

					// We don't need all the flags from NoVerify, so we only register some of them.
					fs.StringVar(&ctx.ParsedFlags.IndexOnlyFiles, "index-only-files", "", "Comma-separated list of files to do indexing")
					fs.BoolVar(&ctx.ParsedFlags.Debug, "debug", false, "Enable debug output")
					fs.DurationVar(&ctx.ParsedFlags.DebugParseDuration, "debug-parse-duration", 0, "Print files that took longer than the specified time to analyse")
					fs.IntVar(&ctx.ParsedFlags.MaxFileSize, "max-sum-filesize", 20*1024*1024, "Max total file size to be parsed concurrently in bytes (limits max memory consumption)")
					fs.IntVar(&ctx.ParsedFlags.MaxConcurrency, "cores", runtime.NumCPU(), "Max cores")
					fs.StringVar(&ctx.ParsedFlags.StubsDir, "stubs-dir", "", "Directory with phpstorm-stubs")
					fs.StringVar(&ctx.ParsedFlags.CacheDir, "cache-dir", DefaultCacheDir(), "Directory for linter cache (greatly improves indexing speed)")
					fs.BoolVar(&ctx.ParsedFlags.DisableCache, "disable-cache", false, "If set, cache is not used and cache-dir is ignored")
					fs.StringVar(&ctx.ParsedFlags.PprofHost, "pprof", "", "HTTP pprof endpoint (e.g. localhost:8080)")
					fs.StringVar(&ctx.ParsedFlags.CPUProfile, "cpuprofile", "", "Write cpu profile to `file`")
					fs.StringVar(&ctx.ParsedFlags.MemProfile, "memprofile", "", "Write memory profile to `file`")
					fs.StringVar(&ctx.ParsedFlags.PhpExtensionsArg, "php-extensions", "php,inc,php5,phtml", "List of PHP extensions to be recognized")

					// Some values need to be set manually.
					ctx.ParsedFlags.AllowAll = true
					ctx.ParsedFlags.ReportsCritical = cmd.AllNonNoticeChecks

					fs.StringVar(&flags.PaletteSrc, "palette", "palette.yaml", "File with color palette")
					fs.StringVar(&flags.ColorTag, "tag", "color", "The tag to be used to set the color in phpdoc")

					ctx.CustomFlags = flags
					return fs
				},
				Action: func(ctx *cmd.AppContext) (int, error) {
					return Check(ctx, context)
				},
			})
		},
	})
	if err != nil {
		log.Println(err)
	}

	os.Exit(status)
}
