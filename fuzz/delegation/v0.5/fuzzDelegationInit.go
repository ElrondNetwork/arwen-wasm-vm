package delegation

import (
	"fmt"
)

func (pfe *fuzzDelegationExecutor) init(args *fuzzDelegationExecutorInitArgs) error {
	pfe.serviceFee = args.serviceFee
	pfe.numBlocksBeforeForceUnstake = args.numBlocksBeforeForceUnstake
	pfe.numBlocksBeforeUnbond = args.numBlocksBeforeUnbond
	pfe.numDelegators = args.numDelegators
	pfe.stakePerNode = args.stakePerNode

	pfe.world.Clear()

	pfe.ownerAddress = "address:fuzz-owner"
	pfe.delegationContractAddress = "sc:fuzz-delegation"
	pfe.auctionMockAddress = "sc:fuzz-auction-mock"
	pfe.faucetAddress = "address:endless-sack-of-erd"
	pfe.withdrawTargetAddress = "address:withdraw-target"
	pfe.stakePurchaseForwardAddress = "address:stake-purchase-forwarded"

	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			},
			"%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {
					"str:stake_per_node": "%d"
				},
				"code": "file:auction-mock.wasm"
			},
			"%s": {
				"nonce": "0",
				"balance": "1,000,000,000,000,000,000,000,000,000,000",
				"storage": {},
				"code": ""
			},
			"%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			},
			"%s": {
				"nonce": "0",
				"balance": "0",
				"storage": {},
				"code": ""
			}
		},
		"newAddresses": [
			{
				"creatorAddress": "%s",
				"creatorNonce": "0",
				"newAddress": "%s"
			}
		]
	}`,
		pfe.ownerAddress,
		pfe.auctionMockAddress,
		pfe.stakePerNode,
		pfe.faucetAddress,
		pfe.withdrawTargetAddress,
		pfe.stakePurchaseForwardAddress,
		pfe.ownerAddress,
		pfe.delegationContractAddress,
	))
	if err != nil {
		return err
	}

	// delegators
	for i := 1; i <= args.numDelegators; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"storage": {},
					"code": ""
				}
			}
		}`,
			pfe.delegatorAddress(i),
		))
		if err != nil {
			return err
		}
	}

	// deploy delegation
	_, err = pfe.executeTxStep(fmt.Sprintf(`
	{
		"step": "scDeploy",
		"txId": "-deploy-",
		"tx": {
			"from": "%s",
			"value": "0",
			"contractCode": "file:delegation.wasm",
			"arguments": [
				"%s",
				"%d",
				"%d",
				"%d",
				"%d",
				"%d"
			],
			"gasLimit": "50,000,000",
			"gasPrice": "0"
		},
		"expect": {
			"out": [],
			"status": "",
			"logs": [],
			"gas": "*",
			"refund": "*"
		}
	}`,
		pfe.ownerAddress,
		pfe.auctionMockAddress,
		args.serviceFee,
		args.ownerMinStake,
		args.numBlocksBeforeUnbond,
		args.minStake,
		args.totalDelegationCap,
	))
	if err != nil {
		return err
	}

	pfe.log("init ok")
	return nil
}
