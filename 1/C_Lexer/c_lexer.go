package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"unicode"
)

// TokenType 标记类型
// 这里使用字符串来表示不同类型的标记，以便于识别和处理各种词法单元
type TokenType string

// Token 标记结构
// 用于存储词法分析过程中识别出的每个标记的类型和对应的值
type Token struct {
	Type  TokenType
	Value string
}

// 定义了一系列常量，表示不同的标记类型，方便在词法分析过程中进行分类和判断
const (
	ILLEGAL        = "ILLEGAL"        // 非法标记，用于表示在词法分析中遇到不符合语法规则的情况
	EOF            = "EOF"            // 文件结束，当读取到输入源代码的末尾时，返回此标记类型
	IDENT          = "IDENT"          // 标识符，用于表示变量名、函数名等自定义的名称
	INT            = "INT"            // 整数，用于表示整数类型的数值
	FLOAT          = "FLOAT"          // 浮点数，用于表示带有小数部分的数值
	STRING         = "STRING"         // 字符串字面量，用于表示用双引号括起来的字符序列
	RETURN         = "RETURN"         // 返回标志，用于表示C语言中的return关键字
	PLUS           = "PLUS"           // 加号运算符
	MINUS          = "MINUS"          // 减号运算符
	STAR           = "STAR"           // 乘号运算符
	SLASH          = "SLASH"          // 除号运算符
	EQUAL          = "EQUAL"          // 等号运算符
	LESS           = "LESS"           // 小于号运算符
	GREATER        = "GREATER"        // 大于号运算符
	SEMICOLON      = "SEMICOLON"      // 分号
	COMMA          = "COMMA"          // 逗号
	LPAREN         = "LPAREN"         // 左括号
	RPAREN         = "RPAREN"         // 右括号
	LBRACE         = "LBRACE"         // 左大括号
	RBRACE         = "RBRACE"         // 右大括号
	PREPROC        = "PREPROC"        // 预处理指令标记
	VOID           = "VOID"           // void关键字标记
	SIGNED         = "SIGNED"         // signed关键字标记
	UNSIGNED       = "UNSIGNED"       // unsigned关键字标记
	SHORT          = "SHORT"          // short关键字标记
	LONG           = "LONG"           // long关键字标记
	DOUBLE         = "DOUBLE"         // double关键字标记
	CHAR           = "CHAR"           // char关键字标记
	ENUM           = "ENUM"           // enum关键字标记
	STRUCT         = "STRUCT"         // struct关键字标记
	UNION          = "UNION"          // union关键字标记
	TYPEDEF        = "TYPEDEF"        // typedef关键字标记
	CONST          = "CONST"          // const关键字标记
	VOLATILE       = "VOLATILE"       // volatile关键字标记
	AUTO           = "AUTO"           // auto关键字标记
	STATIC         = "STATIC"         // static关键字标记
	EXTERN         = "EXTERN"         // extern关键字标记
	REGISTER       = "REGISTER"       // register关键字标记
	SIZEOF         = "SIZEOF"         // sizeof关键字标记
	GOTO           = "GOTO"           // goto关键字标记
	BREAK          = "BREAK"          // break关键字标记
	CONTINUE       = "CONTINUE"       // continue关键字标记
	IF             = "IF"             // if关键字标记
	ELSE           = "ELSE"           // else关键字标记
	SWITCH         = "SWITCH"         // switch关键字标记
	CASE           = "CASE"           // case关键字标记
	DEFAULT        = "DEFAULT"        // default关键字标记
	DO             = "DO"             // do关键字标记
	WHILE          = "WHILE"          // while关键字标记
	FOR            = "FOR"            // for关键字标记
	AND            = "AND"            // &运算符标记（逻辑与、按位与等情况，这里统一先简单标记为AND，可根据后续需求细化）
	OR             = "OR"             // |运算符标记（逻辑或、按位或等情况，这里统一先简单标记为OR，可根据后续需求细化）
	XOR            = "XOR"            // ^运算符标记（逻辑异或、按位异或等情况，这里统一先简单标记为XOR，可根据后续需求细化）
	NOT            = "NOT"            // ~运算符标记（逻辑非、按位取反等情况，这里统一先简单标记为NOT，可根据后续需求细化）
	LSHIFT         = "LSHIFT"         // <<运算符标记（左移运算符）
	RSHIFT         = "RSHIFT"         // >>运算符标记（右移运算符）
	COLON          = "COLON"          // :运算符标记
	ANDAND         = "ANDAND"         // &&运算符标记（逻辑与）
	OROR           = "OROR"           // ||运算符标记（逻辑或）
	NOTNOT         = "NOTNOT"         //!运算符标记（逻辑非）
	EQEQ           = "EQEQ"           // ==运算符标记（等于比较）
	LT             = "LT"             // <运算符标记（小于）
	GT             = "GT"             // >运算符标记（大于）
	LTEQ           = "LTEQ"           // <=运算符标记（小于等于）
	GTEQ           = "GTEQ"           // >=运算符标记（大于等于）
	COMMENT_SINGLE = "COMMENT_SINGLE" // 单行注释标记
	COMMENT_MULTI  = "COMMENT_MULTI"  // 多行注释标记
)

