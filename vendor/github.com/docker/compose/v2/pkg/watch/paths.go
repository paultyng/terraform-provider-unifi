/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package watch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func greatestExistingAncestor(path string) (string, error) {
	if path == string(filepath.Separator) ||
		path == fmt.Sprintf("%s%s", filepath.VolumeName(path), string(filepath.Separator)) {
		return "", fmt.Errorf("cannot watch root directory")
	}

	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return "", errors.Wrapf(err, "os.Stat(%q)", path)
	}

	if os.IsNotExist(err) {
		return greatestExistingAncestor(filepath.Dir(path))
	}

	return path, nil
}

// If we're recursively watching a path, it doesn't
// make sense to watch any of its descendants.
func dedupePathsForRecursiveWatcher(paths []string) []string {
	result := []string{}
	for _, current := range paths {
		isCovered := false
		hasRemovals := false

		for i, existing := range result {
			if IsChild(existing, current) {
				// The path is already covered, so there's no need to include it
				isCovered = true
				break
			}

			if IsChild(current, existing) {
				// Mark the element empty for removal.
				result[i] = ""
				hasRemovals = true
			}
		}

		if !isCovered {
			result = append(result, current)
		}

		if hasRemovals {
			// Remove all the empties
			newResult := []string{}
			for _, r := range result {
				if r != "" {
					newResult = append(newResult, r)
				}
			}
			result = newResult
		}
	}
	return result
}

func IsChild(dir string, file string) bool {
	if dir == "" {
		return false
	}

	dir = filepath.Clean(dir)
	current := filepath.Clean(file)
	child := "."
	for {
		if strings.EqualFold(dir, current) {
			// If the two paths are exactly equal, then they must be the same.
			if dir == current {
				return true
			}

			// If the two paths are equal under case-folding, but not exactly equal,
			// then the only way to check if they're truly "equal" is to check
			// to see if we're on a case-insensitive file system.
			//
			// This is a notoriously tricky problem. See how dep solves it here:
			// https://github.com/golang/dep/blob/v0.5.4/internal/fs/fs.go#L33
			//
			// because you can mount case-sensitive filesystems onto case-insensitive
			// file-systems, and vice versa :scream:
			//
			// We want to do as much of this check as possible with strings-only
			// (to avoid a file system read and error handling), so we only
			// do this check if we have no other choice.
			dirInfo, err := os.Stat(dir)
			if err != nil {
				return false
			}

			currentInfo, err := os.Stat(current)
			if err != nil {
				return false
			}

			if !os.SameFile(dirInfo, currentInfo) {
				return false
			}
			return true
		}

		if len(current) <= len(dir) || current == "." {
			return false
		}

		cDir := filepath.Dir(current)
		cBase := filepath.Base(current)
		child = filepath.Join(cBase, child)
		current = cDir
	}
}
