package MongoDB

import "sync"
import "strings"
import "github.com/boyxp/OnSQL/tokenizer"

var filter sync.Map

type Parser struct{
	index int
	length int
	placeholder int
	tokens []string
}

func (p *Parser) Parse(condition string) map[string]interface{} {
	value, ok := filter.Load(condition);
    if ok {
       	return value.(map[string]interface{})
    }

	tokens  := tokenizer.Tokenize(condition)
    p.index  = 0
    p.length = len(tokens)
    p.tokens = tokens
    p.placeholder = 0

    tree := p._tree()

    if len(tree)==0 {
    	panic("syntax error")
    }

    filter.Store(condition, tree)

    return tree
}

func (p *Parser) _tree() map[string]interface{} {
	state   := 0
	logical := "$and"
	key     := ""
	opr     := ""
	conds   := []map[string]interface{}{}

	var value interface{}

	for ; p.index < p.length; p.index++ {
		token := p.tokens[p.index]
		switch state {
			case 0:
				switch token {
					case "(":
							p.index++
							child := p._tree()
							conds = append(conds, map[string]interface{}{child["logical"].(string): child["conds"].([]map[string]interface{})})
							state = 2
					default:
							key = token
					}
			case 1:
					switch strings.ToLower(token) {
						case "="  :
									opr = "$eq"
						case "!=" :
									opr = "$ne"
						case ">"  :
									opr = "$gt"
						case ">=" :
									opr = "$gte"
						case "<"  :
									opr = "$lt"
						case "<=" :
									opr = "$lte"
						case "like" :
									opr = "$like"
						case "regexp" :
									opr = "$regex"
						case "near" :
									opr = "$near"
						case "in":
									if opr != "" && opr == "not" {
										opr = "$nin"
									} else {
										opr = "$in"
									}
						case "is":
									opr = "is"
									state--
						case "not":
									if opr != "" && opr == "is" {
										opr = "is not"
									} else {
										opr = "not"
									}
									state--
						case "null":
									p.index--
									if opr != "" && opr == "is not" {
										value = true
									} else {
										value = false
									}
									opr = "$exists"
						default:
								panic("syntax error")
					}
			case 2:
					switch strings.ToLower(token) {
						case "null":
						case "?":
								value = token
						default:
								panic("syntax error")
					}

					var placeholder int
					if value=="?" {
						placeholder = p.placeholder
						p.placeholder++
					} else {
						placeholder = -1
					}

					conds = append(conds, map[string]interface{}{key: map[string]interface{}{"opr":opr, "value":value, "placeholder":placeholder}})
			case 3:
					switch strings.ToLower(token) {
						case ")":
								return map[string]interface{}{"logical": logical, "conds": conds}
						case "and":
								logical = "$and"
								state = -1
						case "or":
								logical = "$or"
								state = -1
						default:
								panic("syntax error")
					}
			default:
					panic("syntax error")
		}

		state++
	}

	return map[string]interface{}{logical: conds}
}
