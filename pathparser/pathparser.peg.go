package pathparser

//go:generate peg -inline -switch pathparser.peg

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulepath
	ruleexpr
	rulesegment
	ruleoptional
	rulecapture
	rulename
	rulevalues
	ruleoption
	ruleliteral
	ruleliteral_chars
	ruleslash
	rulebo
	ruleeo
	rulebc
	ruleec
	rulealpha
	rulenum
	ruleset0
	ruleset1
	ruleany
	ruleeos
	ruleAction0
	rulePegText
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8
	ruleAction9
)

var rul3s = [...]string{
	"Unknown",
	"path",
	"expr",
	"segment",
	"optional",
	"capture",
	"name",
	"values",
	"option",
	"literal",
	"literal_chars",
	"slash",
	"bo",
	"eo",
	"bc",
	"ec",
	"alpha",
	"num",
	"set0",
	"set1",
	"any",
	"eos",
	"Action0",
	"PegText",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",
	"Action9",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type Peg struct {
	Path

	Buffer string
	buffer []rune
	rules  [33]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *Peg) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *Peg) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *Peg
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *Peg) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *Peg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.finish()
		case ruleAction1:
			p.addNode(nodeCaptureIdentifier, text)
		case ruleAction2:
			p.addNode(nodeCaptureOption, text)
		case ruleAction3:
			p.addNode(nodeLiteral, text)
		case ruleAction4:
			p.addNode(nodeLiteral, "\\.")
		case ruleAction5:
			p.addNode(nodeLiteral, "/")
		case ruleAction6:
			p.beginOptional()
		case ruleAction7:
			p.endOptional()
		case ruleAction8:
			p.beginCapture()
		case ruleAction9:
			p.endCapture()

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *Peg) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 path <- <(segment+ slash? eos Action0)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulesegment]() {
					goto l0
				}
			l2:
				{
					position3, tokenIndex3 := position, tokenIndex
					if !_rules[rulesegment]() {
						goto l3
					}
					goto l2
				l3:
					position, tokenIndex = position3, tokenIndex3
				}
				{
					position4, tokenIndex4 := position, tokenIndex
					if !_rules[ruleslash]() {
						goto l4
					}
					goto l5
				l4:
					position, tokenIndex = position4, tokenIndex4
				}
			l5:
				{
					position6 := position
					{
						position7, tokenIndex7 := position, tokenIndex
						if !matchDot() {
							goto l7
						}
						goto l0
					l7:
						position, tokenIndex = position7, tokenIndex7
					}
					add(ruleeos, position6)
				}
				{
					add(ruleAction0, position)
				}
				add(rulepath, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 expr <- <((&('{') capture) | (&('[') optional) | (&('%' | '+' | ',' | '-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') literal))> */
		func() bool {
			position9, tokenIndex9 := position, tokenIndex
			{
				position10 := position
				{
					switch buffer[position] {
					case '{':
						{
							position12 := position
							{
								position13 := position
								if buffer[position] != rune('{') {
									goto l9
								}
								position++
								{
									add(ruleAction8, position)
								}
								add(rulebc, position13)
							}
							{
								position15 := position
								{
									position16 := position
									if !_rules[rulealpha]() {
										goto l9
									}
								l17:
									{
										position18, tokenIndex18 := position, tokenIndex
										if !_rules[ruleset0]() {
											goto l18
										}
										goto l17
									l18:
										position, tokenIndex = position18, tokenIndex18
									}
									add(rulePegText, position16)
								}
								{
									add(ruleAction1, position)
								}
								add(rulename, position15)
							}
							{
								position20, tokenIndex20 := position, tokenIndex
								if buffer[position] != rune(':') {
									goto l20
								}
								position++
								{
									position22 := position
									if !_rules[ruleoption]() {
										goto l20
									}
								l23:
									{
										position24, tokenIndex24 := position, tokenIndex
										if buffer[position] != rune('|') {
											goto l24
										}
										position++
										if !_rules[ruleoption]() {
											goto l24
										}
										goto l23
									l24:
										position, tokenIndex = position24, tokenIndex24
									}
									add(rulevalues, position22)
								}
								goto l21
							l20:
								position, tokenIndex = position20, tokenIndex20
							}
						l21:
							{
								position25 := position
								if buffer[position] != rune('}') {
									goto l9
								}
								position++
								{
									add(ruleAction9, position)
								}
								add(ruleec, position25)
							}
							add(rulecapture, position12)
						}
						break
					case '[':
						{
							position27 := position
							{
								position28 := position
								if buffer[position] != rune('[') {
									goto l9
								}
								position++
								{
									add(ruleAction6, position)
								}
								add(rulebo, position28)
							}
							if !_rules[ruleexpr]() {
								goto l9
							}
						l30:
							{
								position31, tokenIndex31 := position, tokenIndex
								if !_rules[rulesegment]() {
									goto l31
								}
								goto l30
							l31:
								position, tokenIndex = position31, tokenIndex31
							}
							{
								position32 := position
								if buffer[position] != rune(']') {
									goto l9
								}
								position++
								{
									add(ruleAction7, position)
								}
								add(ruleeo, position32)
							}
							add(ruleoptional, position27)
						}
						break
					default:
						{
							position34 := position
							{
								position37 := position
								{
									position38, tokenIndex38 := position, tokenIndex
									{
										position40 := position
										if !_rules[ruleset1]() {
											goto l39
										}
									l41:
										{
											position42, tokenIndex42 := position, tokenIndex
											if !_rules[ruleset1]() {
												goto l42
											}
											goto l41
										l42:
											position, tokenIndex = position42, tokenIndex42
										}
										add(rulePegText, position40)
									}
									{
										add(ruleAction3, position)
									}
									goto l38
								l39:
									position, tokenIndex = position38, tokenIndex38
									if buffer[position] != rune('.') {
										goto l9
									}
									position++
									{
										add(ruleAction4, position)
									}
								}
							l38:
								add(ruleliteral_chars, position37)
							}
						l35:
							{
								position36, tokenIndex36 := position, tokenIndex
								{
									position45 := position
									{
										position46, tokenIndex46 := position, tokenIndex
										{
											position48 := position
											if !_rules[ruleset1]() {
												goto l47
											}
										l49:
											{
												position50, tokenIndex50 := position, tokenIndex
												if !_rules[ruleset1]() {
													goto l50
												}
												goto l49
											l50:
												position, tokenIndex = position50, tokenIndex50
											}
											add(rulePegText, position48)
										}
										{
											add(ruleAction3, position)
										}
										goto l46
									l47:
										position, tokenIndex = position46, tokenIndex46
										if buffer[position] != rune('.') {
											goto l36
										}
										position++
										{
											add(ruleAction4, position)
										}
									}
								l46:
									add(ruleliteral_chars, position45)
								}
								goto l35
							l36:
								position, tokenIndex = position36, tokenIndex36
							}
							add(ruleliteral, position34)
						}
						break
					}
				}

				add(ruleexpr, position10)
			}
			return true
		l9:
			position, tokenIndex = position9, tokenIndex9
			return false
		},
		/* 2 segment <- <(slash expr)> */
		func() bool {
			position53, tokenIndex53 := position, tokenIndex
			{
				position54 := position
				if !_rules[ruleslash]() {
					goto l53
				}
				if !_rules[ruleexpr]() {
					goto l53
				}
				add(rulesegment, position54)
			}
			return true
		l53:
			position, tokenIndex = position53, tokenIndex53
			return false
		},
		/* 3 optional <- <(bo expr segment* eo)> */
		nil,
		/* 4 capture <- <(bc name (':' values)? ec)> */
		nil,
		/* 5 name <- <(<(alpha set0*)> Action1)> */
		nil,
		/* 6 values <- <(option ('|' option)*)> */
		nil,
		/* 7 option <- <(<any+> Action2)> */
		func() bool {
			position59, tokenIndex59 := position, tokenIndex
			{
				position60 := position
				{
					position61 := position
					{
						position64 := position
						{
							switch buffer[position] {
							case '\\':
								if buffer[position] != rune('\\') {
									goto l59
								}
								position++
								break
							case '.':
								if buffer[position] != rune('.') {
									goto l59
								}
								position++
								break
							default:
								if !_rules[ruleset1]() {
									goto l59
								}
								break
							}
						}

						add(ruleany, position64)
					}
				l62:
					{
						position63, tokenIndex63 := position, tokenIndex
						{
							position66 := position
							{
								switch buffer[position] {
								case '\\':
									if buffer[position] != rune('\\') {
										goto l63
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l63
									}
									position++
									break
								default:
									if !_rules[ruleset1]() {
										goto l63
									}
									break
								}
							}

							add(ruleany, position66)
						}
						goto l62
					l63:
						position, tokenIndex = position63, tokenIndex63
					}
					add(rulePegText, position61)
				}
				{
					add(ruleAction2, position)
				}
				add(ruleoption, position60)
			}
			return true
		l59:
			position, tokenIndex = position59, tokenIndex59
			return false
		},
		/* 8 literal <- <literal_chars+> */
		nil,
		/* 9 literal_chars <- <((<set1+> Action3) / ('.' Action4))> */
		nil,
		/* 10 slash <- <('/' Action5)> */
		func() bool {
			position71, tokenIndex71 := position, tokenIndex
			{
				position72 := position
				if buffer[position] != rune('/') {
					goto l71
				}
				position++
				{
					add(ruleAction5, position)
				}
				add(ruleslash, position72)
			}
			return true
		l71:
			position, tokenIndex = position71, tokenIndex71
			return false
		},
		/* 11 bo <- <('[' Action6)> */
		nil,
		/* 12 eo <- <(']' Action7)> */
		nil,
		/* 13 bc <- <('{' Action8)> */
		nil,
		/* 14 ec <- <('}' Action9)> */
		nil,
		/* 15 alpha <- <([a-z] / [A-Z])> */
		func() bool {
			position78, tokenIndex78 := position, tokenIndex
			{
				position79 := position
				{
					position80, tokenIndex80 := position, tokenIndex
					if c := buffer[position]; c < rune('a') || c > rune('z') {
						goto l81
					}
					position++
					goto l80
				l81:
					position, tokenIndex = position80, tokenIndex80
					if c := buffer[position]; c < rune('A') || c > rune('Z') {
						goto l78
					}
					position++
				}
			l80:
				add(rulealpha, position79)
			}
			return true
		l78:
			position, tokenIndex = position78, tokenIndex78
			return false
		},
		/* 16 num <- <[0-9]> */
		nil,
		/* 17 set0 <- <((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') num) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') alpha))> */
		func() bool {
			position83, tokenIndex83 := position, tokenIndex
			{
				position84 := position
				{
					switch buffer[position] {
					case '_':
						if buffer[position] != rune('_') {
							goto l83
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						{
							position86 := position
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l83
							}
							position++
							add(rulenum, position86)
						}
						break
					default:
						if !_rules[rulealpha]() {
							goto l83
						}
						break
					}
				}

				add(ruleset0, position84)
			}
			return true
		l83:
			position, tokenIndex = position83, tokenIndex83
			return false
		},
		/* 18 set1 <- <((&('%') '%') | (&(',') ',') | (&('+') '+') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') set0))> */
		func() bool {
			position87, tokenIndex87 := position, tokenIndex
			{
				position88 := position
				{
					switch buffer[position] {
					case '%':
						if buffer[position] != rune('%') {
							goto l87
						}
						position++
						break
					case ',':
						if buffer[position] != rune(',') {
							goto l87
						}
						position++
						break
					case '+':
						if buffer[position] != rune('+') {
							goto l87
						}
						position++
						break
					case '-':
						if buffer[position] != rune('-') {
							goto l87
						}
						position++
						break
					default:
						if !_rules[ruleset0]() {
							goto l87
						}
						break
					}
				}

				add(ruleset1, position88)
			}
			return true
		l87:
			position, tokenIndex = position87, tokenIndex87
			return false
		},
		/* 19 any <- <((&('\\') '\\') | (&('.') '.') | (&('%' | '+' | ',' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') set1))> */
		nil,
		/* 20 eos <- <!.> */
		nil,
		/* 22 Action0 <- <{ p.finish() }> */
		nil,
		nil,
		/* 24 Action1 <- <{ p.addNode(nodeCaptureIdentifier, text) }> */
		nil,
		/* 25 Action2 <- <{ p.addNode(nodeCaptureOption, text) }> */
		nil,
		/* 26 Action3 <- <{ p.addNode(nodeLiteral, text) }> */
		nil,
		/* 27 Action4 <- <{ p.addNode(nodeLiteral, "\\.") }> */
		nil,
		/* 28 Action5 <- <{ p.addNode(nodeLiteral, "/") }> */
		nil,
		/* 29 Action6 <- <{ p.beginOptional() }> */
		nil,
		/* 30 Action7 <- <{ p.endOptional() }> */
		nil,
		/* 31 Action8 <- <{ p.beginCapture() }> */
		nil,
		/* 32 Action9 <- <{ p.endCapture() }> */
		nil,
	}
	p.rules = _rules
}
