package version

import (
	"regexp"
	"strconv"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
)

type (
	// Info is used to store application version information. After being
	// initialized, ParseVersion must be run to populate the information.
	Info struct {
		version            string
		numCommitsSinceTag int
		commitID           string
		isDirty            bool
		release            string
		releaseCandidate   int
		branch             string
		versionString      string
	}
)

var (
	gitVersion string
	gitBranch  string
)

func (v *Info) setVersion(version string) {
	v.version = version
}

// GetDotVersion returns a string of the application version, "0.0.0", if it
// has one (requires a tag that represents a version.)
func (v Info) GetDotVersion() string {
	return v.version
}

func (v *Info) setcommitID(commitID string) {
	v.commitID = commitID
}

// GetGitCommit returns a string for the commit id where this build came from.
func (v Info) GetGitCommit() string {
	return v.commitID
}

func (v *Info) setNumCommitsSinceTag(numCommits int) {
	v.numCommitsSinceTag = numCommits
}

// GetNumCommitsSinceTag will return an int greater than 0 if there have been any
// commits since the last tag.
func (v Info) GetNumCommitsSinceTag() int {
	return v.numCommitsSinceTag
}

func (v *Info) setIsDirty(dirty bool) {
	v.isDirty = dirty
}

// GetIsDirty will return true if the git repo had local modification/uncommited
// changes, otherwise false.
func (v Info) GetIsDirty() bool {
	return v.isDirty
}

func (v *Info) setRelease(release string) {
	v.release = release
}

// GetRelease will return a string of either "Dev", "Release", "RC" depending on
// how the application was built.
func (v Info) GetRelease() string {
	return v.release
}

func (v *Info) setReleaseCandidate(rc int) {
	v.releaseCandidate = rc
}

// GetReleaseCandidate will return an int greater than 0 if this build is a RC
// where the value is which RC it is.
func (v Info) GetReleaseCandidate() int {
	return v.releaseCandidate
}

func (v *Info) setBranch(branch string) {
	v.branch = branch
}

// GetBranch returns the string of the branch this application was built from.
func (v Info) GetBranch() string {
	return v.branch
}

func (v *Info) setDefaultInfo() {
	v.setVersion("0.0.0")
	v.setcommitID("null")
	v.setNumCommitsSinceTag(0)
	v.setIsDirty(true)
	v.setRelease("Dev")
	v.setBranch(gitBranch)
}

// ParseVersion needs to be run to gather the version info and populate the struct.
func (v *Info) ParseVersion() {
	v.generateVersion()
	v.setVersionString()
}

func (v *Info) generateVersion() {
	// Git info not supplied in build env, set default development version
	if gitVersion == "" {
		v.setDefaultInfo()
		return
	}
	versionRE := regexp.MustCompile(`^(?:v(?P<version>\d\.\d\.\d)(?:rc(?P<rc>\d+))?)?(?:-(?P<numcommits>[0-9]*))?(?:(?:-?g)?(?P<commit>[0-9a-f]{6,40}))?(?P<dirty>-dirty)?$`)
	versionGroups := utils.REGetNamedGroupsResults(versionRE, gitVersion)
	// supplied git info didnt contain info needed
	if len(versionGroups) == 0 {
		v.setDefaultInfo()
		return
	}

	v.setBranch(gitBranch)

	if foundVersion, ok := versionGroups["version"]; ok {
		if foundVersion != "" {
			v.setVersion(foundVersion)
		} else {
			v.setVersion("0.0.0")
		}
	}

	if foundCommitID, ok := versionGroups["commit"]; ok {
		if foundCommitID != "" {
			v.setcommitID(foundCommitID)
		} else {
			v.setcommitID("")
		}
	}

	if foundNumCommits, ok := versionGroups["numcommits"]; ok {
		if foundNumCommits != "" {
			i, err := strconv.Atoi(foundNumCommits)
			if err != nil {
				log.Warn("Failed to parse number of git commits in version.")
				v.setNumCommitsSinceTag(0)
			} else {
				v.setNumCommitsSinceTag(i)
			}
		} else {
			v.setNumCommitsSinceTag(-1)
		}
	}

	if foundDirty, ok := versionGroups["dirty"]; ok {
		if foundDirty != "" {
			v.setIsDirty(true)
		} else {
			v.setIsDirty(false)
		}
	}

	if foundReleaseCandidate, ok := versionGroups["rc"]; ok {
		if foundReleaseCandidate != "" {
			i, err := strconv.Atoi(foundReleaseCandidate)
			if err != nil {
				log.Warn("Failed to parse RC in git version info: " + foundReleaseCandidate)
				v.setReleaseCandidate(0)
			} else {
				v.setReleaseCandidate(i)
			}
		} else {
			v.setReleaseCandidate(-1)
		}
	}

	v.setRelease("Dev")
	if v.GetNumCommitsSinceTag() < 0 && !v.GetIsDirty() && v.GetGitCommit() == "" {
		if v.GetBranch() == "master" {
			if v.GetReleaseCandidate() < 0 {
				v.setRelease("Release")
			} else if v.GetReleaseCandidate() > 0 {
				v.setRelease("RC")
			}
		}
	}
}

func (v *Info) setVersionString() {
	var vl []string
	if v.GetRelease() == "RC" {
		vl = []string{v.GetRelease(), strconv.Itoa(v.GetReleaseCandidate()), " ", v.GetDotVersion()}
	} else {
		vl = []string{v.GetRelease(), " ", v.GetDotVersion()}
	}
	if v.GetRelease() == "Dev" {
		if v.GetGitCommit() != "" {
			vl = append(vl, "-", v.GetGitCommit())
		}
		if v.GetIsDirty() {
			vl = append(vl, "+CHANGES")
		}
		if v.GetBranch() != "master" {
			vl = append(vl, "-", v.GetBranch())
		}
	}
	v.versionString = utils.StrConcat(vl)
}

// GetVersionString returns a formatted string of the Info struct for the
// version information in the application.
func (v Info) GetVersionString() string {
	return v.versionString
}
