package dex

import (
	"fmt"
)

func (pfe *fuzzDexExecutor) init(args *fuzzDexExecutorInitArgs) error {
	pfe.wegldTokenId = args.wegldTokenId
	pfe.mexTokenId = args.mexTokenId
	pfe.busdTokenId = args.busdTokenId
	pfe.wemeLpTokenId = args.wemeLpTokenId
	pfe.webuLpTokenId = args.webuLpTokenId
	pfe.wemeFarmTokenId = args.wemeFarmTokenId
	pfe.webuFarmTokenId = args.webuFarmTokenId
	pfe.mexFarmTokenId = args.mexFarmTokenId
	pfe.numUsers = args.numUsers
	pfe.numEvents = args.numEvents
	pfe.removeLiquidityProb = args.removeLiquidityProb
	pfe.addLiquidityProb = args.addLiquidityProb
	pfe.swapProb = args.swapProb
	pfe.queryPairsProb = args.queryPairsProb
	pfe.enterFarmProb = args.enterFarmProb
	pfe.exitFarmProb = args.exitFarmProb
	pfe.claimRewardsProb = args.claimRewardsProb
	pfe.compoundRewardsProb = args.compoundRewardsProb
	pfe.increaseBlockNonceProb = args.increaseBlockNonceProb
	pfe.removeLiquidityMaxValue = args.removeLiquidityMaxValue
	pfe.addLiquidityMaxValue = args.addLiquidityMaxValue
	pfe.swapMaxValue = args.swapMaxValue
	pfe.enterFarmMaxValue = args.enterFarmMaxValue
	pfe.exitFarmMaxValue = args.exitFarmMaxValue
	pfe.claimRewardsMaxValue = args.claimRewardsMaxValue
	pfe.compoundRewardsMaxValue = args.compoundRewardsMaxValue
	pfe.tokenDepositMaxValue = args.tokenDepositMaxValue
	pfe.blockNonceIncrease = args.blockNonceIncrease
	pfe.farmers = make(map[int]FarmerInfo)
	pfe.currentFarmTokenNonce = make(map[string]int)

	pfe.world.Clear()

	pfe.ownerAddress = "address:fuzz_owner"
	pfe.wemeFarmAddress = "sc:fuzz_dex_weme_farm"
	pfe.webuFarmAddress = "sc:fuzz_dex_webu_farm"
	pfe.mexFarmAddress = "sc:fuzz_dex_mex_farm"
	pfe.wemeSwapAddress = "sc:fuzz_dex_weme_swap"
	pfe.webuSwapAddress = "sc:fuzz_dex_webu_swap"

	pfe.currentFarmTokenNonce[pfe.wemeFarmTokenId] = 0
	pfe.currentFarmTokenNonce[pfe.webuFarmTokenId] = 0
	pfe.currentFarmTokenNonce[pfe.mexFarmTokenId] = 0

	pfe.farms[0] = Farm{
		address:      pfe.wemeFarmAddress,
		farmToken:    pfe.wemeFarmTokenId,
		farmingToken: pfe.wemeLpTokenId,
		rewardToken:  pfe.mexTokenId,
	}
	pfe.farms[1] = Farm{
		address:      pfe.webuFarmAddress,
		farmToken:    pfe.webuFarmTokenId,
		farmingToken: pfe.webuLpTokenId,
		rewardToken:  pfe.mexTokenId,
	}
	pfe.farms[2] = Farm{
		address:      pfe.mexFarmAddress,
		farmToken:    pfe.mexFarmTokenId,
		farmingToken: pfe.mexTokenId,
		rewardToken:  pfe.mexTokenId,
	}

	pfe.swaps[0] = SwapPair{
		address:     pfe.wemeSwapAddress,
		lpToken:     pfe.wemeLpTokenId,
		firstToken:  pfe.wegldTokenId,
		secondToken: pfe.mexTokenId,
	}
	pfe.swaps[1] = SwapPair{
		address:     pfe.webuSwapAddress,
		lpToken:     pfe.webuLpTokenId,
		firstToken:  pfe.wegldTokenId,
		secondToken: pfe.busdTokenId,
	}

	// users
	esdtString := pfe.fullOfEsdtWalletString()
	for i := 1; i <= args.numUsers; i++ {
		err := pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"storage": {},
					"esdt": {
						%s
					},
					"code": ""
				}
			}
		}`,
			pfe.userAddress(i),
			esdtString,
		))
		if err != nil {
			return err
		}
	}
	err := pfe.executeStep(fmt.Sprintf(`
	{
		"step": "setState",
		"accounts": {
			"%s": {
				"nonce": "0",
				"balance": "1,000,000,000,000,000,000,000,000,000,000",
				"storage": {},
				"code": ""
			}
		}
	}`,
		pfe.ownerAddress,
	))
	if err != nil {
		return err
	}

	// swaps
	err = pfe.setupPair(pfe.wemeSwapAddress, pfe.wegldTokenId, pfe.mexTokenId, pfe.wemeLpTokenId, pfe.ownerAddress, pfe.ownerAddress)
	if err != nil {
		return err
	}

	err = pfe.setupPair(pfe.webuSwapAddress, pfe.wegldTokenId, pfe.busdTokenId, pfe.webuLpTokenId, pfe.ownerAddress, pfe.ownerAddress)
	if err != nil {
		return err
	}

	// farms
	err = pfe.setupFarm(pfe.wemeFarmAddress, pfe.wemeFarmTokenId, pfe.wemeLpTokenId, pfe.mexTokenId, pfe.ownerAddress, pfe.ownerAddress)
	if err != nil {
		return err
	}

	err = pfe.setupFarm(pfe.webuFarmAddress, pfe.webuFarmTokenId, pfe.webuLpTokenId, pfe.mexTokenId, pfe.ownerAddress, pfe.ownerAddress)
	if err != nil {
		return err
	}

	err = pfe.setupFarm(pfe.mexFarmAddress, pfe.mexFarmTokenId, pfe.mexTokenId, pfe.mexTokenId, pfe.ownerAddress, pfe.ownerAddress)
	if err != nil {
		return err
	}

	// configs
	err = pfe.setFeeOn(pfe.wemeSwapAddress, pfe.wemeFarmAddress, pfe.mexTokenId, pfe.ownerAddress)
	if err != nil {
		return err
	}

	err = pfe.setFeeOn(pfe.webuSwapAddress, pfe.webuFarmAddress, pfe.mexTokenId, pfe.ownerAddress)
	if err != nil {
		return err
	}

	err = pfe.whitelist(pfe.ownerAddress, pfe.wemeSwapAddress, pfe.webuSwapAddress)
	if err != nil {
		return err
	}

	err = pfe.addTrustedSwapPair(pfe.ownerAddress, pfe.webuSwapAddress, pfe.wemeSwapAddress, pfe.wegldTokenId, pfe.mexTokenId)
	if err != nil {
		return err
	}

	pfe.log("init ok")
	return nil
}

func (pfe *fuzzDexExecutor) setupPair(swapAddress, firstTokenId, secondTokenId, lpTokenId, routerAddress, owner string) error {
	return pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"esdt": {
						"str:%s": {
							"roles": [
								"ESDTRoleLocalMint",
								"ESDTRoleLocalBurn"
							]
						}
					},
					"storage": {
						"str:first_token_id": "str:%s",
						"str:second_token_id": "str:%s",
						"str:state": "1",
						"str:lpTokenIdentifier": "str:%s",
						"str:router_address": "%s",
						"str:router_owner_address": "%s",
						"str:total_fee_percent": "300",
						"str:special_fee_percent": "050",
						"str:burn_tokens_gas_limit": "500000",
						"str:mint_tokens_gas_limit": "500000",
						"str:send_fee_gas_limit": "25,000,000",
						"str:extern_swap_gas_limit": "50,000,000"
					},
					"code": "file:elrond_dex_pair.wasm",
					"owner": "%s"
				}
			}
		}`,
		swapAddress,
		lpTokenId,
		firstTokenId,
		secondTokenId,
		lpTokenId,
		routerAddress,
		owner,
		owner,
	))
}

