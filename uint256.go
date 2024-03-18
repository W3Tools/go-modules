package gm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
)

// Uint256 is like `u256` in move.
type Uint256 struct {
	lo uint64
	hi uint64
}

var (
	_ json.Marshaler   = (*Uint256)(nil)
	_ json.Unmarshaler = (*Uint256)(nil)
	_ Marshaler        = (*Uint256)(nil)
)

func (u Uint256) Big() *big.Int {
	loBig := NewBigIntFromUint64(u.lo)
	hiBig := NewBigIntFromUint64(u.hi)
	hiBig = hiBig.Lsh(hiBig, 64)

	return hiBig.Add(hiBig, loBig)
}

func (i Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Big().String())
}

var maxU256 = (&big.Int{}).Lsh(big.NewInt(1), 256)

func checkUint256(bigI *big.Int) error {
	if bigI.Sign() < 0 {
		return fmt.Errorf("%s is negative", bigI.String())
	}

	if bigI.Cmp(maxU256) >= 0 {
		return fmt.Errorf("%s is greater than Max Uint 256", bigI.String())
	}

	return nil
}

func (u *Uint256) SetBigInt(bigI *big.Int) error {
	if err := checkUint256(bigI); err != nil {
		return err
	}

	r := make([]byte, 0, 16)
	bs := bigI.Bytes()
	for i := 0; i+len(bs) < 16; i++ {
		r = append(r, 0)
	}
	r = append(r, bs...)

	hi := binary.BigEndian.Uint64(r[0:8])
	lo := binary.BigEndian.Uint64(r[8:16])

	u.hi = hi
	u.lo = lo

	return nil
}

func (u *Uint256) UnmarshalText(data []byte) error {
	bigI := &big.Int{}
	_, ok := bigI.SetString(string(data), 10)
	if !ok {
		return fmt.Errorf("failed to parse %s as an integer", string(data))
	}

	return u.SetBigInt(bigI)
}

func (i *Uint256) UnmarshalJSON(data []byte) error {
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	bigI := &big.Int{}
	_, ok := bigI.SetString(dataStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse %s as an integer", dataStr)
	}

	return i.SetBigInt(bigI)
}

func NewUint256FromBigInt(bigI *big.Int) (*Uint256, error) {
	i := &Uint256{}

	if err := i.SetBigInt(bigI); err != nil {
		return nil, err
	}

	return i, nil
}

func NewUint256(s string) (*Uint256, error) {
	r := &big.Int{}
	r, ok := r.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s as an integer", s)
	}

	return NewUint256FromBigInt(r)
}

func (i Uint256) MarshalBCS() ([]byte, error) {
	r := make([]byte, 16)

	binary.LittleEndian.PutUint64(r, i.lo)
	binary.LittleEndian.PutUint64(r[8:], i.hi)

	return r, nil
}

func (i *Uint256) Cmp(j *Uint256) int {
	switch {
	case i.hi > j.hi || (i.hi == j.hi && i.lo > j.lo):
		return 1
	case i.hi == j.hi && i.lo == j.lo:
		return 0
	default:
		return -1
	}
}

func NewUint256FromUint64(lo, hi uint64) *Uint256 {
	return &Uint256{
		lo: lo,
		hi: hi,
	}
}

func (u Uint256) String() string {
	return u.Big().String()
}
