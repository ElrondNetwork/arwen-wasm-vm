package gasschedules

import (
	"fmt"

	"github.com/ElrondNetwork/arwen-wasm-vm/config"
	"github.com/pelletier/go-toml"
)

// LoadGasScheduleConfig parses and prepares a gas schedule read from file.
func LoadGasScheduleConfig(fileContents string) (config.GasScheduleMap, error) {
	loadedTree, err := toml.Load(fileContents)
	if err != nil {
		fmt.Printf("cannot interpret file contents as toml: %s", err.Error())
		return nil, err
	}

	gasScheduleConfig := loadedTree.ToMap()

	flattenedGasSchedule := make(config.GasScheduleMap)
	for libType, costs := range gasScheduleConfig {
		flattenedGasSchedule[libType] = make(map[string]uint64)
		costsMap := costs.(map[string]interface{})
		for operationName, cost := range costsMap {
			flattenedGasSchedule[libType][operationName] = uint64(cost.(int64))
		}
	}

	return flattenedGasSchedule, nil
}
