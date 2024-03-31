package merkledag

import (
	"bytes"
	"encoding/json"
)

// Hash to file
func Hash2File(store KVStore, hash []byte, path string, hp HashPool) []byte {
	data, _ := store.Get(hash)
	if data == nil {
		return nil
	}

	var obj Object
	json.Unmarshal(data, &obj)

	if len(obj.Links) == 0 {
		return obj.Data
	}

	var buffer bytes.Buffer
	for _, link := range obj.Links {
		childData := Hash2File(store, link.Hash, path+"/"+link.Name, hp)
		if childData != nil {
			buffer.Write(childData)
		}
	}

	return buffer.Bytes()
}