// 创建一个映射表，将C语言中的关键字映射到对应的标记类型，以便快速判断一个标识符是否为关键字
var keywords map[string]TokenType = map[string]TokenType{
	"int":      INT,
	"float":    FLOAT,
	"return":   RETURN,
	"include":  PREPROC,
	"void":     VOID,
	"signed":   SIGNED,
	"unsigned": UNSIGNED,
	"short":    SHORT,
	"long":     LONG,
	"double":   DOUBLE,
	"char":     CHAR,
	"enum":     ENUM,
	"struct":   STRUCT,
	"union":    UNION,
	"typedef":  TYPEDEF,
	"const":    CONST,
	"volatile": VOLATILE,
	"auto":     AUTO,
	"static":   STATIC,
	"extern":   EXTERN,
	"register": REGISTER,
	"sizeof":   SIZEOF,
	"goto":     GOTO,
	"break":    BREAK,
	"continue": CONTINUE,
	"if":       IF,
	"else":     ELSE,
	"switch":   SWITCH,
	"case":     CASE,
	"default":  DEFAULT,
	"do":       DO,
	"while":    WHILE,
	"for":      FOR,
}

func main() {
	furl := `test.c` //输入文件的路径
	file, err := os.Open(furl)
	if err != nil {
		fmt.Println("未找到对应文件" + furl)
		// 处理打开文件失败的情况
	}
	defer file.Close()

	var buf bytes.Buffer
	reader := bufio.NewReader(file)
	_, err = io.Copy(&buf, reader)
	if err != nil {
		fmt.Println("文件中内容无法读取")
		// 处理复制内容失败的情况
	}

	input := buf.String()
	// 创建一个新的Lexer实例，传入输入字符串，用于对该输入进行词法分析
	lexer := NewLexer(input)

	// 通过循环不断获取下一个标记，直到遇到文件结束标记（EOF）
	for {
		token := lexer.NextToken()
		if token.Type == EOF {
			break
		}
		// 输出每个识别出的标记的类型和值
		fmt.Println(token)
	}
}

// Lexer 结构
// 词法分析器的结构体，用于存储分析过程中的相关状态和数据
type Lexer struct {
	input    string // 输入源代码，即要进行词法分析的C语言代码字符串
	position int    // 当前读取的位置，用于跟踪在输入字符串中的位置
	readPos  int    // 读取字符的位置，通常比position略超前，用于预读取字符
	ch       uint8  // 当前读取的字符，以字节形式的存储
}

// 初始化一个Lexer结构体实例，并读取输入字符串的第一个字符
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 读取第一个字符，使词法分析器处于初始读取状态
	return l
}

func (l *Lexer) readChar() {
	// 如果读取位置已经超过输入字符串的长度，说明已经到达末尾，将当前字符设为0（表示结束）
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		// 否则，从输入字符串中获取当前位置的字符，并转换为uint8类型存储在ch字段中
		l.ch = uint8(l.input[l.readPos])
	}
	l.readPos++
	l.position++
}

