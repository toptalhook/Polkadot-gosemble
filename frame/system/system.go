package system

import (
	"bytes"
	"math"

	"github.com/LimeChain/gosemble/constants/system"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/timestamp"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/storage"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/trie"
	"github.com/LimeChain/gosemble/primitives/types"
)

var Module = support.ModuleMetadata{
	Index: system.ModuleIndex,
	Functions: map[string]support.FunctionMetadata{
		"remark": {Index: system.FunctionRemarkIndex, Func: Remark},
	},
}

func Remark(args sc.Sequence[sc.U8]) {
	// TODO:
}

func Finalize() types.Header {
	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Clear(append(systemHash, executionPhaseHash...))

	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	storage.Clear(append(systemHash, allExtrinsicsLenHash...))

	numberHash := hashing.Twox128(constants.KeyNumber)
	blockNumber := storage.GetDecode(append(systemHash, numberHash...), sc.DecodeU32)

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	parentHash := storage.GetDecode(append(systemHash, parentHashKey...), types.DecodeBlake2bHash)

	digestHash := hashing.Twox128(constants.KeyDigest)
	digest := storage.GetDecode(append(systemHash, digestHash...), types.DecodeDigest)

	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)
	extrinsicCount := storage.TakeDecode(append(systemHash, extrinsicCountHash...), sc.DecodeU32)

	var extrinsics []byte
	extrinsicDataPrefixHash := append(systemHash, hashing.Twox128(constants.KeyExtrinsicData)...)

	for i := 0; i < int(extrinsicCount); i++ {
		sci := sc.U32(i)
		hashIndex := hashing.Twox64(sci.Bytes())

		extrinsicDataHashIndexHash := append(extrinsicDataPrefixHash, hashIndex...)
		extrinsic := storage.TakeBytes(append(extrinsicDataHashIndexHash, sci.Bytes()...))

		extrinsics = append(extrinsics, extrinsic...)
	}

	buf := &bytes.Buffer{}
	extrinsicsRootBytes := trie.Blake2256OrderedRoot(append(sc.ToCompact(uint64(extrinsicCount)).Bytes(), extrinsics...), constants.StorageVersion)
	buf.Write(extrinsicsRootBytes)
	extrinsicsRoot := types.DecodeH256(buf)
	buf.Reset()

	// saturating_sub
	toRemove := blockNumber - constants.BlockHashCount - 1
	if toRemove > blockNumber {
		toRemove = 0
	}

	if toRemove != 0 {
		blockNumHash := hashing.Twox64(toRemove.Bytes())
		blockNumKey := append(systemHash, hashing.Twox128(constants.KeyBlockHash)...)
		blockNumKey = append(blockNumKey, blockNumHash...)
		blockNumKey = append(blockNumKey, toRemove.Bytes()...)

		storage.Clear(blockNumKey)
	}

	storageRootBytes := storage.Root(int32(constants.RuntimeVersion.StateVersion))
	buf.Write(storageRootBytes)
	storageRoot := types.DecodeH256(buf)
	buf.Reset()

	return types.Header{
		ExtrinsicsRoot: extrinsicsRoot,
		StateRoot:      storageRoot,
		ParentHash:     parentHash,
		Number:         blockNumber,
		Digest:         digest,
	}
}

func Initialize(blockNumber types.BlockNumber, parentHash types.Blake2bHash, digest types.Digest) {
	initializationPhase := sc.U32(constants.ExecutionPhaseInitialization)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())

	storage.Set(constants.KeyExtrinsicIndex, sc.U32(0).Bytes())

	numberHash := hashing.Twox128(constants.KeyNumber)
	storage.Set(append(systemHash, numberHash...), blockNumber.Bytes())

	digestHash := hashing.Twox128(constants.KeyDigest)
	storage.Set(append(systemHash, digestHash...), digest.Bytes())

	parentHashKey := hashing.Twox128(constants.KeyParentHash)
	storage.Set(append(systemHash, parentHashKey...), parentHash.Bytes())

	blockHashKeyHash := hashing.Twox128(constants.KeyBlockHash)
	prevBlock := blockNumber - 1
	blockNumHash := hashing.Twox64(prevBlock.Bytes())
	blockNumKey := append(systemHash, blockHashKeyHash...)
	blockNumKey = append(blockNumKey, blockNumHash...)
	blockNumKey = append(blockNumKey, prevBlock.Bytes()...)

	storage.Set(blockNumKey, parentHash.Bytes())

	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)
	storage.Clear(append(systemHash, blockWeightHash...))
}

