package version

import (
	"testing"
)

type (
	versionParsingTest struct {
		gitVersion string
		gitBranch  string

		versionResult            string
		numCommitsSinceTagResult int
		commitIDResult           string
		isDirtyResult            bool
		releaseResult            string
		releaseCandidateResult   int
		branchResult             string
		versionStringResult      string
	}
)

var (
	versionParsingTestData = []versionParsingTest{
		{"3aab5a6-dirty", "master", "0.0.0", 0, "3aab5a6", true, "Dev", 0, "master", "Dev 0.0.0-3aab5a6+CHANGES"},
		{"3aab5a6", "something", "0.0.0", 0, "3aab5a6", false, "Dev", 0, "something", "Dev 0.0.0-3aab5a6-something"},
		{"3aab5a6-dirty", "something", "0.0.0", 0, "3aab5a6", true, "Dev", 0, "something", "Dev 0.0.0-3aab5a6+CHANGES-something"},
		{"", "master", "0.0.0", 0, "null", true, "Dev", 0, "master", "Dev 0.0.0-null+CHANGES"},
		{"v0.1.0-3-g3aab5a6-dirty", "master", "0.1.0", 3, "3aab5a6", true, "Dev", 0, "master", "Dev 0.1.0-3aab5a6+CHANGES"},
		{"v0.1.0", "master", "0.1.0", 0, "", false, "Release", 0, "master", "Release 0.1.0"},
		{"v0.1.0-dirty", "master", "0.1.0", 0, "", true, "Dev", 0, "master", "Dev 0.1.0+CHANGES"},
		{"v0.1.0", "test", "0.1.0", 0, "", false, "Dev", 0, "test", "Dev 0.1.0-test"},
		{"v0.1.0-dirty", "test", "0.1.0", 0, "", false, "Dev", 0, "test", "Dev 0.1.0+CHANGES-test"},
		{"v0.1.0rc3", "master", "0.1.0", 0, "", false, "RC", 3, "master", "RC3 0.1.0"},
		{"v0.1.0rc3-dirty", "master", "0.1.0", 0, "", true, "RC", 3, "master", "Dev 0.1.0+CHANGES"},
	}
)

func TestParseVersion(t *testing.T) {
	for _, testData := range versionParsingTestData {
		gitVersion = testData.gitVersion
		gitBranch = testData.gitBranch
		versionTest := Info{}
		versionTest.ParseVersion()
		if versionTest.GetVersionString() != testData.versionStringResult {
			t.Error(
				"Git Version: ", testData.gitVersion,
				"Git Branch: ", testData.gitBranch,
				"Got: ", versionTest.GetVersionString(),
			)
		}
	}
}
