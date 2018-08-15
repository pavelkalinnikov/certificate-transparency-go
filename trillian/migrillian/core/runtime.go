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
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/google/certificate-transparency-go/trillian/migrillian/election"
	"github.com/google/certificate-transparency-go/util/exepool"
)

// RuntimeOptions are the options used to create a Runtime.
type RuntimeOptions struct {
	Fetchers      int // The number of Fetcher workers.
	Submitters    int // The number of Submitter workers.
	SubmittersBuf int // The buffer size for the Pool of submitters.
}

// Runtime holds objects that Controllers use at run time.
type Runtime struct {
	ef  election.Factory
	fxp *exepool.Pool // Fetchers execution Pool.
	sxp *exepool.Pool // Submitters execution Pool.
}

// NewRuntime creates a Runtime.
func NewRuntime(ef election.Factory, opts RuntimeOptions) (*Runtime, error) {
	if ef == nil {
		return nil, errors.New("empty election factory")
	}
	fxp, err := exepool.New(opts.Fetchers, 0)
	if err != nil {
		return nil, fmt.Errorf("fetchers Pool.New(): %v", err)
	}
	sxp, err := exepool.New(opts.Submitters, opts.SubmittersBuf)
	if err != nil {
		return nil, fmt.Errorf("submitters Pool.New(): %v", err)
	}
	return &Runtime{ef: ef, fxp: fxp, sxp: sxp}, nil
}

// Start starts the Runtime.
func (rt *Runtime) Start() {
	rt.fxp.Start()
	rt.sxp.Start()
}

// Stop stops the Runtime.
func (rt *Runtime) Stop() {
	if err := rt.fxp.Stop(); err != nil {
		glog.Errorf("Runtime: Fetchers Pool.Stop() failed: %v", err)
	}
	if err := rt.sxp.Stop(); err != nil {
		glog.Errorf("Runtime: Submitters Pool.Stop() failed: %v", err)
	}
}
