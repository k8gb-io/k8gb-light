package mapper

/*
Copyright 2022 The k8gb Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/

import (
	"sort"
	"strings"

	"github.com/k8gb-io/k8gb-light/controllers/utils"
)

type PrimaryGeotag []string

// NewPrimaryGeotag returns regions sorted by priority as they will be used in failover.
// It takes geotags as they are defined in annotation followed by the rest of ExtGEoTags sorted alphabetically
// for primaryGeoTag "us, eu" with extClusterGeotags = []{"za","cz","us","uk","eu"} returns
// []{"us","eu","cz","uk","za"}
func NewPrimaryGeotag(primaryGeoTagList string, extClusterGeoTags []string) PrimaryGeotag {
	sort := func(tags []string) []string {
		sort.Slice(tags, func(i, j int) bool {
			return tags[i] < tags[j]
		})
		return tags
	}
	var pg []string
	existsInPrimaryGeoTagList := utils.AsMap(extClusterGeoTags)
	extClusterGeoTagsSorted := sort(extClusterGeoTags)

	if primaryGeoTagList != "" {
		for _, v := range strings.Split(primaryGeoTagList, ",") {
			v = strings.TrimLeft(strings.TrimRight(v, " "), " ")
			pg = append(pg, v)
			existsInPrimaryGeoTagList[v] = false
		}
	}
	for _, v := range extClusterGeoTagsSorted {
		if existsInPrimaryGeoTagList[v] {
			pg = append(pg, v)
		}
	}

	return pg
}

func (pg PrimaryGeotag) String() string {
	return strings.Join(pg, ",")
}
