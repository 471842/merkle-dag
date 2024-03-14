package merkledag

import (
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
	default:
		return nil
	}
}

func StoreFile(store KVStore, file File, h hash.Hash) []byte {
	data := file.Bytes()
	if file.Size() <= BLOCK_SIZE {
		h.Reset()
		h.Write(data)
		hashSum := h.Sum(nil)

		if err := store.Put(hashSum, data); err != nil {
			panic(err)
		}

		return hashSum
	} else {
		t := []byte("blob")
		if file.Size() > BLOCK_SIZE {
			t = []byte("list")
		}
		return t
	}
}

func StoreDir(store KVStore, dir Dir, h hash.Hash) []byte {
	it := dir.It()
	var links []Link

	for it.Next() {
		node := it.Node()
		link := Link{
			Hash: Add(store, node, h),
			Size: int(node.Size()),
		}
		links = append(links, link)
	}

	obj := Object{
		Links: links,
	}

	// Define tree and store the object in KVStore
	hashSum := calculateHash(h, obj)
	if err := store.Put(hashSum, encodeObject(obj)); err != nil {
		panic(err)
	}

	return hashSum
}

func calculateHash(h hash.Hash, obj Object) []byte {
	h.Reset()
	for _, link := range obj.Links {
		h.Write(link.Hash)
	}
	return h.Sum(nil)
}

func encodeObject(obj Object) []byte {
	// Encode object as needed for storage in KVStore
	// This is a placeholder function and needs to be implemented based on the actual requirements
	return []byte{}
}
