package blob

import (
	"fmt"

	"github.com/iotaledger/wasp/packages/hashing"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/collections"
	"github.com/iotaledger/wasp/packages/kv/dict"
	"github.com/iotaledger/wasp/packages/util"
)

const varStateDirectory = "d"

func valuesKey(blobHash hashing.HashValue) string {
	return "v" + string(blobHash[:])
}

func sizesKey(blobHash hashing.HashValue) string {
	return "s" + string(blobHash[:])
}

func mustGetBlobHash(fields dict.Dict) (hashing.HashValue, []kv.Key, [][]byte) {
	kSorted := fields.KeysSorted() // mind determinism
	values := make([][]byte, 0, len(kSorted))
	all := make([][]byte, 0, 2*len(kSorted))
	for _, k := range kSorted {
		v := fields.MustGet(k)
		values = append(values, v)
		all = append(all, v)
		all = append(all, []byte(k))
	}
	return hashing.HashData(all...), kSorted, values
}

// MustGetBlobHash deterministically hashes map of binary values
func MustGetBlobHash(fields dict.Dict) hashing.HashValue {
	ret, _, _ := mustGetBlobHash(fields)
	return ret
}

// GetDirectory retrieves the blob directory from the state
func GetDirectory(state kv.KVStore) *collections.Map {
	return collections.NewMap(state, varStateDirectory)
}

// GetDirectoryR retrieves the blob directory from the read-only state
func GetDirectoryR(state kv.KVStoreReader) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, varStateDirectory)
}

// GetBlobValues retrieves the blob field-value map from the state
func GetBlobValues(state kv.KVStore, blobHash hashing.HashValue) *collections.Map {
	return collections.NewMap(state, valuesKey(blobHash))
}

// GetBlobValuesR retrieves the blob field-value map from the read-only state
func GetBlobValuesR(state kv.KVStoreReader, blobHash hashing.HashValue) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, valuesKey(blobHash))
}

// GetBlobSizes retrieves the writeable blob field-size map from the state
func GetBlobSizes(state kv.KVStore, blobHash hashing.HashValue) *collections.Map {
	return collections.NewMap(state, sizesKey(blobHash))
}

// GetBlobSizesR retrieves the blob field-size map from the read-only state
func GetBlobSizesR(state kv.KVStoreReader, blobHash hashing.HashValue) *collections.ImmutableMap {
	return collections.NewMapReadOnly(state, sizesKey(blobHash))
}

func LocateProgram(state kv.KVStoreReader, programHash hashing.HashValue) (string, []byte, error) {
	blbValues := GetBlobValuesR(state, programHash)
	programBinary := blbValues.MustGetAt([]byte(VarFieldProgramBinary))
	if programBinary == nil {
		return "", nil, fmt.Errorf("can't find program binary for hash %s", programHash.String())
	}
	v := blbValues.MustGetAt([]byte(VarFieldVMType))
	vmType := "wasmtimevm"
	if v != nil {
		vmType = string(v)
	}
	return vmType, programBinary, nil
}

func EncodeSize(size uint32) []byte {
	return util.Uint32To4Bytes(size)
}

func DecodeSize(size []byte) (uint32, error) {
	return util.Uint32From4Bytes(size)
}

func DecodeSizesMap(sizes dict.Dict) (map[string]uint32, error) {
	ret := make(map[string]uint32)
	for field, size := range sizes {
		v, err := DecodeSize(size)
		if err != nil {
			return nil, err
		}
		ret[string(field)] = uint32(v)
	}
	return ret, nil
}

func DecodeDirectory(blobs dict.Dict) (map[hashing.HashValue]uint32, error) {
	ret := make(map[hashing.HashValue]uint32)
	for hash, size := range blobs {
		v, err := DecodeSize(size)
		if err != nil {
			return nil, err
		}
		h, _, err := codec.DecodeHashValue([]byte(hash))
		if err != nil {
			return nil, err
		}
		ret[h] = v
	}
	return ret, nil
}
