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

package utils

import "reflect"

// Contains helps to detect if a non-comparable struct is part of an array
// only use this method if you can't rely on existing golang Contains function of slices (https://pkg.go.dev/golang.org/x/exp/slices#Contains)
func Contains[T any](origin []T, element T) bool {
	for _, v := range origin {
		if reflect.DeepEqual(v, element) {
			return true
		}
	}
	return false
}

// Remove removes all elements from origin slice
func Remove[T any](origin []T, elements ...T) []T {
	var filtered []T
	for _, v := range origin {
		if !Contains(elements, v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func Filter[T any](elements []T, predicate func(T) bool) []T {
	var filtered []T
	for _, v := range elements {
		if predicate(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
