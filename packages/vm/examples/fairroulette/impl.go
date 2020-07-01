package fairroulette

import (
	"bytes"
	"encoding/json"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/address"
	"github.com/iotaledger/goshimmer/dapps/valuetransfers/packages/balance"
	"github.com/iotaledger/wasp/packages/sctransaction"
	"github.com/iotaledger/wasp/packages/util"
	"github.com/iotaledger/wasp/packages/vm/vmtypes"
	"github.com/mr-tron/base58"
	"io"
	"sort"
)

type fairRouletteProcessor map[sctransaction.RequestCode]fairRouletteEntryPoint

type fairRouletteEntryPoint func(ctx vmtypes.Sandbox)

const (
	RequestPlaceBet          = sctransaction.RequestCode(uint16(1))
	RequestVoteForPlay       = sctransaction.RequestCode(uint16(2))
	RequestPlayAndDistribute = sctransaction.RequestCode(uint16(3))
)

var entryPoints = fairRouletteProcessor{
	RequestPlaceBet:          placeBet,
	RequestVoteForPlay:       vote,
	RequestPlayAndDistribute: playAndDistribute,
}

const (
	ProgramHash = "3wo28GRrJu37v6D4xkjZsRLiVQrk3iMn7PifpMFoJEiM"

	ReqVarColor              = "color"
	StateVarBets             = "bets"
	StateVarNumBets          = "numBets"
	StateVarLockedBets       = "lockedBest"
	StateVarNumLockedBets    = "numBets"
	StateVarNumVotes         = "numvotes"
	StateVarLastWinningColor = "lastWinningColor"

	NumColors       = 8
	NumVotesForPlay = 10
)

type betInfo struct {
	player address.Address
	reqId  sctransaction.RequestId
	sum    int64
	color  byte
}

// all strings base58
type betInfoJson struct {
	playerAddr string `json:"player_addr"`
	reqId      string `json:"req_id"`
	sum        int64  `json:"sum"`
	color      byte   `json:"color"`
}

func GetProcessor() vmtypes.Processor {
	return entryPoints
}

func (f fairRouletteProcessor) GetEntryPoint(code sctransaction.RequestCode) (vmtypes.EntryPoint, bool) {
	ep, ok := entryPoints[code]
	return ep, ok
}

func (f fairRouletteEntryPoint) WithGasLimit(i int) vmtypes.EntryPoint {
	return f
}

func (f fairRouletteEntryPoint) Run(ctx vmtypes.Sandbox) {
	f(ctx)
}

// the request places bet into the smart contract
// all current bets are kept as one marshalled binary blob in state variable 'StateVarBets'
func placeBet(ctx vmtypes.Sandbox) {
	// take senders. Must be exactly 1
	senders := ctx.AccessRequest().Senders()
	if len(senders) != 1 {
		return
	}
	sender := senders[0]
	// look if there're some iotas left for the bet.
	// it is after min rewards. Here we accessing only part which is coming with the current request
	sum := ctx.AccessAccount().AvailableBalanceFromRequest(&balance.ColorIOTA)
	if sum == 0 {
		// nothing to bet
		return
	}
	// see if there's a color among args
	col, ok := ctx.AccessRequest().GetInt64(ReqVarColor)
	if !ok {
		// wrong request, no color specified
		return
	}
	data, ok := ctx.AccessState().Get(StateVarBets)
	var bets []*betInfo
	if ok {
		bets = decodeBets(data)
	} else {
		bets = make([]*betInfo, 0)
	}
	bets = append(bets, &betInfo{
		player: sender,
		sum:    sum,
		reqId:  ctx.AccessRequest().ID(),
		color:  byte(col % NumColors),
	})
	ctx.AccessState().Set(StateVarBets, encodeBets(bets))
	ctx.AccessState().SetInt64(StateVarNumBets, int64(len(bets)))
}

// anyone can vote, they can't predict the outcome anyway
// alternatively, only betters could be allowed to bet --> need for hashmap structure
func vote(ctx vmtypes.Sandbox) {
	numVotes, _, _ := ctx.AccessState().GetInt64(StateVarNumVotes)
	if numVotes+1 < NumVotesForPlay {
		ctx.AccessState().SetInt64(StateVarNumVotes, numVotes+1)
		return
	}
	// number of votes reached NumVotesForPlay.
	// Lock current bets and send the 'PlayAndDistribute' request to itself
	// get locked bets
	lockedBetsData, ok := ctx.AccessState().Get(StateVarLockedBets)
	var lockedBets []*betInfo
	if ok {
		lockedBets = decodeBets(lockedBetsData)
	} else {
		lockedBets = make([]*betInfo, 0)
	}
	// get current bets
	data, ok := ctx.AccessState().Get(StateVarBets)
	var bets []*betInfo
	if ok {
		bets = decodeBets(data)
	} else {
		bets = make([]*betInfo, 0)
	}
	// append current bets to locked bets
	lockedBets = append(lockedBets, bets...)
	// store locked bets
	ctx.AccessState().Set(StateVarLockedBets, encodeBets(lockedBets))
	ctx.AccessState().SetInt64(StateVarNumLockedBets, int64(len(lockedBets)))
	// clear current bets
	ctx.AccessState().Del(StateVarBets)
	ctx.AccessState().SetInt64(StateVarNumBets, 0)

	ctx.SendRequestToSelf(RequestPlayAndDistribute, nil)
	ctx.AccessState().SetInt64(StateVarNumVotes, 0)
}

