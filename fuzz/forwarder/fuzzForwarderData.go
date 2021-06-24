package fuzzForwarder

type programmedCallType int

const (
	syncCall programmedCallType = iota
	asyncCall
	transferExecute
)

func (ct programmedCallType) String() string {
	switch ct {
	case syncCall:
		return "syncCall"
	case asyncCall:
		return "asyncCall"
	case transferExecute:
		return "transferExecute"
	default:
		panic("unknown programmedCallType")
	}
}

type programmedCall struct {
	callType  programmedCallType
	fromIndex int
	toIndex   int
	token     string
	nonce     int
	amount    string
}

type fuzzData struct {
	mainCallerAddress     string
	numForwarders         int
	maxCallDepth          int
	programmedCalls       map[int][]*programmedCall
	numFungibleTokens     int
	numSemiFungibleTokens int
}
