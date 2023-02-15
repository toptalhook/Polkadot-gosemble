package inherent

import (
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/types"
)

func CheckInherents(data types.InherentData, block types.Block) types.CheckInherentsResult {
	result := types.NewCheckInherentsResult()

	for _, extrinsic := range block.Extrinsics {
		// Inherents are before any other extrinsics.
		// And signed extrinsics are not inherents.
		if extrinsic.IsSigned() {
			break
		}

		isInherent := false

		call := extrinsic.Function

		switch call.CallIndex.ModuleIndex {
		case timestamp.ModuleIndex:
			if call.CallIndex.FunctionIndex == timestamp.FunctionIndex {
				isInherent = true
				err := timestamp.CheckInherent(call, data)
				if err != nil {
					err := result.PutError(timestamp.InherentIdentifier, err)
					if err != nil {
						panic(err)
					}

					if result.FatalError {
						return result
					}
				}
			}
		}

		// Inherents are before any other extrinsics.
		// No module marked it as inherent thus it is not.
		if !isInherent {
			break
		}
	}

	// TODO: go through all required pallets with required inherents

	return result
}