func playAndDistribute(ctx vmtypes.Sandbox) {
	if !ctx.AccessRequest().IsAuthorisedByAddress(ctx.GetOwnAddress()) {
		// ignore if request is not from itself
		return
	}
	numLocked, _, _ := ctx.AccessState().GetInt64(StateVarNumLockedBets)
	if numLocked == 0 {
		// nothing is to play
		return
	}

	// entropy includes signature of the locked bets. It was not possible to predict it
	// at the moment of locking
	entropy := ctx.GetEntropy()
	winningColor := byte(util.Uint64From8Bytes(entropy[:8]) / NumColors)
	ctx.AccessState().SetInt64(StateVarLastWinningColor, int64(winningColor))

	// take locked bets
	lockedBetsData, ok := ctx.AccessState().Get(StateVarLockedBets)
	var lockedBets []*betInfo
	if ok {
		lockedBets = decodeBets(lockedBetsData)
	} else {
		lockedBets = make([]*betInfo, 0)
	}
	totalLockedAmount := int64(0)
	for _, bet := range lockedBets {
		totalLockedAmount += bet.sum
	}
	// select bets on winning color
	winningBets := lockedBets[:0] // same underlying array
	for _, bet := range lockedBets {
		if bet.color == winningColor {
			winningBets = append(winningBets, bet)
		}
	}

	ctx.AccessState().Del(StateVarLockedBets)
	ctx.AccessState().SetInt64(StateVarNumLockedBets, 0)
	ctx.AccessState().SetInt64(StateVarNumVotes, 0)

	if len(winningBets) == 0 {
		// nobody played on winning color -> all sums stay in smart contract
		// move tokens to itself in order to compress number of outputs in the address
		if !ctx.AccessAccount().MoveTokens(ctx.GetOwnAddress(), &balance.ColorIOTA, totalLockedAmount) {
			ctx.Rollback()
			return
		}
	}

	if !distributeLockedAmount(ctx, winningBets, totalLockedAmount) {
		ctx.Rollback()
		return
	}
}

func distributeLockedAmount(ctx vmtypes.Sandbox, bets []*betInfo, totalLockedAmount int64) bool {
	sumsByPlayers := make(map[address.Address]int64)
	totalWinningAmount := int64(0)
	for _, bet := range bets {
		if _, ok := sumsByPlayers[bet.player]; !ok {
			sumsByPlayers[bet.player] = 0
		}
		sumsByPlayers[bet.player] += bet.sum
		totalWinningAmount += bet.sum
	}

	// NOTE 1: float64 was avoided for determinism reasons
	// NOTE: beware overflows
	for player, sum := range sumsByPlayers {
		sumsByPlayers[player] = (totalLockedAmount * sum) / totalWinningAmount
	}
	// make deterministic sequence
	seqPlayers := make([]address.Address, 0, len(sumsByPlayers))
	resulSum := int64(0)
	for player, sum := range sumsByPlayers {
		seqPlayers = append(seqPlayers, player)
		resulSum += sum
	}
	sort.Slice(seqPlayers, func(i, j int) bool {
		return bytes.Compare(seqPlayers[i][:], seqPlayers[j][:]) < 0
	})

	if resulSum > totalLockedAmount {
		sumsByPlayers[seqPlayers[0]] -= resulSum - totalLockedAmount
	}
	finalWinners := seqPlayers[:0]
	for _, player := range seqPlayers {
		if sumsByPlayers[player] <= 0 {
			continue
		}
		finalWinners = append(finalWinners, player)
	}
	for _, player := range finalWinners {
		if !ctx.AccessAccount().MoveTokens(&player, &balance.ColorIOTA, sumsByPlayers[player]) {
			return false
		}
	}
	return true
}

func (bi *betInfo) Write(w io.Writer) error {
	_, _ = w.Write(bi.player[:])
	_ = util.WriteInt64(w, bi.sum)
	_ = util.WriteByte(w, bi.color)
	return nil
}

func (bi *betInfo) Read(r io.Reader) error {
	var err error
	if err = util.ReadAddress(r, &bi.player); err != nil {
		return err
	}
	if err = util.ReadInt64(r, &bi.sum); err != nil {
		return err
	}
	if bi.color, err = util.ReadByte(r); err != nil {
		return err
	}
	return nil
}

func toJsonable(bi *betInfo) *betInfoJson {
	return &betInfoJson{
		playerAddr: bi.player.String(),
		reqId:      base58.Encode(bi.reqId[:]),
		sum:        bi.sum,
		color:      bi.color,
	}
}

func fromJsonable(biJson *betInfoJson) *betInfo {
	playerAddr, err := address.FromBase58(biJson.playerAddr)
	if err != nil {
		playerAddr = address.Address{}
	}
	reqId, err := sctransaction.NewRequestIdFromString(biJson.reqId)
	if err != nil {
		reqId = sctransaction.RequestId{}
	}

	return &betInfo{
		player: playerAddr,
		reqId:  reqId,
		sum:    biJson.sum,
		color:  biJson.color,
	}
}

func encodeBets(bets []*betInfo) []byte {
	betsJson := make([]*betInfoJson, len(bets))
	for i, bi := range bets {
		betsJson[i] = toJsonable(bi)
	}
	data, _ := json.Marshal(betsJson)
	return data
}

func decodeBets(data []byte) []*betInfo {
	tmpLst := make([]*betInfoJson, 0)
	if err := json.Unmarshal(data, &tmpLst); err != nil {
		return []*betInfo{}
	}

	ret := make([]*betInfo, len(tmpLst))
	for i := range ret {
		ret[i] = fromJsonable(tmpLst[i])
	}
	return ret
}