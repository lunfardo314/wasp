package sbtests

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/testutil"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/root"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDoNothing(t *testing.T) { run2(t, testDoNothing) }
func testDoNothing(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncDoNothing).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	t.Logf("dump accounts:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 42)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-42-extraToken)
}

func TestDoNothingUser(t *testing.T) { run2(t, testDoNothingUser) }
func testDoNothingUser(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)

	userAgentID := coretypes.NewAgentIDFromAddress(user.Address())

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncDoNothing).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, user)
	require.NoError(t, err)

	t.Logf("dump accounts:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 42)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(user.Address(), balance.ColorIOTA, testutil.RequestFundsAmount-1-42)
}

func TestWithdrawToAddress(t *testing.T) { run2(t, testWithdrawToAddress) }
func testWithdrawToAddress(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)
	t.Logf("contract agentID: %s", cAID)

	userAddress := user.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)

	req := solo.NewCallParams(SandboxSCName, sbtestsc.FuncDoNothing).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, user)
	require.NoError(t, err)

	t.Logf("dump accounts 1:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 42)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-42)

	req = solo.NewCallParams(SandboxSCName, sbtestsc.FuncSendToAddress,
		sbtestsc.ParamAddress, userAddress,
	)
	_, err = chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	t.Logf("dump accounts 2:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 5+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-5-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1)
}

func TestDoPanicUser(t *testing.T) { run2(t, testDoPanicUser) }
func testDoPanicUser(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)

	userAddress := user.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)

	t.Logf("dump accounts 1:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)

	req := solo.NewCallParams(sbtestsc.Interface.Name, sbtestsc.FuncPanicFullEP).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, user)
	require.Error(t, err)

	t.Logf("dump accounts 2:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1)
}

func TestDoPanicUserFeeless(t *testing.T) { run2(t, testDoPanicUserFeeless) }
func testDoPanicUserFeeless(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)

	userAddress := user.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)

	t.Logf("dump accounts 1:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)

	req := solo.NewCallParams(sbtestsc.Interface.Name, sbtestsc.FuncPanicFullEP).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, user)
	require.Error(t, err)

	t.Logf("dump accounts 2:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1)

	req = solo.NewCallParams(accounts.Interface.Name, accounts.FuncWithdrawToAddress)
	_, err = chain.PostRequestSync(req, user)
	require.NoError(t, err)

	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)
}

func TestDoPanicUserFee(t *testing.T) { run2(t, testDoPanicUserFee) }
func testDoPanicUserFee(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)

	userAddress := user.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)

	t.Logf("dump accounts 1:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)

	req := solo.NewCallParams(root.Interface.Name, root.FuncSetContractFee,
		root.ParamHname, cID.Hname(),
		root.ParamOwnerFee, 10,
	)
	_, err := chain.PostRequestSync(req, nil)
	require.NoError(t, err)

	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 5+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-5-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)

	req = solo.NewCallParams(sbtestsc.Interface.Name, sbtestsc.FuncPanicFullEP).
		WithTransfer(balance.ColorIOTA, 42)
	_, err = chain.PostRequestSync(req, user)
	require.Error(t, err)

	t.Logf("dump accounts 2:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 5+10+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-5-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-10)
}

func TestRequestToView(t *testing.T) { run2(t, testRequestToView) }
func testRequestToView(t *testing.T, w bool) {
	env, chain := setupChain(t, nil)
	cID, extraToken := setupTestSandboxSC(t, chain, nil, w)
	cAID := coretypes.NewAgentIDFromContractID(cID)
	user := setupDeployer(t, chain)

	userAddress := user.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)

	t.Logf("dump accounts 1:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount)

	// sending request to the view entry point should return an error and invoke fallback for tokens
	req := solo.NewCallParams(sbtestsc.Interface.Name, sbtestsc.FuncJustView).
		WithTransfer(balance.ColorIOTA, 42)
	_, err := chain.PostRequestSync(req, user)
	require.Error(t, err)

	t.Logf("dump accounts 2:\n%s", chain.DumpAccounts())
	chain.AssertAccountBalance(chain.OriginatorAgentID, balance.ColorIOTA, 4+extraToken)
	chain.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain.AssertAccountBalance(cAID, balance.ColorIOTA, 0)
	env.AssertAddressBalance(chain.OriginatorAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1-4-extraToken)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, testutil.RequestFundsAmount-1)
}
