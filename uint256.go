package gm

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
)

// Uint256 is like `u256` in move.
type Uint256 struct {
	lo Uint128
	hi Uint128
}

var (
	_ json.Marshaler   = (*Uint256)(nil)
	_ json.Unmarshaler = (*Uint256)(nil)
	_ Marshaler        = (*Uint256)(nil)
)

func (u Uint256) Big() *big.Int {
	loBig := u.lo.Big()
	hiBig := u.hi.Big()
	hiBig = hiBig.Lsh(hiBig, 128)

	return hiBig.Add(hiBig, loBig)
}

func (u Uint256) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Big().String())
}

func checkUint256(bigI *big.Int) error {
	if bigI.Sign() < 0 {
		return fmt.Errorf("%s is negative", bigI.String())
	}

	if bigI.BitLen() > 256 {
		return fmt.Errorf("%s is greater than Max Uint 256", bigI.String())
	}

	return nil
}

func (u *Uint256) SetBigInt(bigI *big.Int) error {
	if err := checkUint256(bigI); err != nil {
		return err
	}

	r := make([]byte, 0, 32)
	bs := bigI.Bytes()
	for i := 0; i+len(bs) < 32; i++ {
		r = append(r, 0)
	}
	r = append(r, bs...)

	lo := Uint128{
		lo: binary.BigEndian.Uint64(r[24:32]),
		hi: binary.BigEndian.Uint64(r[16:24]),
	}
	hi := Uint128{
		lo: binary.BigEndian.Uint64(r[8:16]),
		hi: binary.BigEndian.Uint64(r[0:8]),
	}

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

func (u *Uint256) UnmarshalJSON(data []byte) error {
	var dataStr string
	if err := json.Unmarshal(data, &dataStr); err != nil {
		return err
	}

	bigI := &big.Int{}
	_, ok := bigI.SetString(dataStr, 10)
	if !ok {
		return fmt.Errorf("failed to parse %s as an integer", dataStr)
	}

	return u.SetBigInt(bigI)
}

func NewUint256FromBigInt(bigI *big.Int) (*Uint256, error) {
	u := &Uint256{}

	if err := u.SetBigInt(bigI); err != nil {
		return nil, err
	}

	return u, nil
}

func NewUint256(s string) (*Uint256, error) {
	b := &big.Int{}
	b, ok := b.SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse %s as an integer", s)
	}

	return NewUint256FromBigInt(b)
}

func (u Uint256) MarshalBCS() ([]byte, error) {
	r := make([]byte, 32)

	binary.LittleEndian.PutUint64(r, u.hi.lo)
	binary.LittleEndian.PutUint64(r[8:], u.hi.hi)
	binary.LittleEndian.PutUint64(r[16:], u.lo.lo)
	binary.LittleEndian.PutUint64(r[32:], u.lo.hi)

	return r, nil
}

func (i *Uint256) Cmp(j *Uint256) int {
	switch {
	case i.hi.Cmp(&j.hi) > 0 || (i.hi.Cmp(&j.hi) == 0 && i.lo.Cmp(&j.lo) > 0):
		return 1
	case i.hi.Cmp(&j.hi) == 0 && i.lo.Cmp(&j.lo) == 0:
		return 0
	default:
		return -1
	}
}

func (u Uint256) String() string {
	return u.Big().String()
}
