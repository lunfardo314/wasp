package blob

import (
	"fmt"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/coretypes/assert"
	"github.com/iotaledger/wasp/packages/kv/kvdecoder"

	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/kv/codec"
	"github.com/iotaledger/wasp/packages/kv/dict"
)

func initialize(ctx coretypes.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("blob.initialize.success hname = %s", Interface.Hname().String())
	return nil, nil
}

// storeBlob treats parameters as names of fields and field values
// it stores it in the state in deterministic binary representation
// Returns hash of the blob
func storeBlob(ctx coretypes.Sandbox) (dict.Dict, error) {
	ctx.Log().Debugf("blob.storeBlob.begin")
	state := ctx.State()
	params := ctx.Params()
	// calculate a deterministic hash of all blob fields
	blobHash, kSorted, values := mustGetBlobHash(params)

	directory := GetDirectory(state)
	assert.NewAssert(ctx.Log()).Require(!directory.MustHasAt(blobHash[:]),
		"blob.storeBlob.fail: blob with hash %s already exist", blobHash.String())

	// get a record by blob hash
	blbValues := GetBlobValues(state, blobHash)
	blbSizes := GetBlobSizes(state, blobHash)

	totalSize := uint32(0)

	// save record of the blob. In parallel save record of sizes of blob fields
	sizes := make([]uint32, len(kSorted))
	for i, k := range kSorted {
		size := uint32(len(values[i]))

		blbValues.MustSetAt([]byte(k), values[i])
		blbSizes.MustSetAt([]byte(k), EncodeSize(size))
		sizes[i] = size
		totalSize += size
	}

	ret := dict.New()
	ret.Set(ParamHash, codec.EncodeHashValue(blobHash))

	directory.MustSetAt(blobHash[:], EncodeSize(totalSize))

	ctx.Event(fmt.Sprintf("[blob] hash: %s, field sizes: %+v", blobHash.String(), sizes))
	return ret, nil
}

// getBlobInfo return lengths of all fields in the blob
func getBlobInfo(ctx coretypes.SandboxView) (dict.Dict, error) {
	ctx.Log().Debugf("blob.getBlobInfo.begin")

	params := kvdecoder.New(ctx.Params(), ctx.Log())
	blobHash := params.MustGetHashValue(ParamHash)

	blbSizes := GetBlobSizesR(ctx.State(), blobHash)
	ret := dict.New()
	blbSizes.MustIterate(func(field []byte, value []byte) bool {
		ret.Set(kv.Key(field), value)
		return true
	})
	return ret, nil
}

func getBlobField(ctx coretypes.SandboxView) (dict.Dict, error) {
	ctx.Log().Debugf("blob.getBlobField.begin")
	state := ctx.State()

	params := kvdecoder.New(ctx.Params(), ctx.Log())
	blobHash := params.MustGetHashValue(ParamHash)
	field := params.MustGetBytes(ParamField)

	blobValues := GetBlobValuesR(state, blobHash)
	if blobValues.MustLen() == 0 {
		return nil, fmt.Errorf("blob with hash %s has not been found", blobHash.String())
	}
	value := blobValues.MustGetAt(field)
	if value == nil {
		return nil, fmt.Errorf("'blob field %s value not found", string(field))
	}
	ret := dict.New()
	ret.Set(ParamBytes, value)
	return ret, nil
}

func listBlobs(ctx coretypes.SandboxView) (dict.Dict, error) {
	ctx.Log().Debugf("blob.listBlobs.begin")
	ret := dict.New()
	GetDirectoryR(ctx.State()).MustIterate(func(hash []byte, totalSize []byte) bool {
		ret.Set(kv.Key(hash), totalSize)
		return true
	})
	return ret, nil
}
