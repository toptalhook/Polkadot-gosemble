package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/aura"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/constants/timestamp"
	am "github.com/LimeChain/gosemble/frame/aura/module"
	bm "github.com/LimeChain/gosemble/frame/balances/module"
	sm "github.com/LimeChain/gosemble/frame/system/module"
	tsm "github.com/LimeChain/gosemble/frame/timestamp/module"
	"github.com/LimeChain/gosemble/primitives/support"
)

var Modules = map[sc.U8]support.ModuleMetadata{
	system.ModuleIndex:    sm.NewSystemModule(),
	aura.ModuleIndex:      am.NewAuraModule(),
	timestamp.ModuleIndex: tsm.NewTimestampModule(),
	balances.ModuleIndex:  bm.NewBalancesModule(),
}
