package gmsui

type SuiMoveId struct {
	Id string `json:"id"`
}

type SuiMoveTable struct {
	Fields SuiMoveTableFields `json:"fields"`
	Type   string             `json:"type"`
}

type SuiMoveTableFields struct {
	Id   SuiMoveId `json:"id"`
	Size string    `json:"size"`
}

type SuiMoveString struct {
	Type   string              `json:"type"`
	Fields SuiMoveStringFields `json:"fields"`
}

type SuiMoveStringFields struct {
	Name string `json:"name"`
}

type SuiMoveDynamicField[TypeFields any, TypeName any] struct {
	Id    SuiMoveId                            `json:"id"`
	Name  TypeName                             `json:"name"`
	Value SuiMoveDynamicFieldValue[TypeFields] `json:"value"`
}

type SuiMoveDynamicFieldValue[TypeFields any] struct {
	Fields TypeFields `json:"fields"`
	Type   string     `json:"type"`
}
