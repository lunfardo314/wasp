// Copyright 2020 IOTA Stiftung
// SPDX-License-Identifier: Apache-2.0

package wasmlib

type MapKey interface {
	KeyId() Key32
}

type Key string

func (key Key) KeyId() Key32 {
	return GetKeyIdFromString(string(key))
}

type Key32 int32

func (key Key32) KeyId() Key32 {
	return key
}

// \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\ // \\

const (
	KeyAddress         = Key32(-1)
	KeyAggregateBls    = Key32(-2)
	KeyBalances        = Key32(-3)
	KeyBase58Bytes     = Key32(-4)
	KeyBase58String    = Key32(-5)
	KeyCall            = Key32(-6)
	KeyCaller          = Key32(-7)
	KeyChainOwnerId    = Key32(-8)
	KeyColor           = Key32(-9)
	KeyContractCreator = Key32(-10)
	KeyContractId      = Key32(-11)
	KeyData            = Key32(-12)
	KeyDeploy          = Key32(-13)
	KeyEvent           = Key32(-14)
	KeyExports         = Key32(-15)
	KeyHashBlake2b     = Key32(-16)
	KeyHashSha3        = Key32(-17)
	KeyHname           = Key32(-18)
	KeyIncoming        = Key32(-19)
	KeyLength          = Key32(-20)
	KeyLog             = Key32(-21)
	KeyMaps            = Key32(-22)
	KeyName            = Key32(-23)
	KeyPanic           = Key32(-24)
	KeyParams          = Key32(-25)
	KeyPost            = Key32(-26)
	KeyRandom          = Key32(-27)
	KeyResults         = Key32(-28)
	KeyReturn          = Key32(-29)
	KeyState           = Key32(-30)
	KeyTimestamp       = Key32(-31)
	KeyTrace           = Key32(-32)
	KeyTransfers       = Key32(-33)
	KeyUtility         = Key32(-34)
	KeyValid           = Key32(-35)
	KeyValidBls        = Key32(-36)
	KeyValidEd25519    = Key32(-37)
	KeyZzzzzzz         = Key32(-98)
)
