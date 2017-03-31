// Copyright 2014 Hu Cong. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package processor

import (
	"github.com/viixv/crawler/core/commons/page"
)

type PageProcessor interface {
	Process(p *page.Page)
	Finish()
}