// 在输入字符串中，遇到空白字符（空格、制表符、换行符、回车符）时，不断读取下一个字符，直到遇到非空白字符
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 使用unicode包中的IsLetter函数，通过将uint8类型的字符转换为rune类型，判断是否为字母字符
func (l *Lexer) isLetter(ch uint8) bool {
	return unicode.IsLetter(rune(ch))
}

// 使用unicode包中的IsDigit函数，通过将uint8类型的字符转换为rune类型，判断是否为数字字符
func (l *Lexer) isDigit(ch uint8) bool {
	return unicode.IsDigit(rune(ch))
}

// 从当前位置开始，读取连续的字母、数字或下划线字符组成的字符串，作为标识符或关键字
func (l *Lexer) readIdent() string {
	position := l.readPos
	for l.isLetter(l.ch) || l.isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[position-1 : l.readPos-1]
}

// peekChar 函数用于查看下一个字符，但不移动读取位置
func (l *Lexer) peekChar() uint8 {
	if l.readPos >= len(l.input) {
		return 0
	}
	return uint8(l.input[l.readPos])
}

// 获取当前行号和列号的函数
func (l *Lexer) getPositionInfo() (int, int) {
	// 这里简单假设输入字符串是按行存储的，通过换行符来计算行号
	// 列号则通过当前读取位置减去上一个换行符位置来计算（需要记录上一个换行符位置）

	lineNumber := 1
	columnNumber := 1
	lastNewlinePos := -1

	for i, char := range l.input {
		if char == '\n' {
			lineNumber++
			lastNewlinePos = i
		}
	}

	// 如果还未遇到换行符，列号就是当前读取位置 + 1
	if lastNewlinePos == -1 {
		columnNumber = l.position + 1
	} else {
		columnNumber = l.position - lastNewlinePos
	}

	return lineNumber, columnNumber
}

