package types

type ProtocolConfigValue interface {
	isProtocolConfigValue()
}

type ProtocolConfigValue_u32 struct {
	U32 string `json:"u32"`
}

type ProtocolConfigValue_u64 struct {
	U64 string `json:"u64"`
}

type ProtocolConfigValue_f64 struct {
	F64 string `json:"f64"`
}

func (ProtocolConfigValue_u32) isProtocolConfigValue() {}
func (ProtocolConfigValue_u64) isProtocolConfigValue() {}
func (ProtocolConfigValue_f64) isProtocolConfigValue() {}
