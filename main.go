package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	semver.Sort(releases)

	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	for _, r := range releases {
		if !r.LessThan(*minVersion) {
			versionSlice = append(versionSlice, semver.New(r.String()))
		}
	}

	var results []*semver.Version
	var currentVersion = versionSlice[0]
	var currentMajor = currentVersion.Major
	var currentMinor = currentVersion.Minor

	for _, v := range versionSlice {
		if v.Major != currentMajor || v.Minor != currentMinor {
			currentMajor = v.Major
			currentMinor = v.Minor
			results = append(results, currentVersion)
		}
		currentVersion = v
	}
	results = append(results, currentVersion)

	reverse(results)

	return results
}

func reverse(a []*semver.Version) {
	for i := len(a)/2-1; i >= 0; i-- {
		opp := len(a)-1-i
		a[i], a[opp] = a[opp], a[i]
	}
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func getReleases(repository string, minVersion *semver.Version) []*semver.Version {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	repo := strings.Split(repository, "/")
	releases, _, err := client.Repositories.ListReleases(ctx, repo[0], repo[1], opt)
	if err != nil {
		// return nil // is this really a good way?
		fmt.Println(err.Error())
	}

	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	return allReleases
}
func main() {
	// Github
	minVersion := semver.New("1.8.0")
	allReleases := getReleases("kubernetes/kubernetes", minVersion)
	versionSliceK := LatestVersions(allReleases, minVersion)

	minVersion = semver.New("2.2.0")
	allReleases = getReleases("prometheus/prometheus", minVersion)
	versionSliceP := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of kubernetes/kubernetes: %s", versionSliceK)
	fmt.Printf("\nlatest versions of prometheus/prometheus: %s", versionSliceP)

}