// 这是词法分析器的核心函数，用于分析输入字符串并返回下一个词法单元（标记）
func (l *Lexer) NextToken() Token {
	var token Token

	l.skipWhitespace() // 跳过空白字符，确保从非空白字符开始分析

	if l.ch == 0 { // 检查是否到达文件末尾
		token.Type = EOF
		token.Value = ""
		return token
	}

	if l.isLetter(l.ch) { // 标识符或关键字
		token.Type = IDENT
		token.Value = l.readIdent()
		// 判断读取到的标识符是否为关键字，如果是，则更新标记类型为对应的标记类型
		if tokenType, ok := keywords[token.Value]; ok {
			token.Type = tokenType
		}
		return token
	}

	if l.isDigit(l.ch) { // 整数标记处理
		position := l.readPos
		for l.isDigit(l.ch) {
			l.readChar()
		}
		// 判断下一个字符是否为小数点，如果是则继续处理浮点数
		nextCh := l.ch
		if nextCh == '.' {
			l.readChar()
			// 进入浮点数处理逻辑
			token.Type = FLOAT
			hasDot := false
			for l.isDigit(l.ch) || (l.ch == '.' && !hasDot) {
				if l.ch == '.' {
					hasDot = true
				}
				l.readChar()
			}
			token.Value = l.input[position-1 : l.readPos-1]
			return token
		} else {
			token.Type = INT
			token.Value = l.input[position-1 : l.readPos-1]
			return token
		}
	}

	if l.ch == '.' {
		nextCh := l.peekChar()
		if nextCh == '.' || (!l.isDigit(nextCh) && nextCh != '.') {
			// 单独的 '.'，表示浮点数 0.0
			token.Type = FLOAT
			token.Value = "0.0"
			l.readChar()
			return token
		} else if l.isDigit(nextCh) {
			token.Type = FLOAT
			position := l.readPos
			hasDot := false
			for l.isDigit(l.ch) || (l.ch == '.' && !hasDot) {
				if l.ch == '.' {
					hasDot = true
				}
				l.readChar()
			}
			token.Value = l.input[position-1 : l.readPos-1]
			return token
		}
	}
	if l.ch == '"' { // 字符串标记处理
		token.Type = STRING
		l.readChar() // 跳过开头的双引号
		position := l.readPos
		for l.ch != '"' {
			l.readChar()
		}
		l.readChar() // 跳过结尾的双引号
		token.Value = l.input[position-1 : l.readPos-2]
		return token
	}

	if l.ch == '+' {
		token.Type = PLUS
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '-' {
		token.Type = MINUS
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '*' {
		token.Type = STAR
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '/' {
		// 处理注释情况
		if l.peekChar() == '/' {
			token.Type = COMMENT_SINGLE
			l.readChar()
			l.readChar()
			position := l.readPos
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			token.Value = l.input[position-1 : l.readPos-2]
			return token
		} else if l.peekChar() == '*' {
			token.Type = COMMENT_MULTI
			l.readChar()
			l.readChar()
			position := l.readPos
			for l.peekChar() != 0 && !(l.ch == '*' && l.peekChar() == '/') {
				l.readChar()
			}
			if l.peekChar() == 0 {
				token.Type = ILLEGAL
				token.Value = "多行注释未关闭"
				return token
			} else {
				l.readChar()
				l.readChar()
				token.Value = l.input[position-1 : l.readPos-3]
				return token
			}
		} else {
			token.Type = SLASH
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '=' {
		// 判断是否为 == 运算符
		if l.peekChar() == '=' {
			token.Type = EQEQ
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = EQUAL
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '<' {
		// 判断是否为 <= 运算符
		if l.peekChar() == '=' {
			token.Type = LTEQ
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = LESS
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '>' {
		// 判断是否为 >= 运算符
		if l.peekChar() == '=' {
			token.Type = GTEQ
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = GREATER
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == ';' {
		token.Type = SEMICOLON
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == ',' {
		token.Type = COMMA
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '(' {
		token.Type = LPAREN
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == ')' {
		token.Type = RPAREN
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '{' {
		token.Type = LBRACE
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '}' {
		token.Type = RBRACE
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '#' {
		token.Type = PREPROC
		token.Value = string(l.ch)
		l.readChar()
		position := l.readPos
		for ; l.ch != 0 && !(l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r'); l.readChar() {
		}
		token.Value = l.input[position-1 : l.readPos-1]
		return token
	}

	if l.ch == '&' {
		// 判断是否为 && 运算符
		if l.peekChar() == '&' {
			token.Type = ANDAND
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = AND
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '|' {
		// 判断是否为 || 运算符
		if l.peekChar() == '|' {
			token.Type = OROR
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = OR
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '^' {
		token.Type = XOR
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '~' {
		token.Type = NOT
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == ':' {
		token.Type = COLON
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	if l.ch == '!' {
		// 判断是否为!运算符
		if l.peekChar() == '=' {
			token.Type = NOTNOT
			l.readChar()
			l.readChar()
			return token
		} else {
			token.Type = NOT
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == '%' {
		token.Type = AND
		token.Value = string(l.ch)
		l.readChar()
		return token
	}

	// 处理 << 运算符
	if l.ch == '<' {
		nextCh := l.peekChar()
		if nextCh == '<' {
			token.Type = LSHIFT
			l.readChar()
			l.readChar()
			token.Value = "<<"
			return token
		} else {
			token.Type = LESS
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	// 处理 >> 运算符
	if l.ch == '>' {
		nextCh := l.peekChar()
		if nextCh == '>' {
			token.Type = RSHIFT
			l.readChar()
			l.readChar()
			token.Value = ">>"
			return token
		} else {
			token.Type = GREATER
			token.Value = string(l.ch)
			l.readChar()
			return token
		}
	}

	if l.ch == 0 {
		token.Type = EOF
	} else {
		token.Type = ILLEGAL
		// 记录当前行号和列号，以便提供更详细的错误信息
		lineNumber, columnNumber := l.getPositionInfo()
		token.Value = fmt.Sprintf("非法字符 '%s' 在第 %d 行，第 %d 列", string(l.ch), lineNumber, columnNumber)
	}

	l.readChar() // 读取下一个字符，为下一次分析做准备
	return token
}