func IdleAndFinalizeHook(blockNumber types.BlockNumber) {
	systemHash := hashing.Twox128(constants.KeySystem)
	blockWeightHash := hashing.Twox128(constants.KeyBlockWeight)

	storage.Get(append(systemHash, blockWeightHash...))

	// TODO: weights
	/**
	let weight = <frame_system::Pallet<System>>::block_weight();
	let max_weight = <System::BlockWeights as frame_support::traits::Get<_>>::get().max_block;
	let remaining_weight = max_weight.saturating_sub(weight.total());

	if remaining_weight.all_gt(Weight::zero()) {
		let used_weight = <AllPalletsWithSystem as OnIdle<System::BlockNumber>>::on_idle(
			block_number,
			remaining_weight,
		);
		<frame_system::Pallet<System>>::register_extra_weight_unchecked(
			used_weight,
			DispatchClass::Mandatory,
		);
	}
	// Each pallet (babe, grandpa) has its own on_finalize that has to be implemented once it is supported
	<AllPalletsWithSystem as OnFinalize<System::BlockNumber>>::on_finalize(block_number);
	*/
	timestamp.OnFinalize()
}

func NoteFinishedInitialize() {
	initializationPhase := sc.U32(constants.ExecutionPhaseApplyExtrinsic)

	systemHash := hashing.Twox128(constants.KeySystem)
	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(systemHash, executionPhaseHash...), initializationPhase.Bytes())
}

func NoteFinishedExtrinsics() {
	extrinsicIndex := storage.TakeDecode(constants.KeyExtrinsicIndex, sc.DecodeU32)

	systemHash := hashing.Twox128(constants.KeySystem)
	extrinsicCountHash := hashing.Twox128(constants.KeyExtrinsicCount)

	storage.Set(append(systemHash, extrinsicCountHash...), extrinsicIndex.Bytes())

	executionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	finalizationPhase := sc.U32(constants.ExecutionPhaseFinalization)

	storage.Set(append(systemHash, executionPhaseHash...), finalizationPhase.Bytes())
}

func ResetEvents() {
	systemHash := hashing.Twox128(constants.KeySystem)
	eventsHash := hashing.Twox128(constants.KeyEvents)
	eventCountHash := hashing.Twox128(constants.KeyEventCount)
	eventTopicHash := hashing.Twox128(constants.KeyEventTopic)

	storage.Clear(append(systemHash, eventsHash...))
	storage.Clear(append(systemHash, eventCountHash...))

	limit := sc.NewOption[sc.U32](sc.U32(math.MaxUint32))
	storage.ClearPrefix(append(systemHash, eventTopicHash...), limit.Bytes())
}

// Note what the extrinsic data of the current extrinsic index is.
//
// This is required to be called before applying an extrinsic. The data will used
// in [`finalize`] to calculate the correct extrinsics root.
func NoteExtrinsic(encodedExt []byte) {
	keySystemHash := hashing.Twox128(constants.KeySystem)
	keyExtrinsicData := hashing.Twox128(constants.KeyExtrinsicData)

	keyExtrinsicDataPrefixHash := append(keySystemHash, keyExtrinsicData...)
	extrinsicIndex := extrinsicIndexValue()

	hashIndex := hashing.Twox64(extrinsicIndex.Bytes())

	keySystemExtrinsicDataHashIndex := append(keyExtrinsicDataPrefixHash, hashIndex...)
	storage.Set(append(keySystemExtrinsicDataHashIndex, extrinsicIndex.Bytes()...), encodedExt)
}

// To be called immediately after an extrinsic has been applied.
//
// Emits an `ExtrinsicSuccess` or `ExtrinsicFailed` event depending on the outcome.
// The emitted event contains the post-dispatch corrected weight including
// the base-weight for its dispatch class.
func NoteAppliedExtrinsic(r *types.DispatchResultWithPostInfo[types.PostDispatchInfo], info types.DispatchInfo) {
	// TODO:
	// info.Weight = extract_actual_weight(r, &info).saturating_add(T::BlockWeights::get().get(info.class).base_extrinsic)
	// info.PaysFee = extract_actual_pays_fee(r, &info)

	// Self::deposit_event(match r {
	// 	Ok(_) => Event::ExtrinsicSuccess { dispatch_info: info },
	// 	Err(err) => {
	// 		log::trace!(
	// 			target: LOG_TARGET,
	// 			"Extrinsic failed at block({:?}): {:?}",
	// 			Self::block_number(),
	// 			err,
	// 		);
	// 		Event::ExtrinsicFailed { dispatch_error: err.error, dispatch_info: info }
	// 	},
	// });

	nextExtrinsicIndex := extrinsicIndexValue() + sc.U32(1)

	keySystemHash := hashing.Twox128(constants.KeySystem)

	storage.Set(constants.KeyExtrinsicIndex, nextExtrinsicIndex.Bytes())

	keyExecutionPhaseHash := hashing.Twox128(constants.KeyExecutionPhase)
	storage.Set(append(keySystemHash, keyExecutionPhaseHash...), (types.NewExtrinsicPhase(types.PhaseApplyExtrinsic, nextExtrinsicIndex)).Bytes())
}

