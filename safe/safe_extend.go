package safe

import (
	"github.com/injoyai/base/chans"
	"github.com/injoyai/base/maps"
)

type (
	Map  = maps.Safe
	Chan = chans.Safe
)

var (
	NewMap  = maps.NewSafe
	NewChan = chans.NewSafe
)
