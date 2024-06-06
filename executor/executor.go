package parser

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"

	p "github.com/mandriota/mareal/parser"
	u "github.com/mandriota/mareal/utils"
)

var (
	Scanner = bufio.NewScanner(os.Stdin)
	oWriter = bufio.NewWriter(os.Stdout)
)

type Executer struct {
	global map[string]*p.Node
}

func Execute(src string) (err error) {
	e := new(Executer)
	e.init()

	srcTree, err := p.New(src).Parse()
	if err != nil {
		return err
	}

	retTrees, err := e.executeTree(e.global, e.global, srcTree)
	if err != nil {
		return err
	}

	if retTrees != nil {
		e.executeWrite(e.global, e.global, retTrees...)
	}

	return nil
}

func (e *Executer) init() {
	e.global = map[string]*p.Node{
		"nl": {Token: p.Token{Typ: p.TkStr, Val: "\n"}},
		"pi": {Token: p.Token{Typ: p.TkNum, Val: math.Pi}},
		"e":  {Token: p.Token{Typ: p.TkNum, Val: math.E}},
	}
}

func (e *Executer) numMap(superScope, localScope map[string]*p.Node, srcTree *p.Node, mapper func(acc, v float64) float64) (p.Buff, error) {
	if len(srcTree.Component) < 3 {
		return nil, fmt.Errorf("function \"%v\": expected at least 2 args", srcTree.Component[0])
	}

	fnName := srcTree.Component[0].Val.(string)

	argRet := new(p.Node)
	accRets, err := e.executeTree(superScope, localScope, srcTree.Component[1])
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": error during arg parsing: %v", fnName, err)
	}

	if err := accRets.Sub(argRet); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	if err := argRet.Typ.Assert(p.TkNum); err != nil {
		return nil, fmt.Errorf("function \"%s\": error during arg parsing: %v", fnName, err)
	}

	accRet := argRet.Val.(float64)

	for _, el := range srcTree.Component[2:] {
		argRets, err := e.executeTree(superScope, localScope, el)
		if err != nil {
			return nil, fmt.Errorf("function \"%s\": error during arg parsing: %v", fnName, err)
		}

		if err := argRets.Sub(argRet); err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		if err := argRet.Typ.Assert(p.TkNum); err != nil {
			return nil, fmt.Errorf("function \"%s\": error during arg parsing: %v", fnName, err)
		}

		accRet = mapper(accRet, argRet.Val.(float64))
	}

	return p.Buff{&p.Node{Token: p.Token{Typ: p.TkNum, Val: accRet}}}, err
}