// Gets the index of extrinsic that is currently executing.
func extrinsicIndexValue() sc.U32 {
	return storage.GetDecode(constants.KeyExtrinsicIndex, sc.DecodeU32)
}

func EnsureInherentsAreFirst(block types.Block) int {
	signedExtrinsicFound := false

	for i, extrinsic := range block.Extrinsics {
		isInherent := false

		if extrinsic.IsSigned() {
			// Signed extrinsics are not inherents
			isInherent = false
		} else {
			call := extrinsic.Function
			// Iterate through all calls and check if the given call is inherent
			switch call.CallIndex.ModuleIndex {
			case timestamp.Module.Index:
				for funcKey := range timestamp.Module.Functions {
					if call.CallIndex.FunctionIndex == timestamp.Module.Functions[funcKey].Index {
						isInherent = true
					}
				}

			}
		}

		if !isInherent {
			signedExtrinsicFound = true
		}

		if signedExtrinsicFound && isInherent {
			return i
		}
	}

	return -1
}

// The current block number being processed. Set by `execute_block`.
func StorageGetBlockNumber() types.BlockNumber {
	systemHash := hashing.Twox128(constants.KeySystem)
	numberHash := hashing.Twox128(constants.KeyNumber)
	return storage.GetDecode(append(systemHash, numberHash...), sc.DecodeU32)
}

func StorageGetAllExtrinsicsLen() sc.U32 {
	systemHash := hashing.Twox128(constants.KeySystem)
	allExtrinsicsLenHash := hashing.Twox128(constants.KeyAllExtrinsicsLen)
	return storage.GetDecode(append(systemHash, allExtrinsicsLenHash...), sc.DecodeU32)
}

func StorageGetAccount(who types.PublicKey) types.AccountInfo {
	systemHash := hashing.Twox128(constants.KeySystem)
	accountHash := hashing.Twox128(constants.KeyAccount)

	whoBytes := sc.FixedSequenceU8ToBytes(who)

	key := append(systemHash, accountHash...)
	key = append(key, hashing.Blake128(whoBytes)...)
	key = append(key, whoBytes...)

	return storage.GetDecode(key, types.DecodeAccountInfo)
}

func StorageSetAccount(who types.PublicKey, account types.AccountInfo) {
	systemHash := hashing.Twox128(constants.KeySystem)
	accountHash := hashing.Twox128(constants.KeyAccount)

	whoBytes := sc.FixedSequenceU8ToBytes(who)

	key := append(systemHash, accountHash...)
	key = append(key, hashing.Blake128(whoBytes)...)
	key = append(key, whoBytes...)

	storage.Set(key, account.Bytes())
}

// Map of block numbers to block hashes.
func StorageGetBlockHash(blockNumber sc.U32) types.Blake2bHash {
	// Module prefix
	systemHash := hashing.Twox128(constants.KeySystem)
	// Storage prefix
	blockHashHash := hashing.Twox128(constants.KeyBlockHash)
	// Block number hash
	blockNumHash := hashing.Twox64(blockNumber.Bytes())

	key := append(systemHash, blockHashHash...)
	key = append(key, blockNumHash...)
	key = append(key, blockNumber.Bytes()...)

	return storage.GetDecode(key, types.DecodeBlake2bHash)
}

// Map of block numbers to block hashes.
func StorageExistsBlockHash(blockNumber sc.U32) sc.Bool {
	// Module prefix
	systemHash := hashing.Twox128(constants.KeySystem)
	// Storage prefix
	blockHashHash := hashing.Twox128(constants.KeyBlockHash)
	// Block number hash
	blockNumHash := hashing.Twox64(blockNumber.Bytes())

	key := append(systemHash, blockHashHash...)
	key = append(key, blockNumHash...)
	key = append(key, blockNumber.Bytes()...)

	return storage.Exists(key) == 1
}
