package sync

import (
	"vicarnet/db"
)

func sendCharSync(id string, data string) {
	db.Cache.Set("sync_char_"+id, data, nil)
}

func sendCharLevelSync(id string, data string) {
	db.Cache.Set("sync_char_level_"+id, data, nil)
}

func retrieveCharSync(ids []string) map[string]map[string]string {
	result := make(map[string]map[string]string)

	for _, id := range ids {
		data, ok := getCharSyncData(id)
		if !ok {
			continue
		}

		result[id] = data
	}

	return result
}

func getCharSyncData(id string) (map[string]string, bool) {
	data := make(map[string]string)
	success := false

	charData, ok := db.Cache.GetOnce("sync_char_" + id)
	if ok {
		data["c"] = charData.(string)
		success = true
	}

	levelData, ok := db.Cache.GetOnce("sync_char_level_" + id)
	if ok {
		data["l"] = levelData.(string)
		success = true
	}

	return data, success

}
