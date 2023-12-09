package evalx

import (
	"errors"
	"fmt"
	"github.com/avicd/go-utilx/conv"
	"github.com/avicd/go-utilx/refx"
	"github.com/avicd/go-utilx/tokx"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

func Eval(text string, cts ...Context) (any, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("empty expression")
	}
	expr := tokx.DoubleQuota(text)
	var stack *Stack
	if len(cts) > 0 && cts[0] != nil {
		stack = StackOf(cts[0])
	} else {
		stack = StackOf(NewScope())
	}
	var astExpr ast.Expr
	if ae, ok := stack.Ctx.CacheOf(expr); ok {
		astExpr = ae
	} else {
		var err error
		astExpr, err = parser.ParseExpr(expr)
		if err != nil {
			return nil, err
		} else {
			stack.Ctx.Cache(expr, astExpr)
		}
	}
	val := evalExpr(astExpr, stack)
	if stack.Error != nil {
		return nil, stack.Error
	}
	return val, nil
}

func evalExpr(buf ast.Expr, stack *Stack) any {
	defer func() {
		rc := recover()
		if rc != nil {
			stack.Error = refx.AsError(rc)
		}
	}()
	switch expr := buf.(type) {
	case *ast.SelectorExpr:
		return evalSelector(expr, stack)
	case *ast.Ident:
		return evalIdent(expr, stack, false)
	case *ast.BasicLit:
		return evalBasicLit(expr)
	case *ast.UnaryExpr:
		return evalUnary(expr, stack)
	case *ast.ParenExpr:
		return evalExpr(expr.X, stack)
	case *ast.IndexExpr:
		return evalIndex(expr, stack)
	case *ast.CallExpr:
		return evalCall(expr, stack)
	case *ast.BinaryExpr:
		return evalBinary(expr, stack)
	}
	return nil
}

func evalSelector(expr *ast.SelectorExpr, stack *Stack) any {
	target := stack.PopTarget()
	var parent any
	switch exprX := expr.X.(type) {
	case *ast.SelectorExpr:
		return evalSelector(exprX, stack)
	case *ast.Ident:
		parent = evalIdent(exprX, stack, true)
	default:
		parent = evalExpr(expr.X, stack)
	}
	if parent != nil {
		if target == METHOD {
			if method, ok := refx.MethodOf(parent, expr.Sel.Name); ok {
				return method
			}
		}
		if val, ok := refx.PropOf(parent, expr.Sel.Name); ok {
			if target == METHOD {
				if refx.IsFunc(val) {
					return val
				}
			} else {
				return val
			}
		}
	}
	return nil
}

func evalIndex(expr *ast.IndexExpr, stack *Stack) any {
	target := stack.PopTarget()
	obj := evalExpr(expr.X, stack)
	index := evalExpr(expr.Index, stack)
	if target == METHOD {
		if method, ok := refx.MethodOf(obj, index); ok {
			return method
		}
	}
	if val, ok := refx.PropOf(obj, index); ok {
		if target == METHOD {
			if refx.IsFunc(val) {
				return val
			}
		} else {
			return val
		}
	}
	return nil
}

func evalIdent(expr *ast.Ident, stack *Stack, sel bool) any {
	target := stack.PopTarget()
	if !sel {
		switch expr.Name {
		case "nil", "null":
			return nil
		case "true":
			return true
		case "false":
			return false
		}
	}
	if target == METHOD {
		if method, ok := stack.Ctx.MethodOf(expr.Name); ok {
			return method
		}
	}
	if val, ok := stack.Ctx.ValueOf(expr.Name); ok {
		if target == METHOD {
			if refx.IsFunc(val) {
				return val
			}
		} else {
			return val
		}
	}
	return nil
}

func evalUnary(expr *ast.UnaryExpr, stack *Stack) any {
	ret := evalExpr(expr.X, stack)
	switch expr.Op {
	case token.NOT:
		return !refx.AsBool(ret)
	case token.SUB:
		if refx.IsInteger(ret) {
			return -refx.AsInt64(ret)
		} else if refx.IsUInteger(ret) {
			return -refx.AsUint64(ret)
		} else if refx.IsFloat(ret) {
			return -refx.AsFloat64(ret)
		} else {
			panic(fmt.Errorf("evalx: invalid operator '%s' on %v", expr.Op, expr.X))
		}
	case token.ADD:
		if refx.IsNumber(ret) {
			panic(fmt.Errorf("evalx: invalid operator '%s' on %v", expr.Op, expr.X))
		}
	}
	return ret
}

