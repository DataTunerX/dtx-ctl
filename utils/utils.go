package utils

import "github.com/blang/semver/v4"

func ParseDataTunerXVersion(version string) (semver.Version, error) {
	return semver.ParseTolerant(version)
}
