// Package contract evaluates a subset of formal pseudocode PRE/POST expressions.
// [IMPL-SPEC_CTL] [ARCH-SPEC_DSL_AND_ORACLE] [REQ-PSEUDOCODE_FORMAL_VERIFICATION]
package contract

import (
	"fmt"
	"strconv"
	"strings"
)

// Env holds named values for expression evaluation (lowercase keys).
type Env map[string]any

// EvalBool evaluates a language-agnostic contract expression against env.
func EvalBool(expr string, env Env) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true, nil
	}
	v, err := parseOr(expr, env)
	if err != nil {
		return false, err
	}
	b, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("expression not boolean: %q", expr)
	}
	return b, nil
}

func parseOr(expr string, env Env) (any, error) {
	parts := splitTop(expr, " OR ")
	if len(parts) > 1 {
		for _, p := range parts {
			v, err := parseAnd(strings.TrimSpace(p), env)
			if err != nil {
				return nil, err
			}
			if b, ok := v.(bool); ok && b {
				return true, nil
			}
		}
		return false, nil
	}
	return parseAnd(expr, env)
}

func parseAnd(expr string, env Env) (any, error) {
	parts := splitTop(expr, " AND ")
	if len(parts) > 1 {
		for _, p := range parts {
			v, err := parseCompare(strings.TrimSpace(p), env)
			if err != nil {
				return nil, err
			}
			if b, ok := v.(bool); ok && !b {
				return false, nil
			}
		}
		return true, nil
	}
	return parseCompare(expr, env)
}

func parseCompare(expr string, env Env) (any, error) {
	for _, op := range []string{"==", "!="} {
		if i := strings.Index(expr, op); i >= 0 {
			left := strings.TrimSpace(expr[:i])
			right := strings.TrimSpace(expr[i+len(op):])
			lv, err := parsePrimary(left, env)
			if err != nil {
				return nil, err
			}
			rv, err := parsePrimary(right, env)
			if err != nil {
				return nil, err
			}
			eq, err := valuesEqual(lv, rv)
			if err != nil {
				return nil, err
			}
			if op == "==" {
				return eq, nil
			}
			return !eq, nil
		}
	}
	return parsePrimary(expr, env)
}

func parsePrimary(expr string, env Env) (any, error) {
	expr = strings.TrimSpace(expr)
	if strings.HasPrefix(expr, "NOT ") {
		v, err := parsePrimary(strings.TrimSpace(expr[4:]), env)
		if err != nil {
			return nil, err
		}
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("NOT applied to non-bool: %q", expr)
		}
		return !b, nil
	}
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return parseOr(expr[1:len(expr)-1], env)
	}
	if v, ok := parseLiteral(expr); ok {
		return v, nil
	}
	if strings.Contains(expr, "(") && strings.HasSuffix(expr, ")") {
		return evalCall(expr, env)
	}
	if v, ok := lookupPath(expr, env); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unknown identifier: %q", expr)
}

func parseLiteral(s string) (any, bool) {
	switch strings.ToLower(s) {
	case "true":
		return true, true
	case "false":
		return false, true
	case "ok":
		return "ok", true
	}
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		return s[1 : len(s)-1], true
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n, true
	}
	return nil, false
}

func lookupPath(path string, env Env) (any, bool) {
	parts := strings.Split(path, ".")
	var cur any = env
	for _, p := range parts {
		m, ok := cur.(Env)
		if !ok {
			if p == parts[0] {
				if v, ok := env[strings.ToLower(path)]; ok {
					return v, true
				}
			}
			return nil, false
		}
		key := strings.ToLower(p)
		v, ok := m[key]
		if !ok {
			// camelCase field alias: HasHandle -> hashandle
			key = strings.ToLower(p)
			v, ok = m[key]
		}
		if !ok {
			return nil, false
		}
		cur = v
	}
	return cur, true
}

func evalCall(expr string, env Env) (any, error) {
	i := strings.Index(expr, "(")
	if i < 0 {
		return nil, fmt.Errorf("invalid call: %q", expr)
	}
	name := strings.TrimSpace(expr[:i])
	argsStr := strings.TrimSpace(expr[i+1 : len(expr)-1])
	var args []string
	if argsStr != "" {
		args = []string{strings.TrimSpace(argsStr)}
	}
	switch name {
	case "VALIDATE_PATH":
		if len(args) != 1 {
			return nil, fmt.Errorf("VALIDATE_PATH wants 1 arg")
		}
		path, err := parsePrimary(args[0], env)
		if err != nil {
			return nil, err
		}
		ps, _ := path.(string)
		if ps == "" || strings.Contains(ps, "..") {
			return "error", nil
		}
		return "ok", nil
	default:
		return nil, fmt.Errorf("unknown call: %s", name)
	}
}

func valuesEqual(a, b any) (bool, error) {
	if sa, ok := a.(string); ok {
		if sb, ok := b.(string); ok {
			return sa == sb, nil
		}
	}
	if ba, ok := a.(bool); ok {
		if bb, ok := b.(bool); ok {
			return ba == bb, nil
		}
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b), nil
}

func splitTop(expr, sep string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i <= len(expr)-len(sep); i++ {
		switch expr[i] {
		case '(':
			depth++
		case ')':
			depth--
		}
		if depth == 0 && strings.EqualFold(expr[i:i+len(sep)], sep) {
			parts = append(parts, expr[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	parts = append(parts, expr[start:])
	if len(parts) == 1 {
		return parts
	}
	return parts
}

// WriterEnv builds contract.Env from atomic writer oracle fields.
func WriterEnv(targetPath, tempPath string, hasHandle, committed, closed bool) Env {
	return Env{
		"writer": Env{
			"targetpath": targetPath,
			"temppath":   tempPath,
			"hashandle":  hasHandle,
			"committed":  committed,
			"closed":     closed,
		},
		"targetpath": targetPath,
	}
}
