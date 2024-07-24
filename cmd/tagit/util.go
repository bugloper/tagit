package tagit

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func GetStagingPrefix(env string) string {
	if match, _ := regexp.MatchString(`^s(\d+)$`, env); match {
		return env + "v"
	}
	panic("Invalid environment format for staging")
}

func LatestTag(prefix string) string {
	out, _ := exec.Command("git", "tag", "-l", prefix+"*").Output()
	tags := strings.Split(string(out), "\n")
	sort.Strings(tags)
	if len(tags) > 0 {
		return tags[len(tags)-1]
	}
	return ""
}

func IncrementVersion(version string, tagType string) string {
	parts := strings.Split(version, ".")
	major, minor, patch := toInt(parts[0]), toInt(parts[1]), toInt(parts[2])

	// Semantic Versioning
	switch tagType {
	case "x": // major release resets both minor and patch to 0
		major++
		minor, patch = 0, 0
	case "y": // major release resets patch to 0
		minor++
		patch = 0
	case "z":
		patch++
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
