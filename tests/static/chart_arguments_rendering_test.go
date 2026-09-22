/*
 * Copyright (c) 2026, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package static_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/yaml"
)

func TestChartArgumentsRenderContract(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
	}{
		{name: "numeric collection interval", args: []string{"-c", "1000"}},
		{name: "YAML scalars and empty strings", args: []string{"true", "null", "", "0123"}},
		{name: "YAML punctuation", args: []string{"path: value", "#comment", "[value]", "quote\"and\\slash"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			release, err := runHelmInstall(t, map[string]interface{}{"arguments": tc.args})
			require.NoError(t, err)
			data, err := yaml.YAMLToJSON(requireManifestKind(t, release.Manifest, "DaemonSet"))
			require.NoError(t, err)
			var daemonSet appsv1.DaemonSet
			require.NoError(t, json.Unmarshal(data, &daemonSet))
			require.Equal(t, tc.args, daemonSet.Spec.Template.Spec.Containers[0].Args)
		})
	}
}
