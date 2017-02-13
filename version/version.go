package version

import (
	"regexp"
	"strconv"

	"github.com/shawnsilva/piccolo/log"
	"github.com/shawnsilva/piccolo/utils"
)

type (
	VersionInfo struct {
		version            string
		numCommitsSinceTag int
		commitId           string
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

func (v *VersionInfo) setVersion(version string) {
	v.version = version
}

func (v VersionInfo) GetDotVersion() string {
	return v.version
}

func (v *VersionInfo) setCommitId(commitId string) {
	v.commitId = commitId
}

func (v VersionInfo) GetGitCommit() string {
	return v.commitId
}

func (v *VersionInfo) setNumCommitsSinceTag(numCommits int) {
	v.numCommitsSinceTag = numCommits
}

func (v VersionInfo) GetNumCommitsSinceTag() int {
	return v.numCommitsSinceTag
}

func (v *VersionInfo) setIsDirty(dirty bool) {
	v.isDirty = dirty
}

func (v VersionInfo) GetIsDirty() bool {
	return v.isDirty
}

func (v *VersionInfo) setRelease(release string) {
	v.release = release
}

func (v VersionInfo) GetRelease() string {
	return v.release
}

func (v *VersionInfo) setReleaseCandidate(rc int) {
	v.releaseCandidate = rc
}

func (v VersionInfo) GetReleaseCandidate() int {
	return v.releaseCandidate
}

func (v *VersionInfo) setBranch(branch string) {
	v.branch = branch
}

func (v VersionInfo) GetBranch() string {
	return v.branch
}

func (v *VersionInfo) setDefaultInfo() {
	v.setVersion("0.0.0")
	v.setCommitId("null")
	v.setNumCommitsSinceTag(0)
	v.setIsDirty(true)
	v.setRelease("Dev")
	v.setBranch(gitBranch)
}

func (v *VersionInfo) ParseVersion() {
	v.generateVersion()
}

func (v *VersionInfo) generateVersion() {
	// Git info not supplied in build env, set default development version
	if gitVersion == "" {
		v.setDefaultInfo()
		return
	}
	versionRE := regexp.MustCompile(`^(?:v(?P<version>\d\.\d\.\d)(?P<rc>rc\d+)?)?(?:-?(?P<numcommits>[0-9]*))?(?:(?:-?g)?(?P<commit>[0-9a-f]{6,40}))?(?P<dirty>-dirty)?$`)
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

	if foundCommitId, ok := versionGroups["commit"]; ok {
		if foundCommitId != "" {
			v.setCommitId(foundCommitId)
		} else {
			v.setCommitId("")
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
				log.Warn("Failed to parse RC in git version info.")
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
	v.setVersionString()
}

func (v *VersionInfo) setVersionString() {
	vl := []string{v.GetRelease(), " ", v.GetDotVersion()}
	if v.GetRelease() == "RC" {
		vl = append(vl, string(v.GetReleaseCandidate()))
	}
	if v.GetRelease() == "Dev" {
		vl = append(vl, "-", v.GetGitCommit())
		if v.GetIsDirty() {
			vl = append(vl, "+CHANGES")
		}
		if v.GetBranch() != "master" {
			vl = append(vl, "-", v.GetBranch())
		}
	}
	v.versionString = utils.StrConcat(vl)
}

func (v VersionInfo) GetVersionString() string {
	return v.versionString
}