func (pfe *fuzzDexExecutor) setupFarm(farmAddress, farmTokenId, enterFarmTokenId, rewardTokenId, routerAddress, owner string) error {
	return pfe.executeStep(fmt.Sprintf(`
		{
			"step": "setState",
			"accounts": {
				"%s": {
					"nonce": "0",
					"balance": "0",
					"esdt": {
						"str:%s": {
							"roles": [
								"ESDTRoleNFTCreate",
								"ESDTRoleNFTAddQuantity",
								"ESDTRoleNFTBurn"
							]
						},
						"str:%s": {
							"roles": [
								"ESDTRoleLocalMint",
								"ESDTRoleLocalBurn"
							]
						},
						"str:%s": {
							"roles": [
								"ESDTRoleLocalMint",
								"ESDTRoleLocalBurn"
							]
						}
					},
					"storage": {
						"str:farming_token_id": "str:%s",
						"str:farm_token_id": "str:%s",
						"str:reward_token_id": "str:%s",
						"str:router_address": "%s",
						"str:state": "1",
						"str:owner": "%s",
						"str:minimum_farming_epochs": "2",
						"str:burn_tokens_gas_limit": "5,000,000",
						"str:mint_tokens_gas_limit": "5,000,000",
						"str:locked_rewards_apr_multiplier": "2",
						"str:division_safety_constant": "1000000000000",
						"str:create_farm_tokens_gas_limit": "5000000",
						"str:produce_rewards_enabled": "1",
						"str:per_block_reward_amount": "10000000000000000",
						"str:penalty_percent": "0",
						"str:nft_deposit_max_len": "10",
						"str:nft_deposit_accepted_token_ids.info": "u32:1|u32:1|u32:1|u32:1",
						"str:nft_deposit_accepted_token_ids.node_id|str:%s": "1",
						"str:nft_deposit_accepted_token_ids.node_links": "0",
						"str:nft_deposit_accepted_token_ids.value|u32:1": "str:%s"
					},
					"code": "file:elrond_dex_farm.wasm",
					"owner": "%s"
				}
			}
		}`,
		farmAddress,
		farmTokenId,
		enterFarmTokenId,
		rewardTokenId,
		enterFarmTokenId,
		farmTokenId,
		rewardTokenId,
		routerAddress,
		owner,
		farmTokenId,
		farmTokenId,
		owner,
	))
}

