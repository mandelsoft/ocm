// Copyright 2022 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package mime

import (
	"strings"
)

func IsJSON(mime string) bool {
	if mime == MIME_JSON || mime == MIME_JSON_ALT {
		return true
	}
	if strings.HasSuffix(mime, "+json") {
		return true
	}
	return false
}

func IsYAML(mime string) bool {
	if mime == MIME_YAML || mime == MIME_YAML_ALT {
		return true
	}
	if strings.HasSuffix(mime, "+yaml") {
		return true
	}
	return false
}

func BaseType(mime string) string {
	i := strings.Index(mime, "+")
	if i > 0 {
		return mime[:i]
	}
	return mime
}

func IsGZip(mime string) bool {
	return strings.HasSuffix(mime, "/gzip") || strings.HasSuffix(mime, "+gzip")
}
