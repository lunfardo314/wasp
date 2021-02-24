package sbtests

import (
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChainOwnerIDView(t *testing.T) { run2(t, testChainOwnerIDView) }
func testChainOwnerIDView(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncChainOwnerIDView)
	require.NoError(t, err)

	c := ret.MustGet(sbtestsc.ParamChainOwnerID)

	require.EqualValues(t, chain.OriginatorAgentID.Bytes(), c)
}

func TestChainOwnerIDFull(t *testing.T) { run2(t, testChainOwnerIDFull) }
func testChainOwnerIDFull(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncChainOwnerIDFull)
	ret, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	c := ret.MustGet(sbtestsc.ParamChainOwnerID)
	require.EqualValues(t, chain.OriginatorAgentID.Bytes(), c)
}

func TestContractIDView(t *testing.T) { run2(t, testContractIDView) }
func testContractIDView(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncContractIDView)
	require.NoError(t, err)
	cID := coretypes.NewContractID(chain.ChainID, coretypes.Hn(SandboxSCName))
	require.EqualValues(t, cID[:], ret.MustGet(sbtestsc.VarContractID))
}

func TestContractIDFull(t *testing.T) { run2(t, testContractIDFull) }
func testContractIDFull(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncContractIDFull)
	ret, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)
	cID := coretypes.NewContractID(chain.ChainID, coretypes.Hn(SandboxSCName))
	require.EqualValues(t, cID[:], ret.MustGet(sbtestsc.VarContractID))
}

func TestSandboxCall(t *testing.T) { run2(t, testSandboxCall) }
func testSandboxCall(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)

	ret, err := chain.CallView(SandboxSCName, sbtestsc.FuncSandboxCall)
	require.NoError(t, err)

	d := ret.MustGet(sbtestsc.VarSandboxCall)
	require.EqualValues(t, "'solo' testing chain", string(d))
}