func (e *Executer) executeTree(superScope, localScope map[string]*p.Node, srcTree *p.Node, args ...*p.Node) (retTrees p.Buff, err error) {
	switch srcTree.Typ {
	case p.TkNum, p.TkStr, p.TkArr:
		return p.Buff{srcTree}, nil
	case p.TkIdent:
		let, ok := localScope[srcTree.Val.(string)]
		if !ok {
			if let, ok = superScope[srcTree.Val.(string)]; !ok {
				return nil, fmt.Errorf("variable \"%s\" is not declared in the current scope", srcTree.String())
			}
		}

		return e.executeTree(superScope, localScope, let)
	case p.TkRoutine:
		if len(srcTree.Component) == 0 {
			return nil, fmt.Errorf("function name expected")
		}

		fnName := srcTree.Component[0]
		if err := fnName.Typ.Assert(p.TkIdent); err != nil {
			return nil, fmt.Errorf("parsing function name: %v", err)
		}

		switch fnName := fnName.Val.(string); fnName {
		case "_":
			return e.executeReturn(fnName, superScope, localScope, srcTree)
		case "quote":
			return srcTree.Component[1:], nil
		case "lambda":
			return e.executeLambda(fnName, superScope, localScope, srcTree, args...)
		case "def":
			return e.executeDef(fnName, superScope, localScope, srcTree)
		case "set":
			return e.executeSet(fnName, superScope, localScope, srcTree)
		case "rep":
			return e.executeRep(fnName, superScope, localScope, srcTree)
		case "map":
			return e.executeMap(fnName, superScope, localScope, srcTree)
		case "if":
			return e.executeIf(fnName, superScope, localScope, srcTree)
		case "+":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return acc + v
			})
		case "-":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return acc - v
			})
		case "*":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return acc * v
			})
		case "/":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return acc / v
			})
		case "%":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return math.Mod(acc, v)
			})
		case "^":
			return e.numMap(superScope, localScope, srcTree, func(acc, v float64) float64 {
				return math.Pow(acc, v)
			})
		case "get":
			if err := e.executeWrite(superScope, localScope, srcTree.Component[1:]...); err != nil {
				return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
			}

			Scanner.Scan()
			return p.Buff{&p.Node{Token: p.Token{Typ: p.TkStr, Val: Scanner.Text()}}}, nil
		case "put":
			return nil, e.executeWrite(superScope, localScope, srcTree.Component[1:]...)
		default:
			fnBody, ok := localScope[fnName]
			if !ok {
				if fnBody, ok = superScope[fnName]; !ok {
					return nil, fmt.Errorf("function \"%s\": not declared in the current scope", fnName)
				}
			}

			if err := fnBody.Typ.Assert(p.TkRoutine); err != nil {
				return nil, fmt.Errorf("function \"%s\": %v", fnBody, err)
			}

			childArgs := make(p.Buff, 0)

			for _, arg := range srcTree.Component[1:] {
				argRets, err := e.executeTree(superScope, localScope, arg)
				if err != nil {
					return nil, err
				}
				childArgs.Add(argRets...)
			}

			u.MapCopyNoOverwrite(localScope, superScope)

			if retTrees, err = e.executeTree(localScope, make(map[string]*p.Node), fnBody, childArgs...); err != nil {
				return nil, fmt.Errorf("function \"%s\": %v", fnBody, err)
			}
		}
	}

	return
}

func (e *Executer) executeReturn(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	u.MapCopyNoOverwrite(localScope, superScope)

	retTrees := make(p.Buff, 0)

	for _, arg := range srcTree.Component[1:] {
		argRets, err := e.executeTree(localScope, make(map[string]*p.Node), arg)
		if err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		retTrees.Add(argRets...)
	}

	return retTrees, nil
}

func (e *Executer) executeLambda(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node, args ...*p.Node) (p.Buff, error) {
	if nArgs := len(srcTree.Component); nArgs != 3 {
		return nil, fmt.Errorf("function \"%s\": expects 3 arguments: received %d", fnName, nArgs)
	}

	err := errors.Join(srcTree.Component[1].Typ.Assert(p.TkRoutine), srcTree.Component[2].Typ.Assert(p.TkRoutine))
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	sgnRets, err := e.executeTree(localScope, localScope, srcTree.Component[1])
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	if len(sgnRets) != len(args) {
		return nil, fmt.Errorf("function \"%s\": expected %d args: found %d", fnName, len(sgnRets), len(args))
	}

	childScope := make(map[string]*p.Node, len(sgnRets))
	for i, sgnRet := range sgnRets {
		if err := sgnRet.Typ.Assert(p.TkIdent); err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		childScope[sgnRet.Val.(string)] = args[i]
	}

	u.MapCopyNoOverwrite(localScope, superScope)
	return e.executeTree(localScope, childScope, srcTree.Component[2])
}

func (e *Executer) executeDef(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	if len(srcTree.Component)%2 == 0 {
		return nil, fmt.Errorf("function \"%s\": missing assigned value", fnName)
	}

	for i := 1; i+1 < len(srcTree.Component); i += 2 {
		if err := srcTree.Component[i].Typ.Assert(p.TkIdent); err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		superScope[srcTree.Component[i].Val.(string)] = srcTree.Component[i+1]
	}

	return nil, nil
}

