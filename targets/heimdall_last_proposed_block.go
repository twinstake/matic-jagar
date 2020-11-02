package targets

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetLatestProposedBlockAndTime to get latest proposed block height and time
func GetLatestProposedBlockAndTime(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var blockResp LatestBlock
	err = json.Unmarshal(resp.Body, &blockResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	time := blockResp.Block.Header.Time
	blockTime := GetUserDateFormat(time)
	blockHeight := blockResp.Block.Header.Height
	log.Printf("last proposed block time : %s,  height : %s", blockTime, blockHeight)

	if cfg.ValDetails.ValidatorHexAddress == blockResp.Block.Header.ProposerAddress {
		fields := map[string]interface{}{
			"height":     blockHeight,
			"block_time": blockTime,
		}
		_ = writeToInfluxDb(c, bp, "heimdall_last_proposed_block", map[string]string{}, fields)
	}

	_ = writeToInfluxDb(c, bp, "heimdall_lastest_block", map[string]string{}, map[string]interface{}{"height": blockHeight, "block_time": time})

	// Store chainID in database
	chainID := blockResp.Block.Header.ChainID
	_ = writeToInfluxDb(c, bp, "heimdall_chain_id", map[string]string{}, map[string]interface{}{"chain_id": chainID})
	log.Printf("Chain ID : %s ", chainID)
}

// GetPrevBlockTime returns time of the pevious block
func GetPrevBlockTime(cfg *config.Config, c client.Client, height string) string {
	var t string
	q := client.NewQuery(fmt.Sprintf("SELECT last(block_time) FROM heimdall_lastest_block WHERE height = '%s'", height), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						t = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return t
}