func evalCall(expr *ast.CallExpr, stack *Stack) any {
	stack.PushTarget(METHOD)
	val := evalExpr(expr.Fun, stack)
	if val == nil {
		panic(fmt.Errorf("evalx: invalid function"))
		return nil
	}
	method := refx.ValueOf(val)
	var args []reflect.Value
	for _, buf := range expr.Args {
		argVal := evalExpr(buf, stack)
		args = append(args, refx.ValueOf(argVal))
	}
	values := method.Call(args)
	if len(values) > 0 {
		return values[0].Interface()
	}
	return nil
}

func evalBinary(expr *ast.BinaryExpr, stack *Stack) any {
	left := evalExpr(expr.X, stack)
	if expr.Op == token.LOR && refx.AsBool(left) {
		return true
	}
	var right any
	if expr.Op == token.ADD {
		if refx.IsString(left) {
			right = evalString(expr.Y, stack)
		} else {
			right = evalExpr(expr.Y, stack)
			if refx.IsString(right) {
				left = evalString(expr.X, stack)
			}
		}
	} else {
		right = evalExpr(expr.Y, stack)
	}
	stack.X = left
	stack.Y = right
	var ret any
	switch expr.Op {
	case token.LOR:
		ret = refx.AsBool(right)
	case token.LAND:
		ret = refx.AsBool(left) && refx.AsBool(right)
	case token.EQL, token.NEQ, token.GTR, token.LSS, token.GEQ, token.LEQ:
		ret = evalCmpBool(expr, left, right)
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		ret = evalArithmetic(expr, stack)
	case token.AND, token.OR, token.XOR, token.SHL, token.SHR, token.AND_NOT:
		ret = evalBitOpr(expr, stack)
	}
	return ret
}

func evalBasicLit(expr *ast.BasicLit) any {
	switch expr.Kind {
	case token.STRING:
		return strings.Trim(expr.Value, "\"")
	case token.CHAR:
		return strings.Trim(expr.Value, "'")
	case token.INT:
		return conv.ParseInt(expr.Value)
	case token.FLOAT:
		return conv.ParseFloat(expr.Value)
	default:
		return nil
	}
}

func evalString(input ast.Expr, stack *Stack) string {
	if expr, ok := input.(*ast.BinaryExpr); ok {
		if expr.Op == token.ADD {
			left := evalString(expr.X, stack)
			right := evalString(expr.Y, stack)
			return left + right
		}
	}
	return refx.AsString(evalExpr(input, stack))
}

func evalCmpBool(expr *ast.BinaryExpr, x, y any) bool {
	var ret bool
	cmp := refx.Cmp(x, y)
	switch expr.Op {
	case token.EQL:
		ret = cmp == refx.CmpEq
	case token.NEQ:
		ret = cmp != refx.CmpEq
	case token.GTR:
		ret = cmp == refx.CmpGtr
	case token.LSS:
		ret = cmp == refx.CmpLss
	case token.GEQ:
		ret = cmp == refx.CmpGtr || cmp == refx.CmpEq
	case token.LEQ:
		ret = cmp == refx.CmpLss || cmp == refx.CmpEq
	}
	return ret
}

func evalArithmetic(expr *ast.BinaryExpr, stack *Stack) any {
	var ret any
	x := stack.X
	y := stack.Y
	if expr.Op == token.ADD && (refx.IsString(x) || refx.IsString(y)) {
		return refx.AsString(x) + refx.AsString(y)
	}
	if refx.IsNumber(x) && refx.IsNumber(y) {
		if refx.IsGeneralInt(x) && refx.IsGeneralInt(y) {
			a := refx.AsInt64(x)
			b := refx.AsInt64(y)
			switch expr.Op {
			case token.ADD:
				ret = a + b
			case token.SUB:
				ret = a - b
			case token.MUL:
				ret = a * b
			case token.QUO:
				ret = a / b
			case token.REM:
				ret = a % b
			}
		} else {
			a := refx.AsFloat64(x)
			b := refx.AsFloat64(y)
			switch expr.Op {
			case token.ADD:
				ret = a + b
			case token.SUB:
				ret = a - b
			case token.MUL:
				ret = a * b
			case token.QUO:
				ret = a / b
			case token.REM:
				panic(fmt.Errorf("evalx: invalid operator '%s' on float", expr.Op))
			}
		}
	}
	return ret
}

func evalBitOpr(expr *ast.BinaryExpr, stack *Stack) any {
	var val uint64
	x := stack.X
	y := stack.Y
	if refx.IsGeneralInt(x) && refx.IsGeneralInt(y) {
		a := refx.AsUint64(x)
		b := refx.AsUint64(y)
		switch expr.Op {
		case token.AND:
			val = a & b
		case token.OR:
			val = a | b
		case token.XOR:
			val = a ^ b
		case token.SHL:
			val = a << b
		case token.SHR:
			val = a >> b
		case token.AND_NOT:
			val = a &^ b
		}
	} else {
		panic(fmt.Errorf("evalx: bit operator work only on integer"))
	}
	return val
}
