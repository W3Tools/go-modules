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

type SuiMoveDynamicField[T any] struct {
	Id    SuiMoveId                   `json:"id"`
	Name  string                      `json:"name"`
	Value SuiMoveDynamicFieldValue[T] `json:"value"`
}

type SuiMoveDynamicFieldValue[T any] struct {
	Fields T      `json:"fields"`
	Type   string `json:"type"`
}
