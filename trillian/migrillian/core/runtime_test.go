// Copyright 2018 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"strings"
	"testing"

	"github.com/google/certificate-transparency-go/trillian/migrillian/election"
)

func TestRuntime(t *testing.T) {
	noop := election.NoopFactory{}
	for _, tc := range []struct {
		desc    string
		ef      election.Factory
		opts    RuntimeOptions
		wantErr string
	}{
		{desc: "empty-election-factory", ef: nil, wantErr: "empty election factory"},
		{
			desc:    "wrong-fetchers",
			ef:      noop,
			opts:    RuntimeOptions{Fetchers: 0},
			wantErr: "fetchers Pool",
		},
		{
			desc:    "wrong-submitters",
			ef:      noop,
			opts:    RuntimeOptions{Fetchers: 1, Submitters: -1},
			wantErr: "submitters Pool",
		},
		{
			desc:    "wrong-buf",
			ef:      noop,
			opts:    RuntimeOptions{Fetchers: 3, Submitters: 2, SubmittersBuf: -10},
			wantErr: "submitters Pool",
		},
		{desc: "ok-1", ef: noop, opts: RuntimeOptions{Fetchers: 3, Submitters: 2}},
		{desc: "ok-2", ef: noop, opts: RuntimeOptions{Fetchers: 3, Submitters: 2, SubmittersBuf: 10}},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			rt, err := NewRuntime(tc.ef, tc.opts)
			if len(tc.wantErr) > 0 && (err == nil || !strings.Contains(err.Error(), tc.wantErr)) {
				t.Errorf("New(): err=%v, want err containing %v", err, tc.wantErr)
			} else if len(tc.wantErr) == 0 && err != nil {
				t.Errorf("New(): err=%v, want nil", err)
			}
			if err != nil {
				return
			}
			rt.Start()
			rt.Stop()
		})
	}
}
