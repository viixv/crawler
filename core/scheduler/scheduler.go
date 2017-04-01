// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
package scheduler

import (
	"github.com/viixv/crawler/core/commons/request"
)

type Scheduler interface {
	Push(req *request.Request)
	Poll() *request.Request
	Count() int
}
