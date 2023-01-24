package storage

import (
	"github.com/LimeChain/gosemble/env"
	"github.com/LimeChain/gosemble/utils"
)

func ChangesRoot(parent_hash int64) int64 { panic("Not implemented!") }

func Clear(key []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	env.ExtStorageClearVersion1(keyOffsetSize)
}

func ClearPrefix(key []byte, limit []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	limitOffsetSize := utils.BytesToOffsetAndSize(limit)
	env.ExtStorageClearPrefixVersion2(keyOffsetSize, limitOffsetSize)
}

func Exists(key int64) int32 { panic("Not implemented!") }

func Get(key []byte) []byte {
	psKey := utils.BytesToOffsetAndSize(key)
	psValue := env.ExtStorageGetVersion1(psKey)
	offset, size := utils.Int64ToOffsetAndSize(psValue)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

func NextKey(key int64) int64 { panic("Not implemented!") }

func Read(key int64, value_out int64, offset int32) int64 { panic("Not implemented!") }

func Root(key []byte) []byte {
	keyOffsetSize := utils.SliceToOffset(key)
	valueOffsetSize := env.ExtStorageRootVersion2(int32(keyOffsetSize))
	offset, size := utils.Int64ToOffsetAndSize(valueOffsetSize)
	value := utils.ToWasmMemorySlice(offset, size)
	return value
}

func Set(key []byte, value []byte) {
	keyOffsetSize := utils.BytesToOffsetAndSize(key)
	valueOffsetSize := utils.BytesToOffsetAndSize(value)
	env.ExtStorageSetVersion1(keyOffsetSize, valueOffsetSize)
}
