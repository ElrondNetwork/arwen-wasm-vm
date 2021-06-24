package fuzzForwarder

import (
	"fmt"
	"strings"
)

func (pfe *fuzzExecutor) initData() error {
	pfe.data = &fuzzData{
		mainCallerAddress:     "address:main_caller",
		numForwarders:         5,
		maxCallDepth:          9,
		programmedCalls:       make(map[int][]*programmedCall),
		numFungibleTokens:     3,
		numSemiFungibleTokens: 3,
	}

	pfe.world.Clear()
	return nil
}

func (pfe *fuzzExecutor) setUp() error {
	err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"esdt": {
						%s
					}
				}
			}
		}`,
		pfe.data.mainCallerAddress,
		pfe.setUpTokens(),
	))
	if err != nil {
		return err
	}
	for i := 1; i <= pfe.data.numForwarders; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "10000000",
					"esdt": {
						%s
					},
					"storage": {},
					"code": "file:../forwarder/output/forwarder.wasm"
				}
			}
		}`,
			pfe.forwarderAddress(i),
			pfe.setUpTokens(),
		))
		if err != nil {
			return err
		}
	}
	return nil
}

func (pfe *fuzzExecutor) setUpTokens() string {
	var sb strings.Builder
	first := true
	for i := 1; i <= pfe.data.numFungibleTokens; i++ {
		if first {
			first = false
		} else {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`
			"str:%s": "1,000,000,000,000"`,
			pfe.fungibleTokenName(i)))
	}
	for i := 1; i <= pfe.data.numSemiFungibleTokens; i++ {
		if first {
			first = false
		} else {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`
			"str:%s": {
				"nonce": "1",
				"balance": "1,000,000,000,000"
			}`,
			pfe.semiFungibleTokenName(i)))
	}
	return sb.String()
}
