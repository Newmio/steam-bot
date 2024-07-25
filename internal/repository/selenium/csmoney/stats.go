package repocsmoney

import (
	"bot/internal/domain/entity"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/Newmio/steam_helper"
)

func (r *csmoney) GetRareFloats(limit, offset int) (map[string][]entity.FloatItem, error) {
	items := make(map[string][]entity.FloatItem)

	param := steam_helper.Param{
		Method:  "GET",
		Url:     "https://cs.money/5.0/load_bots_inventory/730?hasRareFloat=true&limit=60&offset=0&order=asc&priceWithBonus=0&sort=price&withStack=true",
		Headers: map[string]string{"Accept": "application/json"},
	}

	resp, err := r.http.Do(param)
	if err != nil {
		return nil, steam_helper.Trace(err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, steam_helper.Trace(err)
	}

	if _, ok := data["items"]; !ok || resp.StatusCode != 200 {
		return nil, steam_helper.Trace(err, fmt.Sprintf("status %d: %s", resp.StatusCode, string(resp.Body)))
	}

	for _, value := range data["items"].([]interface{}) {
		itemMap := value.(map[string]interface{})

		float, err := strconv.ParseFloat(itemMap["float"].(string), 64)
		if err != nil {
			return nil, steam_helper.Trace(err)
		}

		hashName := url.QueryEscape(itemMap["fullName"].(string))

		items[hashName] = append(items[hashName], entity.FloatItem{
			Float:    float,
			HashName: hashName,
			Cost:     int(itemMap["price"].(float64) * 100),
		})
	}

	return items, nil
}
