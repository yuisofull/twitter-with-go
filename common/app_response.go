package common

type successRes struct {
	Data   interface{} `json:"data"`
	Paging interface{} `json:"paging,omitempty" form:"paging,omitempty"`
	Filter interface{} `json:"filter,omitempty" form:"filter,omitempty"`
}

type simpleSuccessRes struct {
	Data interface{} `json:"data"`
}

func NewSuccessResponse(data, paging, filter interface{}) *successRes {
	return &successRes{Data: data, Paging: paging, Filter: filter}
}

func SimpleNewSuccessResponse(data interface{}) *simpleSuccessRes {
	return &simpleSuccessRes{Data: data}
}
