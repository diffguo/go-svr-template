package common

// 定义10000以上的ErrCode和其含义。 客户端自己使用的范围为 700~10000， 700以下兼容http状态码

const (
	ErrCodeParamErr = iota + 10000
	ErrCodeDBErr
	ErrCodeLogicErr
)

var MapErrCode2Desc = map[int]string{
	ErrCodeParamErr: "参数错误",
	ErrCodeDBErr:    "数据库错误",
	ErrCodeLogicErr: "", // 逻辑错误, 请使用SendResponseImp填充具体的错误原因
}