func (e *Executer) executeSet(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	if len(srcTree.Component)%2 == 0 {
		return nil, fmt.Errorf("function \"%s\": missing assigned value", fnName)
	}

	for i := 1; i+1 < len(srcTree.Component); i += 2 {
		if err := srcTree.Component[i].Typ.Assert(p.TkIdent); err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		litRet := new(p.Node)
		litRets, err := e.executeTree(superScope, localScope, srcTree.Component[i+1])
		if err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		if err := litRets.Sub(litRet); err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		superScope[srcTree.Component[i].Val.(string)] = litRet
	}

	return nil, nil
}

func (e *Executer) executeRep(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	if nArgs := len(srcTree.Component); nArgs != 4 {
		return nil, fmt.Errorf("function \"%s\": expects 4 arguments: received %d", fnName, nArgs)
	}

	err := errors.Join(srcTree.Component[1].Typ.Assert(p.TkIdent), srcTree.Component[3].Typ.Assert(p.TkRoutine))
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	maxRet := new(p.Node)
	maxRets, err := e.executeTree(superScope, localScope, srcTree.Component[2])
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	if err := maxRets.Sub(maxRet); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	if err := maxRet.Typ.Assert(p.TkNum); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	u.MapCopyNoOverwrite(localScope, superScope)

	retTrees := make(p.Buff, 0)
	itrName := srcTree.Component[1].Val.(string)
	rangeToRetN := int(maxRet.Val.(float64))

	for i := 0; i < rangeToRetN; i++ {
		repFnRets, err := e.executeTree(localScope, map[string]*p.Node{
			itrName: &p.Node{Token: p.Token{Typ: p.TkNum, Val: float64(i)}},
		}, srcTree.Component[3])
		if err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		retTrees.Add(repFnRets...)
	}

	return retTrees, nil
}

func (e *Executer) executeMap(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	if nArgs := len(srcTree.Component); nArgs != 4 {
		return nil, fmt.Errorf("function \"%s\" expects 4 arguments: received %d", fnName, nArgs)
	}

	if err := errors.Join(srcTree.Component[1].Typ.Assert(p.TkIdent), srcTree.Component[3].Typ.Assert(p.TkRoutine)); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	arrRets, err := e.executeTree(superScope, localScope, srcTree.Component[2])
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	u.MapCopyNoOverwrite(localScope, superScope)

	retTrees := make(p.Buff, 0)
	itrName := srcTree.Component[1].Val.(string)

	for _, localScope[itrName] = range arrRets {
		mapFnRets, err := e.executeTree(localScope, localScope, srcTree.Component[3])
		if err != nil {
			return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
		}

		retTrees.Add(mapFnRets...)
	}

	return retTrees, nil
}

func (e *Executer) executeIf(fnName string, superScope, localScope map[string]*p.Node, srcTree *p.Node) (p.Buff, error) {
	nArgs := len(srcTree.Component)
	if nArgs < 3 || nArgs > 4 {
		return nil, fmt.Errorf("function \"%s\" expects 3 or 4 arguments: received %d", fnName, nArgs)
	}

	cndRets, err := e.executeTree(superScope, localScope, srcTree.Component[1])
	if err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	cndRet := new(p.Node)
	if err := cndRets.Sub(cndRet); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	if err := cndRet.Typ.Assert(p.TkNum); err != nil {
		return nil, fmt.Errorf("function \"%s\": %v", fnName, err)
	}

	body := (*p.Node)(nil)
	if condRet := cndRet.Val.(float64); condRet <= -1 || condRet >= 1 {
		body = srcTree.Component[2]
	} else if nArgs == 4 {
		body = srcTree.Component[3]
	}

	return e.executeTree(superScope, localScope, body)
}

func (e *Executer) executeWrite(superScope, localScope map[string]*p.Node, args ...*p.Node) error {
	for _, arg := range args {
		argRets, err := e.executeTree(superScope, localScope, arg)
		if err != nil {
			return fmt.Errorf("parsing %v: %v", arg, err)
		}

		for _, argRet := range argRets {
			oWriter.WriteString(argRet.String())
		}
	}

	oWriter.Flush()
	return nil
}
