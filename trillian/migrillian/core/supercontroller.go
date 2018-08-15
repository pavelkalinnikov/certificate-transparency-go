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
	"context"

	"github.com/golang/glog"
	"github.com/google/trillian/monitoring"
)

// SuperController migrates data from any number of CT logs to Trillian.
//
// TODO(pavelkalinnikov): Rename Controller->Migration/Migriller, and
// SuperController->Controller.
type SuperController struct {
	ctrls []*Controller
	mf    monitoring.MetricFactory
}

// NewSuperController creates a SuperController.
func NewSuperController(mf monitoring.MetricFactory) *SuperController {
	return &SuperController{mf: mf}
}

// Add adds the specified migration Controller to the list of migrations
// managed by this SuperController. It must be invoked before Run.
func (sc *SuperController) Add(ctrl *Controller) {
	sc.ctrls = append(sc.ctrls, ctrl)
}

// Run executes log migration for all logs that this SuperController manages.
func (sc *SuperController) Run(ctx context.Context, rt *Runtime) error {
	rt.Start()
	defer rt.Stop()

	for _, ctrl := range sc.ctrls {
		// TODO(pavelkalinnikov): Add running status of each goroutine to
		// monitoring.
		// TODO(pavelkalinnikov): Restart on failure, or maybe Exit. Could also
		// cancel other Controllers for a more graceful termination.
		go func() {
			if err := ctrl.RunWhenMaster(ctx, rt); err != nil {
				glog.Exitf("Controller.RunWhenMaster() returned: %v", err)
			}
		}()
	}
	return nil
}
