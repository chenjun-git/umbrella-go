package monitor

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"umbrella-go/umbrella-common/json"
)

type Version struct {
	GitTag    string `json:"git_tag"`
	GitHash   string `json:"git_hash"`
	BuildTime string `json:"build_time"`
}

var (
	GlobalVersion    Version
	versionJsonCache []byte
)

func InitVersion(v Version) {
	GlobalVersion = v
	bs, _ := json.Marshal(v)
	versionJsonCache = bs
}

func init() {
	MonitorHandlers["/internal/version"] = http.HandlerFunc(GetVersionHandler)

	if err := AddFlags(flag.CommandLine); err != nil {
		panic(err)
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "BuildTime: %s\nGitHash: %s\nGitTag: %s\n\n",
			GlobalVersion.BuildTime, GlobalVersion.GitHash, GlobalVersion.GitTag)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

// AddFlags adds the flags used by this package to the given FlagSet. That's
// useful if working with a custom FlagSet. The init function of this package
// adds the flags to flag.CommandLine anyway. Thus, it's usually enough to call
// flag.Parse() to make the logging flags take effect.
func AddFlags(fs *flag.FlagSet) error {
	fs.Var(
		newVersionVar(false),
		"v",
		"print version",
	)

	return nil
}

type version bool

func newVersionVar(value bool) *version {
	ver := version(value)
	return &ver
}

func (ver *version) String() string {
	return fmt.Sprint(*ver)
}

// Set implements flag.Value.
func (ver *version) Set(value string) error {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	*ver = version(val)
	if *ver {
		fmt.Fprintf(os.Stderr, "BuildTime: %s\nGitHash: %s\nGitTag: %s\n",
			GlobalVersion.BuildTime, GlobalVersion.GitHash, GlobalVersion.GitTag)
		os.Exit(0)
	}
	return nil
}

func (ver *version) IsBoolFlag() bool {
	return true
}

func GetVersionHandler(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", strconv.Itoa(len(versionJsonCache)))
	w.Write(versionJsonCache)
}
