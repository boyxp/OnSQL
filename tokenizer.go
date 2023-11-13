package main

import "fmt"
import "strings"
import "github.com/boyxp/OnSQL/tokenizer"

func main() {
	condition := "a   =  ? and b   >   ? and c >=? and (d<=? or e!=?) and (c in (   ?) or n is     null) and k is not     null"
	tokens := tokenizer.Tokenize(condition)
	fmt.Println(tokens)
	p := parser{}
	fmt.Println(p.parse(tokens))
}

type parser struct{
	index int
	length int
	tokens []string
}

func (p *parser) parse(tokens []string) map[string]interface{} {
    p.index = 0
    p.length = len(tokens)
    p.tokens = tokens

    tree := p._tree()

    return tree
}

func (p *parser) _tree() map[string]interface{} {
	state := 0
	logical := "$and"
	key := ""
	oprts := ""
	var value interface{} = ""
	conds := []map[string]interface{}{}
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
			case "=" :
					oprts = "$eq"
			case "!=" :
					oprts = "$ne"
			case ">" :
					oprts = "$gt"
			case ">=" :
					oprts = "$gte"
			case "<" :
					oprts = "$lt"
			case "<=" :
					oprts = "$lte"
			case "like" :
					oprts = "$like"
			case "regexp" :
					oprts = "$regex"
			case "near" :
					oprts = "$near"
			case "in":
				if oprts != "" && oprts == "not" {
					oprts = "$nin"
				} else {
					oprts = "$in"
				}
			case "is":
				oprts = "is"
				state--
			case "not":
				if oprts != "" && oprts == "is" {
					oprts = "is not"
				} else {
					oprts = "not"
				}
				state--
			case "null":
				p.index--
				if oprts != "" && oprts == "is not" {
					value = 1
				} else {
					value = 0
				}
				oprts = "$exists"
			default:
				panic("syntax error")
			}
		case 2:
			switch token {
			case "null":
			case "?":
				value = token
			default:
				panic("syntax error")
			}

			conds = append(conds, map[string]interface{}{key: map[string]interface{}{oprts: value}})
		case 3:
			switch token {
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
	return map[string]interface{}{"logical": logical, "conds": conds}
}

