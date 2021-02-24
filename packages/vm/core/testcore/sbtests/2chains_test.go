package sbtests

import (
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/coretypes"
	"github.com/iotaledger/wasp/packages/solo"
	"github.com/iotaledger/wasp/packages/vm/core/accounts"
	"github.com/iotaledger/wasp/packages/vm/core/testcore/sbtests/sbtestsc"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test2Chains(t *testing.T) { run2(t, test2Chains) }
func test2Chains(t *testing.T, w bool) {
	env := solo.New(t, false, false)
	chain1 := env.NewChain(nil, "ch1")
	chain2 := env.NewChain(nil, "ch2")
	chain1.CheckAccountLedger()
	chain2.CheckAccountLedger()

	contractID1, _ := setupTestSandboxSC(t, chain1, nil, w)
	contractID2, _ := setupTestSandboxSC(t, chain2, nil, w)

	contractAgentID1 := coretypes.NewAgentIDFromContractID(contractID1)
	contractAgentID2 := coretypes.NewAgentIDFromContractID(contractID2)

	userWallet := env.NewSignatureSchemeWithFunds()
	userAddress := userWallet.Address()
	userAgentID := coretypes.NewAgentIDFromAddress(userAddress)
	env.AssertAddressBalance(userAddress, balance.ColorIOTA, solo.Saldo)

	chain1.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain1.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 0)

	req := solo.NewCallParams(accounts.Interface.Name, accounts.FuncDeposit,
		accounts.ParamAgentID, contractAgentID2,
	).WithTransfer(
		balance.ColorIOTA, 42,
	)
	_, err := chain1.PostRequestSync(req, userWallet)
	require.NoError(t, err)

	accountsAgentID1 := coretypes.NewAgentIDFromContractID(accounts.Interface.ContractID(chain1.ChainID))
	accountsAgentID2 := coretypes.NewAgentIDFromContractID(accounts.Interface.ContractID(chain2.ChainID))

	env.AssertAddressBalance(userAddress, balance.ColorIOTA, solo.Saldo-43)
	chain1.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain2.AssertAccountBalance(userAgentID, balance.ColorIOTA, 0)
	chain1.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain1.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 42)
	chain2.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 0)

	chain1.AssertAccountBalance(accountsAgentID1, balance.ColorIOTA, 0)
	chain1.AssertAccountBalance(accountsAgentID2, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(accountsAgentID1, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(accountsAgentID2, balance.ColorIOTA, 0)

	req = solo.NewCallParams(sbtestsc.Name, sbtestsc.FuncWithdrawToChain,
		sbtestsc.ParamChainID, chain1.ChainID,
	).WithTransfer(
		balance.ColorIOTA, 3,
	)
	_, err = chain2.PostRequestSync(req, userWallet)
	require.NoError(t, err)

	chain1.WaitForEmptyBacklog()
	chain2.WaitForEmptyBacklog()

	env.AssertAddressBalance(userAddress, balance.ColorIOTA, solo.Saldo-47)
	chain1.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain2.AssertAccountBalance(userAgentID, balance.ColorIOTA, 1)
	chain1.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain1.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(contractAgentID1, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(contractAgentID2, balance.ColorIOTA, 43)

	chain1.AssertAccountBalance(accountsAgentID1, balance.ColorIOTA, 1) // !!!! TODO
	chain1.AssertAccountBalance(accountsAgentID2, balance.ColorIOTA, 0)
	chain2.AssertAccountBalance(accountsAgentID1, balance.ColorIOTA, 1) // !!!! TODO
	chain2.AssertAccountBalance(accountsAgentID2, balance.ColorIOTA, 0)
}
