package constant

type ErrorCode int

// https://github.com/DarkEnergyProcessor/NPPS4/blob/master/npps4/idol/error.py
const (
	ErrorCodeUnknown                 ErrorCode = 1
	ErrorCodeNoUnitIsSellable        ErrorCode = 1205
	ErrorCodeLivePreciseListNotFound ErrorCode = 3421
	ErrorCodeEventNoEventData        ErrorCode = 10004
)
