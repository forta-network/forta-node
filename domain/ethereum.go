package domain

//TraceAction is an element of a trace_block Trace response
type TraceAction struct {
	CallType      *string `json:"callType"`
	To            *string `json:"to"`
	Input         *string `json:"input"`
	From          *string `json:"from"`
	Gas           *string `json:"gas"`
	Value         *string `json:"value"`
	Init          *string `json:"init"`
	Address       *string `json:"address"`
	Balance       *string `json:"balance"`
	RefundAddress *string `json:"refundAddress"`
}

//TraceResult is a result element of a trace_block Trace response
type TraceResult struct {
	Output  *string `json:"output"`
	GasUsed *string `json:"gasUsed"`
	Address *string `json:"address"`
	Code    *string `json:"code"`
}

//Trace is a
type Trace struct {
	Action              TraceAction  `json:"action"`
	BlockHash           *string      `json:"blockHash"`
	BlockNumber         *int         `json:"blockNumber"`
	Result              *TraceResult `json:"result"`
	Subtraces           int          `json:"subtraces"`
	TraceAddress        []int        `json:"traceAddress"`
	TransactionHash     *string      `json:"transactionHash"`
	TransactionPosition *int         `json:"transactionPosition"`
	Type                string       `json:"type"`
	Error               *string      `json:"error"`
}
