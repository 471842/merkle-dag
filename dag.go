package merkledag

import (
	"encoding/json"
	"hash"
)

const (
	K          = 1 << 10
	BLOCK_SIZE = 256 * K
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
	switch node.Type() {
	case FILE:
		return StoreFile(store, node.(File), h)
	case DIR:
		return StoreDir(store, node.(Dir), h)
	}
	return nil
}

func StoreFile(store KVStore, file File, h hash.Hash) []byte {
	t := []byte("blob")
	if file.Size() > BLOCK_SIZE {
		t = []byte("list")
	}

	data := file.Bytes()
	if len(data) <= BLOCK_SIZE {
		// 文件小于等于BLOCK_SIZE，直接存储
		h.Reset()
		h.Write(t)
		h.Write(data)
		hash := h.Sum(nil)
		store.Put(hash, data)
		return hash
	} else {
		// 文件大于BLOCK_SIZE，分块存储
		var links []Link
		for i := 0; i < len(data); i += BLOCK_SIZE {
			end := i + BLOCK_SIZE
			if end > len(data) {
				end = len(data)
			}
			block := data[i:end]
			h.Reset()
			h.Write([]byte("blob"))
			h.Write(block)
			hash := h.Sum(nil)
			store.Put(hash, block)
			links = append(links, Link{Hash: hash, Size: len(block)})
		}
		// 存储分块链接
		linksData, _ := json.Marshal(links)
		h.Reset()
		h.Write(t)
		h.Write(linksData)
		hash := h.Sum(nil)
		store.Put(hash, linksData)
		return hash
	}
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) []byte {
	var tree Object
	tree.Links = make([]Link, 0)
	tree.Data = make([]byte, 0)

	it := dir.It()
	for it.Next() {
		node := it.Node()
		hash := Add(store, node, h)
		link := Link{Name: node.Name(), Hash: hash, Size: int(node.Size())}
		tree.Links = append(tree.Links, link)
	}

	treeData, _ := json.Marshal(tree)
	h.Reset()
	h.Write([]byte("tree"))
	h.Write(treeData)
	hash := h.Sum(nil)
	store.Put(hash, treeData)
	return hash
}
