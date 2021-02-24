package sbtests

import (
	"fmt"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address/signaturescheme"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	DEBUG        = false
	ERC20_NAME   = "erc20"
	ERC20_SUPPLY = 100000

	// ERC20 constants
	PARAM_SUPPLY  = "s"
	PARAM_CREATOR = "c"
)

var (
	WasmFileTestcore = "sbtestsc/testcore_bg.wasm"
	WasmFileErc20    = "sbtestsc/erc20_bg.wasm"
	SandboxSCName    = "test_sandbox"
)

func setupChain(t *testing.T, sigSchemeChain signaturescheme.SignatureScheme) (*solo.Solo, *solo.Chain) {
	env := solo.New(t, DEBUG, false)
	chain := env.NewChain(sigSchemeChain, "ch1")
	return env, chain
}

func setupDeployer(t *testing.T, chain *solo.Chain) signaturescheme.SignatureScheme {
	user := chain.Env.NewSignatureSchemeWithFunds()
	chain.Env.AssertAddressBalance(user.Address(), balance.ColorIOTA, testutil.RequestFundsAmount)

	req := solo.NewCallParams(root.Interface.Name, root.FuncGrantDeploy,
		root.ParamDeployer, coretypes.NewAgentIDFromAddress(user.Address()),
	)
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)
	return user
}

func run2(t *testing.T, test func(*testing.T, bool), skipWasm ...bool) {
	t.Run(fmt.Sprintf("run CORE version of %s", t.Name()), func(t *testing.T) {
		test(t, false)
	})
	if len(skipWasm) == 0 || !skipWasm[0] {
		t.Run(fmt.Sprintf("run WASM version of %s", t.Name()), func(t *testing.T) {
			test(t, true)
		})
	} else {
		t.Logf("skipped WASM version of '%s'", t.Name())
	}
}

func setupTestSandboxSC(t *testing.T, chain *solo.Chain, user signaturescheme.SignatureScheme, runWasm bool) (coretypes.ContractID, int64) {
	var err error
	var extraToken int64
	if runWasm {
		err = chain.DeployWasmContract(user, SandboxSCName, WasmFileTestcore)
		extraToken = 1
	} else {
		err = chain.DeployContract(user, SandboxSCName, sbtestsc.Interface.ProgramHash)
		extraToken = 0
	}
	require.NoError(t, err)

	deployed := coretypes.NewContractID(chain.ChainID, coretypes.Hn(sbtestsc.Interface.Name))
	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncDoNothing)
	_, err = chain.PostRequestSync(req, user)
	require.NoError(t, err)
	t.Logf("deployed test_sandbox'%s': %s", SandboxSCName, coretypes.Hn(SandboxSCName))
	return deployed, extraToken
}

func setupERC20(t *testing.T, chain *solo.Chain, user signaturescheme.SignatureScheme, runWasm bool) coretypes.ContractID {
	var err error
	if !runWasm {
		t.Logf("skipped %s. Only for Wasm tests, always loads %s", t.Name(), WasmFileErc20)
		return coretypes.ContractID{}
	}
	var userAgentID coretypes.AgentID
	if user == nil {
		userAgentID = chain.OriginatorAgentID
	} else {
		userAgentID = coretypes.NewAgentIDFromAddress(user.Address())
	}
	err = chain.DeployWasmContract(user, ERC20_NAME, WasmFileErc20,
		PARAM_SUPPLY, 1000000,
		PARAM_CREATOR, userAgentID,
	)
	require.NoError(t, err)

	deployed := coretypes.NewContractID(chain.ChainID, coretypes.Hn(sbtestsc.Interface.Name))
	t.Logf("deployed erc20'%s': %s --  %s", ERC20_NAME, coretypes.Hn(ERC20_NAME), deployed)
	return deployed
}

func TestSetup1(t *testing.T) { run2(t, testSetup1) }
func testSetup1(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	setupTestSandboxSC(t, chain, nil, w)
}

func TestSetup2(t *testing.T) { run2(t, testSetup2) }
func testSetup2(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	user := setupDeployer(t, chain)
	setupTestSandboxSC(t, chain, user, w)
}

func TestSetup3(t *testing.T) { run2(t, testSetup3) }
func testSetup3(t *testing.T, w bool) {
	_, chain := setupChain(t, nil)
	user := setupDeployer(t, chain)
	setupTestSandboxSC(t, chain, user, w)
	//setupERC20(t, chain, user, w)
}
