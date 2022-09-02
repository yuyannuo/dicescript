package dicescript

import (
	"errors"
	"strconv"
)

type ParserData struct {
	counterStack []int64  // f-string 嵌套计数，在解析时中起作用
	varnameStack []string // 另一个解析用栈
}

func (pd *ParserData) init() {
	pd.counterStack = []int64{}
	pd.varnameStack = []string{}
}

func (e *Parser) checkStackOverflow() bool {
	if e.Error != nil {
		return true
	}
	if e.codeIndex >= len(e.code) {
		need := len(e.code) * 2
		if need <= 8192 {
			newCode := make([]ByteCode, need)
			copy(newCode, e.code)
			e.code = newCode
		} else {
			e.Error = errors.New("E1:指令虚拟机栈溢出，请不要发送过长的指令")
			return true
		}
	}
	return false
}

func (e *Parser) WriteCode(T CodeType, value interface{}) {
	if e.checkStackOverflow() {
		return
	}

	c := &e.code[e.codeIndex]
	c.T = T
	c.Value = value
	e.codeIndex += 1
}

func (e *Parser) LMark() {
	e.WriteCode(TypeLeftValueMark, nil)
}

func (e *Parser) AddOp(operator CodeType) {
	e.WriteCode(operator, nil)
}

func (e *Parser) AddLoadName(value string) {
	e.WriteCode(TypeLoadName, value)
}

func (e *Parser) PushIntNumber(value string) {
	val, _ := strconv.ParseInt(value, 10, 64)
	e.WriteCode(TypePushIntNumber, int64(val))
}

func (e *Parser) PushStr(value string) {
	e.WriteCode(TypePushString, value)
}

func (e *Parser) AddFormatString(value string, num int64) {
	//e.PushStr(value)
	e.WriteCode(TypeLoadFormatString, num) // num
}

func (e *Parser) PushFloatNumber(value string) {
	val, _ := strconv.ParseFloat(value, 64)
	e.WriteCode(TypePushFloatNumber, float64(val))
}

func (e *Parser) AddStore(text string) {
	e.WriteCode(TypeStoreName, text)
}

func (e *Parser) NamePush(test string) {
	e.varnameStack = append(e.varnameStack, test)
}

func (e *Parser) NamePop() string {
	last := len(e.varnameStack) - 1
	val := e.varnameStack[last]
	e.varnameStack = e.varnameStack[:last]
	return val
}

func (e *Parser) CounterPush() {
	e.counterStack = append(e.counterStack, 0)
}

func (e *Parser) CounterAdd(offset int64) {
	last := len(e.counterStack) - 1
	if last != -1 {
		e.counterStack[last] += offset
	}
}

func (e *Parser) CounterPop() int64 {
	last := len(e.counterStack) - 1
	num := e.counterStack[last]
	e.counterStack = e.counterStack[:last]
	return num
}