func (pfe *fuzzDexExecutor) setFeeOn(swapAddress, farmAddress, feeToken, ownerAddress string) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "set-fee-on",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "setFeeOn",
				"arguments": [
					"1",
					"%s",
					"str:%s"
				],
				"gasLimit": "10,000,000",
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
		ownerAddress,
		swapAddress,
		farmAddress,
		feeToken,
	))
	return err
}

func (pfe *fuzzDexExecutor) whitelist(ownerAddress, swapAddressToConfig, swapAddressToWhitelist string) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "whitelist",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "whitelist",
				"arguments": [
					"%s"
				],
				"gasLimit": "10,000,000",
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
		ownerAddress,
		swapAddressToConfig,
		swapAddressToWhitelist,
	))
	return err
}

func (pfe *fuzzDexExecutor) addTrustedSwapPair(ownerAddress, swapAddressToConfig, swapAddressToAdd, firstTokenId, secondTokenId string) error {
	_, err := pfe.executeTxStep(fmt.Sprintf(`
		{
			"step": "scCall",
			"txId": "whitelist",
			"tx": {
				"from": "%s",
				"to": "%s",
				"value": "0",
				"function": "addTrustedSwapPair",
				"arguments": [
					"%s",
					"str:%s",
					"str:%s"
				],
				"gasLimit": "10,000,000",
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
		ownerAddress,
		swapAddressToConfig,
		swapAddressToAdd,
		firstTokenId,
		secondTokenId,
	))
	return err
}
