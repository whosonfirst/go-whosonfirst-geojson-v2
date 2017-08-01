package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
)

// type Hash interface {
//     HashFile(string) (string, error)
//     HashBytes([]byte) (string, error)
//     HashFromJSON([]byte) (string, error)
// }

type Hash struct {
	algo string
}

func DefaultHash() (*Hash, error) {
	return NewHash("md5")
}

func NewHash(algo string) (*Hash, error) {

	switch algo {

	case "md5":
		// pass
	default:
		return nil, errors.New("Unsupported hashing algorithm")
	}

	h := Hash{
		algo: algo,
	}

	return &h, nil
}

func (h *Hash) HashFile(path string) (string, error) {

	body, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	return h.HashBytes(body)
}

func (h *Hash) HashFromJSON(raw []byte) (string, error) {

	var geom interface{}

	err := json.Unmarshal(raw, &geom)

	if err != nil {
		return "", err
	}

	body, err := json.Marshal(geom)

	if err != nil {
		return "", err
	}

	return h.HashBytes(body)
}

func (h *Hash) HashBytes(body []byte) (string, error) {

	var hash [16]byte

	switch h.algo {

	case "md5":
		hash = md5.Sum(body)
	default:
		return "", errors.New("How did we even get this far")
	}

	hex := hex.EncodeToString(hash[:])
	return hex, nil
}
