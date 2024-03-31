package merkledag

import (
	"encoding/json"
	"hash"
)

type Link struct {
	Name string
	Hash []byte
	Size int
}

type Object struct {
	Links []Link
	Data  []byte
}

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 如果是文件,直接存储文件内容
	if node.Type() == FILE {
		file, _ := node.(File)
		data := file.Bytes()
		h.Write(data)
		store.Put(h.Sum(nil), data)
		return h.Sum(nil)
	}

	// 如果是目录,则递归存储子节点
	dir, _ := node.(Dir)
	var links []Link
	for it := dir.It(); it.Next(); {
		node := it.Node()
		hash := Add(store, node, h)
		size := node.Size()
		links = append(links, Link{
			Name: node.Name(),
			Hash: hash,
			Size: int(size),
		})
	}

	obj := Object{
		Links: links,
	}
	data, _ := json.Marshal(obj)
	h.Write(data)
	objectHash := h.Sum(nil)
	store.Put(objectHash, data)
	return objectHash
}
