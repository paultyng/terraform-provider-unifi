package ginkgolinter

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/printer"
	"go/token"
	gotypes "go/types"
	"reflect"

	"github.com/go-toolsmith/astcopy"
	"golang.org/x/tools/go/analysis"

	"github.com/nunnatsa/ginkgolinter/ginkgohandler"
	"github.com/nunnatsa/ginkgolinter/gomegahandler"
	"github.com/nunnatsa/ginkgolinter/interfaces"
	"github.com/nunnatsa/ginkgolinter/reverseassertion"
	"github.com/nunnatsa/ginkgolinter/types"
	"github.com/nunnatsa/ginkgolinter/version"
)

// The ginkgolinter enforces standards of using ginkgo and gomega.
//
// For more details, look at the README.md file

const (
	linterName                    = "ginkgo-linter"
	wrongLengthWarningTemplate    = linterName + ": wrong length assertion; consider using `%s` instead"
	wrongNilWarningTemplate       = linterName + ": wrong nil assertion; consider using `%s` instead"
	wrongBoolWarningTemplate      = linterName + ": wrong boolean assertion; consider using `%s` instead"
	wrongErrWarningTemplate       = linterName + ": wrong error assertion; consider using `%s` instead"
	wrongCompareWarningTemplate   = linterName + ": wrong comparison assertion; consider using `%s` instead"
	doubleNegativeWarningTemplate = linterName + ": avoid double negative assertion; consider using `%s` instead"
	valueInEventually             = linterName + ": use a function call in %s. This actually checks nothing, because %s receives the function returned value, instead of function itself, and this value is never changed"
	comparePointerToValue         = linterName + ": comparing a pointer to a value will always fail. consider using `%s` instead"
	missingAssertionMessage       = linterName + `: %q: missing assertion method. Expected "Should()", "To()", "ShouldNot()", "ToNot()" or "NotTo()"`
	missingAsyncAssertionMessage  = linterName + `: %q: missing assertion method. Expected "Should()" or "ShouldNot()"`
	focusContainerFound           = linterName + ": Focus container found. This is used only for local debug and should not be part of the actual source code, consider to replace with %q"
	focusSpecFound                = linterName + ": Focus spec found. This is used only for local debug and should not be part of the actual source code, consider to remove it"
	compareDifferentTypes         = linterName + ": use %[1]s with different types: Comparing %[2]s with %[3]s; either change the expected value type if possible, or use the BeEquivalentTo() matcher, instead of %[1]s()"
	matchErrorArgWrongType        = linterName + ": the MatchError matcher used to assert a non error type (%s)"
	matchErrorWrongTypeAssertion  = linterName + ": MatchError first parameter (%s) must be error, string, GomegaMatcher or func(error)bool are allowed"
	matchErrorMissingDescription  = linterName + ": missing function description as second parameter of MatchError"
	matchErrorRedundantArg        = linterName + ": redundant MatchError arguments; consider removing them; consider using `%s` instead"
	matchErrorNoFuncDescription   = linterName + ": The second parameter of MatchError must be the function description (string)"
)

const ( // gomega matchers
	beEmpty        = "BeEmpty"
	beEquivalentTo = "BeEquivalentTo"
	beFalse        = "BeFalse"
	beIdenticalTo  = "BeIdenticalTo"
	beNil          = "BeNil"
	beNumerically  = "BeNumerically"
	beTrue         = "BeTrue"
	beZero         = "BeZero"
	equal          = "Equal"
	haveLen        = "HaveLen"
	haveOccurred   = "HaveOccurred"
	haveValue      = "HaveValue"
	not            = "Not"
	omega          = "Ω"
	succeed        = "Succeed"
	and            = "And"
	or             = "Or"
	withTransform  = "WithTransform"
	matchError     = "MatchError"
)

const ( // gomega actuals
	expect                 = "Expect"
	expectWithOffset       = "ExpectWithOffset"
	eventually             = "Eventually"
	eventuallyWithOffset   = "EventuallyWithOffset"
	consistently           = "Consistently"
	consistentlyWithOffset = "ConsistentlyWithOffset"
)

// Analyzer is the interface to go_vet
var Analyzer = NewAnalyzer()

type ginkgoLinter struct {
	config *types.Config
}

// NewAnalyzer returns an Analyzer - the package interface with nogo
func NewAnalyzer() *analysis.Analyzer {
	linter := ginkgoLinter{
		config: &types.Config{
			SuppressLen:     false,
			SuppressNil:     false,
			SuppressErr:     false,
			SuppressCompare: false,
			ForbidFocus:     false,
			AllowHaveLen0:   false,
		},
	}

	a := &analysis.Analyzer{
		Name: "ginkgolinter",
		Doc:  fmt.Sprintf(doc, version.Version()),
		Run:  linter.run,
	}

	var ignored bool
	a.Flags.Init("ginkgolinter", flag.ExitOnError)
	a.Flags.Var(&linter.config.SuppressLen, "suppress-len-assertion", "Suppress warning for wrong length assertions")
	a.Flags.Var(&linter.config.SuppressNil, "suppress-nil-assertion", "Suppress warning for wrong nil assertions")
	a.Flags.Var(&linter.config.SuppressErr, "suppress-err-assertion", "Suppress warning for wrong error assertions")
	a.Flags.Var(&linter.config.SuppressCompare, "suppress-compare-assertion", "Suppress warning for wrong comparison assertions")
	a.Flags.Var(&linter.config.SuppressAsync, "suppress-async-assertion", "Suppress warning for function call in async assertion, like Eventually")
	a.Flags.Var(&linter.config.SuppressTypeCompare, "suppress-type-compare-assertion", "Suppress warning for comparing values from different types, like int32 and uint32")
	a.Flags.Var(&linter.config.AllowHaveLen0, "allow-havelen-0", "Do not warn for HaveLen(0); default = false")

	a.Flags.BoolVar(&ignored, "suppress-focus-container", true, "Suppress warning for ginkgo focus containers like FDescribe, FContext, FWhen or FIt. Deprecated and ignored: use --forbid-focus-container instead")
	a.Flags.Var(&linter.config.ForbidFocus, "forbid-focus-container", "trigger a warning for ginkgo focus containers like FDescribe, FContext, FWhen or FIt; default = false.")

	return a
}

const doc = `enforces standards of using ginkgo and gomega

or
       ginkgolinter version

version: %s

currently, the linter searches for following:
* trigger a warning when using Eventually or Constantly with a function call. This is in order to prevent the case when 
  using a function call instead of a function. Function call returns a value only once, and so the original value
  is tested again and again and is never changed. [Bug]

* trigger a warning when comparing a pointer to a value. [Bug]

* trigger a warning for missing assertion method: [Bug]
	Eventually(checkSomething)

* trigger a warning when a ginkgo focus container (FDescribe, FContext, FWhen or FIt) is found. [Bug]

* validate the MatchError gomega matcher [Bug]

* trigger a warning when using the Equal or the BeIdentical matcher with two different types, as these matchers will
  fail in runtime.

* wrong length assertions. We want to assert the item rather than its length. [Style]
For example:
	Expect(len(x)).Should(Equal(1))
This should be replaced with:
	Expect(x)).Should(HavelLen(1))
	
* wrong nil assertions. We want to assert the item rather than a comparison result. [Style]
For example:
	Expect(x == nil).Should(BeTrue())
This should be replaced with:
	Expect(x).Should(BeNil())

* wrong error assertions. For example: [Style]
	Expect(err == nil).Should(BeTrue())
This should be replaced with:
	Expect(err).ShouldNot(HaveOccurred())

* wrong boolean comparison, for example: [Style]
	Expect(x == 8).Should(BeTrue())
This should be replaced with:
	Expect(x).Should(BeEqual(8))

* replaces Equal(true/false) with BeTrue()/BeFalse() [Style]

* replaces HaveLen(0) with BeEmpty() [Style]
`

// main assertion function
func (l *ginkgoLinter) run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		fileConfig := l.config.Clone()

		cm := ast.NewCommentMap(pass.Fset, file, file.Comments)

		fileConfig.UpdateFromFile(cm)

		gomegaHndlr := gomegahandler.GetGomegaHandler(file)
		ginkgoHndlr := ginkgohandler.GetGinkgoHandler(file)

		if gomegaHndlr == nil && ginkgoHndlr == nil { // no gomega or ginkgo imports => no use in gomega in this file; nothing to do here
			continue
		}

		//gomegaMatcher = getMatcherInterface(pass, file)

		ast.Inspect(file, func(n ast.Node) bool {
			if ginkgoHndlr != nil && fileConfig.ForbidFocus {
				spec, ok := n.(*ast.ValueSpec)
				if ok {
					for _, val := range spec.Values {
						if exp, ok := val.(*ast.CallExpr); ok {
							if checkFocusContainer(pass, ginkgoHndlr, exp) {
								return true
							}
						}
					}
				}
			}

			stmt, ok := n.(*ast.ExprStmt)
			if !ok {
				return true
			}

			config := fileConfig.Clone()

			if comments, ok := cm[stmt]; ok {
				config.UpdateFromComment(comments)
			}

			// search for function calls
			assertionExp, ok := stmt.X.(*ast.CallExpr)
			if !ok {
				return true
			}

			if ginkgoHndlr != nil && bool(config.ForbidFocus) && checkFocusContainer(pass, ginkgoHndlr, assertionExp) {
				return true
			}

			// no more ginkgo checks. From here it's only gomega. So if there is no gomega handler, exit here. This is
			// mostly to prevent nil pointer error.
			if gomegaHndlr == nil {
				return true
			}

			assertionFunc, ok := assertionExp.Fun.(*ast.SelectorExpr)
			if !ok {
				checkNoAssertion(pass, assertionExp, gomegaHndlr)
				return true
			}

			if !isAssertionFunc(assertionFunc.Sel.Name) {
				checkNoAssertion(pass, assertionExp, gomegaHndlr)
				return true
			}

			actualExpr := gomegaHndlr.GetActualExpr(assertionFunc)
			if actualExpr == nil {
				return true
			}

			return checkExpression(pass, config, assertionExp, actualExpr, gomegaHndlr)
		})
	}
	return nil, nil
}

func checkFocusContainer(pass *analysis.Pass, ginkgoHndlr ginkgohandler.Handler, exp *ast.CallExpr) bool {
	foundFocus := false
	isFocus, id := ginkgoHndlr.GetFocusContainerName(exp)
	if isFocus {
		reportNewName(pass, id, id.Name[1:], focusContainerFound, id.Name)
		foundFocus = true
	}

	if id != nil && ginkgohandler.IsContainer(id) {
		for _, arg := range exp.Args {
			if ginkgoHndlr.IsFocusSpec(arg) {
				reportNoFix(pass, arg.Pos(), focusSpecFound)
				foundFocus = true
			} else if callExp, ok := arg.(*ast.CallExpr); ok {
				if checkFocusContainer(pass, ginkgoHndlr, callExp) { // handle table entries
					foundFocus = true
				}
			}
		}
	}

	return foundFocus
}

func checkExpression(pass *analysis.Pass, config types.Config, assertionExp *ast.CallExpr, actualExpr *ast.CallExpr, handler gomegahandler.Handler) bool {
	expr := astcopy.CallExpr(assertionExp)
	oldExpr := goFmt(pass.Fset, expr)

	if checkAsyncAssertion(pass, config, expr, actualExpr, handler, oldExpr) {
		return true
	}

	actualArg := getActualArg(actualExpr, handler)
	if actualArg == nil {
		return true
	}

	if !bool(config.SuppressLen) && isActualIsLenFunc(actualArg) {
		return checkLengthMatcher(expr, pass, handler, oldExpr)
	} else if nilable, compOp := getNilableFromComparison(actualArg); nilable != nil {
		if isExprError(pass, nilable) {
			if config.SuppressErr {
				return true
			}
		} else if config.SuppressNil {
			return true
		}

		return checkNilMatcher(expr, pass, nilable, handler, compOp == token.NEQ, oldExpr)

	} else if first, second, op, ok := isComparison(pass, actualArg); ok {
		matcher, shouldContinue := startCheckComparison(expr, handler)
		if !shouldContinue {
			return false
		}
		if !bool(config.SuppressLen) && isActualIsLenFunc(first) {
			if handleLenComparison(pass, expr, matcher, first, second, op, handler, oldExpr) {
				return false
			}
		}
		return bool(config.SuppressCompare) || checkComparison(expr, pass, matcher, handler, first, second, op, oldExpr)

	} else if checkMatchError(pass, assertionExp, actualArg, handler, oldExpr) {
		return false
	} else if isExprError(pass, actualArg) {
		return bool(config.SuppressErr) || checkNilError(pass, expr, handler, actualArg, oldExpr)

	} else if checkPointerComparison(pass, config, assertionExp, expr, actualArg, handler, oldExpr) {
		return false
	} else if !handleAssertionOnly(pass, config, expr, handler, actualArg, oldExpr, true) {
		return false
	} else if !config.SuppressTypeCompare {
		return !checkEqualWrongType(pass, assertionExp, actualArg, handler, oldExpr)
	}

	return true
}

func checkMatchError(pass *analysis.Pass, origExp *ast.CallExpr, actualArg ast.Expr, handler gomegahandler.Handler, oldExpr string) bool {
	matcher, ok := origExp.Args[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	return doCheckMatchError(pass, origExp, matcher, actualArg, handler, oldExpr)
}

func doCheckMatchError(pass *analysis.Pass, origExp *ast.CallExpr, matcher *ast.CallExpr, actualArg ast.Expr, handler gomegahandler.Handler, oldExpr string) bool {
	name, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return false
	}
	switch name {
	case matchError:
	case not:
		nested, ok := matcher.Args[0].(*ast.CallExpr)
		if !ok {
			return false
		}

		return doCheckMatchError(pass, origExp, nested, actualArg, handler, oldExpr)
	case and, or:
		res := true
		for _, arg := range matcher.Args {
			if nested, ok := arg.(*ast.CallExpr); ok {
				if !doCheckMatchError(pass, origExp, nested, actualArg, handler, oldExpr) {
					res = false
				}
			}
		}
		return res
	default:
		return false
	}

	if !isExprError(pass, actualArg) {
		reportNoFix(pass, origExp.Pos(), matchErrorArgWrongType, goFmt(pass.Fset, actualArg))
	}

	expr := astcopy.CallExpr(matcher)

	validAssertion, requiredParams := checkMatchErrorAssertion(pass, matcher)
	if !validAssertion {
		reportNoFix(pass, expr.Pos(), matchErrorWrongTypeAssertion, goFmt(pass.Fset, matcher.Args[0]))
		return false
	}

	numParams := len(matcher.Args)
	if numParams == requiredParams {
		if numParams == 2 {
			t := pass.TypesInfo.TypeOf(matcher.Args[1])
			if !gotypes.Identical(t, gotypes.Typ[gotypes.String]) {
				pass.Reportf(expr.Pos(), matchErrorNoFuncDescription)
				return false
			}
		}
		return true
	}

	if requiredParams == 2 && numParams == 1 {
		reportNoFix(pass, expr.Pos(), matchErrorMissingDescription)
		return false
	}

	var newArgsSuggestion = []ast.Expr{expr.Args[0]}
	if requiredParams == 2 {
		newArgsSuggestion = append(newArgsSuggestion, expr.Args[1])
	}
	expr.Args = newArgsSuggestion
	report(pass, expr, matchErrorRedundantArg, oldExpr)
	return false
}

func checkMatchErrorAssertion(pass *analysis.Pass, matcher *ast.CallExpr) (bool, int) {
	if isErrorMatcherValidArg(pass, matcher.Args[0]) {
		return true, 1
	}

	t1 := pass.TypesInfo.TypeOf(matcher.Args[0])
	if isFuncErrBool(t1) {
		return true, 2
	}

	return false, 0
}

// isFuncErrBool checks if a function is with the signature `func(error) bool`
func isFuncErrBool(t gotypes.Type) bool {
	sig, ok := t.(*gotypes.Signature)
	if !ok {
		return false
	}
	if sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}

	if !interfaces.ImplementsError(sig.Params().At(0).Type()) {
		return false
	}

	b, ok := sig.Results().At(0).Type().(*gotypes.Basic)
	if ok && b.Name() == "bool" && b.Info() == gotypes.IsBoolean && b.Kind() == gotypes.Bool {
		return true
	}

	return false
}

func isErrorMatcherValidArg(pass *analysis.Pass, arg ast.Expr) bool {
	if isExprError(pass, arg) {
		return true
	}

	if t, ok := pass.TypesInfo.TypeOf(arg).(*gotypes.Basic); ok && t.Kind() == gotypes.String {
		return true
	}

	t := pass.TypesInfo.TypeOf(arg)

	return interfaces.ImplementsGomegaMatcher(t)
}

func checkEqualWrongType(pass *analysis.Pass, origExp *ast.CallExpr, actualArg ast.Expr, handler gomegahandler.Handler, old string) bool {
	matcher, ok := origExp.Args[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	return checkEqualDifferentTypes(pass, matcher, actualArg, handler, old, false)
}

func checkEqualDifferentTypes(pass *analysis.Pass, matcher *ast.CallExpr, actualArg ast.Expr, handler gomegahandler.Handler, old string, parentPointer bool) bool {
	matcherFuncName, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return false
	}

	actualType := pass.TypesInfo.TypeOf(actualArg)

	switch matcherFuncName {
	case equal, beIdenticalTo: // continue
	case and, or:
		foundIssue := false
		for _, nestedExp := range matcher.Args {
			nested, ok := nestedExp.(*ast.CallExpr)
			if !ok {
				continue
			}
			if checkEqualDifferentTypes(pass, nested, actualArg, handler, old, parentPointer) {
				foundIssue = true
			}
		}

		return foundIssue
	case withTransform:
		nested, ok := matcher.Args[1].(*ast.CallExpr)
		if !ok {
			return false
		}

		matcherFuncName, ok = handler.GetActualFuncName(nested)
		switch matcherFuncName {
		case equal, beIdenticalTo:
		case not:
			return checkEqualDifferentTypes(pass, nested, actualArg, handler, old, parentPointer)
		default:
			return false
		}

		if t := getFuncType(pass, matcher.Args[0]); t != nil {
			actualType = t
			matcher = nested

			if !ok {
				return false
			}
		} else {
			return checkEqualDifferentTypes(pass, nested, actualArg, handler, old, parentPointer)
		}

	case not:
		nested, ok := matcher.Args[0].(*ast.CallExpr)
		if !ok {
			return false
		}

		return checkEqualDifferentTypes(pass, nested, actualArg, handler, old, parentPointer)

	case haveValue:
		nested, ok := matcher.Args[0].(*ast.CallExpr)
		if !ok {
			return false
		}

		return checkEqualDifferentTypes(pass, nested, actualArg, handler, old, true)
	default:
		return false
	}

	matcherValue := matcher.Args[0]

	switch act := actualType.(type) {
	case *gotypes.Tuple:
		actualType = act.At(0).Type()
	case *gotypes.Pointer:
		if parentPointer {
			actualType = act.Elem()
		}
	}

	matcherType := pass.TypesInfo.TypeOf(matcherValue)

	if !reflect.DeepEqual(matcherType, actualType) {
		// Equal can handle comparison of interface and a value that implements it
		if isImplementing(matcherType, actualType) || isImplementing(actualType, matcherType) {
			return false
		}

		reportNoFix(pass, matcher.Pos(), compareDifferentTypes, matcherFuncName, actualType, matcherType)
		return true
	}

	return false
}

func getFuncType(pass *analysis.Pass, expr ast.Expr) gotypes.Type {
	switch f := expr.(type) {
	case *ast.FuncLit:
		if f.Type != nil && f.Type.Results != nil && len(f.Type.Results.List) > 0 {
			return pass.TypesInfo.TypeOf(f.Type.Results.List[0].Type)
		}
	case *ast.Ident:
		a := pass.TypesInfo.TypeOf(f)
		if sig, ok := a.(*gotypes.Signature); ok && sig.Results().Len() > 0 {
			return sig.Results().At(0).Type()
		}
	}

	return nil
}

func isImplementing(ifs, impl gotypes.Type) bool {
	if gotypes.IsInterface(ifs) {

		var (
			theIfs *gotypes.Interface
			ok     bool
		)

		for {
			theIfs, ok = ifs.(*gotypes.Interface)
			if ok {
				break
			}
			ifs = ifs.Underlying()
		}

		return gotypes.Implements(impl, theIfs)
	}
	return false
}

// be careful - never change origExp!!! only modify its clone, expr!!!
func checkPointerComparison(pass *analysis.Pass, config types.Config, origExp *ast.CallExpr, expr *ast.CallExpr, actualArg ast.Expr, handler gomegahandler.Handler, oldExpr string) bool {
	if !isPointer(pass, actualArg) {
		return false
	}
	matcher, ok := origExp.Args[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	matcherFuncName, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return false
	}

	// not using recurse here, since we need the original expression, in order to get the TypeInfo, while we should not
	// modify it.
	for matcherFuncName == not {
		reverseAssertionFuncLogic(expr)
		expr.Args[0] = expr.Args[0].(*ast.CallExpr).Args[0]
		matcher, ok = matcher.Args[0].(*ast.CallExpr)
		if !ok {
			return false
		}

		matcherFuncName, ok = handler.GetActualFuncName(matcher)
		if !ok {
			return false
		}
	}

	switch matcherFuncName {
	case equal, beIdenticalTo, beEquivalentTo:
		arg := matcher.Args[0]
		if isPointer(pass, arg) {
			return false
		}
		if isNil(arg) {
			return false
		}
		if isInterface(pass, arg) {
			return false
		}
	case beFalse, beTrue, beNumerically:
	default:
		return false
	}

	handleAssertionOnly(pass, config, expr, handler, actualArg, oldExpr, false)

	args := []ast.Expr{astcopy.CallExpr(expr.Args[0].(*ast.CallExpr))}
	handler.ReplaceFunction(expr.Args[0].(*ast.CallExpr), ast.NewIdent(haveValue))
	expr.Args[0].(*ast.CallExpr).Args = args
	report(pass, expr, comparePointerToValue, oldExpr)

	return true

}

// check async assertion does not assert function call. This is a real bug in the test. In this case, the assertion is
// done on the returned value, instead of polling the result of a function, for instance.
func checkAsyncAssertion(pass *analysis.Pass, config types.Config, expr *ast.CallExpr, actualExpr *ast.CallExpr, handler gomegahandler.Handler, oldExpr string) bool {
	funcName, ok := handler.GetActualFuncName(actualExpr)
	if !ok {
		return false
	}

	var funcIndex int
	switch funcName {
	case eventually, consistently:
		funcIndex = 0
	case eventuallyWithOffset, consistentlyWithOffset:
		funcIndex = 1
	default:
		return false
	}

	if !config.SuppressAsync && len(actualExpr.Args) > funcIndex {
		t := pass.TypesInfo.TypeOf(actualExpr.Args[funcIndex])

		// skip context variable, if used as first argument
		if "context.Context" == t.String() {
			funcIndex++
		}

		if len(actualExpr.Args) > funcIndex {
			if fun, funcCall := actualExpr.Args[funcIndex].(*ast.CallExpr); funcCall {
				t = pass.TypesInfo.TypeOf(fun)
				if !isValidAsyncValueType(t) {
					actualExpr = handler.GetActualExpr(expr.Fun.(*ast.SelectorExpr))

					if len(fun.Args) > 0 {
						origArgs := actualExpr.Args
						origFunc := actualExpr.Fun
						actualExpr.Args = fun.Args

						origArgs[funcIndex] = fun.Fun
						call := &ast.SelectorExpr{
							Sel: ast.NewIdent("WithArguments"),
							X: &ast.CallExpr{
								Fun:  origFunc,
								Args: origArgs,
							},
						}

						actualExpr.Fun = call
						actualExpr.Args = fun.Args
					} else {
						actualExpr.Args[funcIndex] = fun.Fun
					}

					handleAssertionOnly(pass, config, expr, handler, actualExpr, oldExpr, false)
					report(pass, expr, fmt.Sprintf(valueInEventually, funcName, funcName)+"; consider using `%s` instead", oldExpr)
					return true
				}
			}
		}
	}

	handleAssertionOnly(pass, config, expr, handler, actualExpr, oldExpr, true)
	return true
}

func isValidAsyncValueType(t gotypes.Type) bool {
	switch t.(type) {
	// allow functions that return function or channel.
	case *gotypes.Signature, *gotypes.Chan, *gotypes.Pointer:
		return true
	case *gotypes.Named:
		return isValidAsyncValueType(t.Underlying())
	}

	return false
}

func startCheckComparison(exp *ast.CallExpr, handler gomegahandler.Handler) (*ast.CallExpr, bool) {
	matcher, ok := exp.Args[0].(*ast.CallExpr)
	if !ok {
		return nil, false
	}

	matcherFuncName, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return nil, false
	}

	switch matcherFuncName {
	case beTrue:
	case beFalse:
		reverseAssertionFuncLogic(exp)
	case equal:
		boolean, found := matcher.Args[0].(*ast.Ident)
		if !found {
			return nil, false
		}

		if boolean.Name == "false" {
			reverseAssertionFuncLogic(exp)
		} else if boolean.Name != "true" {
			return nil, false
		}

	case not:
		reverseAssertionFuncLogic(exp)
		exp.Args[0] = exp.Args[0].(*ast.CallExpr).Args[0]
		return startCheckComparison(exp, handler)

	default:
		return nil, false
	}

	return matcher, true
}

func checkComparison(exp *ast.CallExpr, pass *analysis.Pass, matcher *ast.CallExpr, handler gomegahandler.Handler, first ast.Expr, second ast.Expr, op token.Token, oldExp string) bool {
	fun, ok := exp.Fun.(*ast.SelectorExpr)
	if !ok {
		return true
	}

	call := handler.GetActualExpr(fun)
	if call == nil {
		return true
	}

	switch op {
	case token.EQL:
		handleEqualComparison(pass, matcher, first, second, handler)

	case token.NEQ:
		reverseAssertionFuncLogic(exp)
		handleEqualComparison(pass, matcher, first, second, handler)
	case token.GTR, token.GEQ, token.LSS, token.LEQ:
		if !isNumeric(pass, first) {
			return true
		}
		handler.ReplaceFunction(matcher, ast.NewIdent(beNumerically))
		matcher.Args = []ast.Expr{
			&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s"`, op.String())},
			second,
		}
	default:
		return true
	}

	call.Args = []ast.Expr{first}
	report(pass, exp, wrongCompareWarningTemplate, oldExp)
	return false
}

func handleEqualComparison(pass *analysis.Pass, matcher *ast.CallExpr, first ast.Expr, second ast.Expr, handler gomegahandler.Handler) {
	if isZero(pass, second) {
		handler.ReplaceFunction(matcher, ast.NewIdent(beZero))
		matcher.Args = nil
	} else {
		t := pass.TypesInfo.TypeOf(first)
		if gotypes.IsInterface(t) {
			handler.ReplaceFunction(matcher, ast.NewIdent(beIdenticalTo))
		} else if _, ok := t.(*gotypes.Pointer); ok {
			handler.ReplaceFunction(matcher, ast.NewIdent(beIdenticalTo))
		} else {
			handler.ReplaceFunction(matcher, ast.NewIdent(equal))
		}

		matcher.Args = []ast.Expr{second}
	}
}

func handleLenComparison(pass *analysis.Pass, exp *ast.CallExpr, matcher *ast.CallExpr, first ast.Expr, second ast.Expr, op token.Token, handler gomegahandler.Handler, oldExpr string) bool {
	switch op {
	case token.EQL:
	case token.NEQ:
		reverseAssertionFuncLogic(exp)
	default:
		return false
	}

	var eql *ast.Ident
	if isZero(pass, second) {
		eql = ast.NewIdent(beEmpty)
	} else {
		eql = ast.NewIdent(haveLen)
		matcher.Args = []ast.Expr{second}
	}

	handler.ReplaceFunction(matcher, eql)
	firstLen, ok := first.(*ast.CallExpr) // assuming it's len()
	if !ok {
		return false // should never happen
	}

	val := firstLen.Args[0]
	fun := handler.GetActualExpr(exp.Fun.(*ast.SelectorExpr))
	fun.Args = []ast.Expr{val}

	report(pass, exp, wrongLengthWarningTemplate, oldExpr)
	return true
}

// Check if the "actual" argument is a call to the golang built-in len() function
func isActualIsLenFunc(actualArg ast.Expr) bool {
	lenArgExp, ok := actualArg.(*ast.CallExpr)
	if !ok {
		return false
	}

	lenFunc, ok := lenArgExp.Fun.(*ast.Ident)
	return ok && lenFunc.Name == "len"
}

// Check if matcher function is in one of the patterns we want to avoid
func checkLengthMatcher(exp *ast.CallExpr, pass *analysis.Pass, handler gomegahandler.Handler, oldExp string) bool {
	matcher, ok := exp.Args[0].(*ast.CallExpr)
	if !ok {
		return true
	}

	matcherFuncName, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return true
	}

	switch matcherFuncName {
	case equal:
		handleEqualMatcher(matcher, pass, exp, handler, oldExp)
		return false

	case beZero:
		handleBeZero(pass, exp, handler, oldExp)
		return false

	case beNumerically:
		return handleBeNumerically(matcher, pass, exp, handler, oldExp)

	case not:
		reverseAssertionFuncLogic(exp)
		exp.Args[0] = exp.Args[0].(*ast.CallExpr).Args[0]
		return checkLengthMatcher(exp, pass, handler, oldExp)

	default:
		return true
	}
}

// Check if matcher function is in one of the patterns we want to avoid
func checkNilMatcher(exp *ast.CallExpr, pass *analysis.Pass, nilable ast.Expr, handler gomegahandler.Handler, notEqual bool, oldExp string) bool {
	matcher, ok := exp.Args[0].(*ast.CallExpr)
	if !ok {
		return true
	}

	matcherFuncName, ok := handler.GetActualFuncName(matcher)
	if !ok {
		return true
	}

	switch matcherFuncName {
	case equal:
		handleEqualNilMatcher(matcher, pass, exp, handler, nilable, notEqual, oldExp)

	case beTrue:
		handleNilBeBoolMatcher(pass, exp, handler, nilable, notEqual, oldExp)

	case beFalse:
		reverseAssertionFuncLogic(exp)
		handleNilBeBoolMatcher(pass, exp, handler, nilable, notEqual, oldExp)

	case not:
		reverseAssertionFuncLogic(exp)
		exp.Args[0] = exp.Args[0].(*ast.CallExpr).Args[0]
		return checkNilMatcher(exp, pass, nilable, handler, notEqual, oldExp)

	default:
		return true
	}
	return false
}

func checkNilError(pass *analysis.Pass, assertionExp *ast.CallExpr, handler gomegahandler.Handler, actualArg ast.Expr, oldExpr string) bool {
	if len(assertionExp.Args) == 0 {
		return true
	}

	equalFuncExpr, ok := assertionExp.Args[0].(*ast.CallExpr)
	if !ok {
		return true
	}

	funcName, ok := handler.GetActualFuncName(equalFuncExpr)
	if !ok {
		return true
	}

	switch funcName {
	case beNil: // no additional processing needed.
	case equal:

		if len(equalFuncExpr.Args) == 0 {
			return true
		}

		nilable, ok := equalFuncExpr.Args[0].(*ast.Ident)
		if !ok || nilable.Name != "nil" {
			return true
		}

	case not:
		reverseAssertionFuncLogic(assertionExp)
		assertionExp.Args[0] = assertionExp.Args[0].(*ast.CallExpr).Args[0]
		return checkNilError(pass, assertionExp, handler, actualArg, oldExpr)
	default:
		return true
	}

	var newFuncName string
	if _, ok := actualArg.(*ast.CallExpr); ok {
		newFuncName = succeed
	} else {
		reverseAssertionFuncLogic(assertionExp)
		newFuncName = haveOccurred
	}

	handler.ReplaceFunction(equalFuncExpr, ast.NewIdent(newFuncName))
	equalFuncExpr.Args = nil

	report(pass, assertionExp, wrongErrWarningTemplate, oldExpr)
	return false
}

// handleAssertionOnly checks use-cases when the actual value is valid, but only the assertion should be fixed
// it handles:
//
//	Equal(nil) => BeNil()
//	Equal(true) => BeTrue()
//	Equal(false) => BeFalse()
//	HaveLen(0) => BeEmpty()
func handleAssertionOnly(pass *analysis.Pass, config types.Config, expr *ast.CallExpr, handler gomegahandler.Handler, actualArg ast.Expr, oldExpr string, shouldReport bool) bool {
	if len(expr.Args) == 0 {
		return true
	}

	equalFuncExpr, ok := expr.Args[0].(*ast.CallExpr)
	if !ok {
		return true
	}

	funcName, ok := handler.GetActualFuncName(equalFuncExpr)
	if !ok {
		return true
	}

	switch funcName {
	case equal:
		if len(equalFuncExpr.Args) == 0 {
			return true
		}

		tkn, ok := equalFuncExpr.Args[0].(*ast.Ident)
		if !ok {
			return true
		}

		var replacement string
		var template string
		switch tkn.Name {
		case "nil":
			if config.SuppressNil {
				return true
			}
			replacement = beNil
			template = wrongNilWarningTemplate
		case "true":
			replacement = beTrue
			template = wrongBoolWarningTemplate
		case "false":
			if isNegativeAssertion(expr) {
				reverseAssertionFuncLogic(expr)
				replacement = beTrue
			} else {
				replacement = beFalse
			}
			template = wrongBoolWarningTemplate
		default:
			return true
		}

		handler.ReplaceFunction(equalFuncExpr, ast.NewIdent(replacement))
		equalFuncExpr.Args = nil

		if shouldReport {
			report(pass, expr, template, oldExpr)
		}

		return false

	case beFalse:
		if isNegativeAssertion(expr) {
			reverseAssertionFuncLogic(expr)
			handler.ReplaceFunction(equalFuncExpr, ast.NewIdent(beTrue))
			if shouldReport {
				report(pass, expr, doubleNegativeWarningTemplate, oldExpr)
			}
		}
		return false

	case haveLen:
		if config.AllowHaveLen0 {
			return true
		}

		if len(equalFuncExpr.Args) > 0 {
			if isZero(pass, equalFuncExpr.Args[0]) {
				handler.ReplaceFunction(equalFuncExpr, ast.NewIdent(beEmpty))
				equalFuncExpr.Args = nil
				if shouldReport {
					report(pass, expr, wrongLengthWarningTemplate, oldExpr)
				}
				return false
			}
		}

		return true

	case not:
		reverseAssertionFuncLogic(expr)
		expr.Args[0] = expr.Args[0].(*ast.CallExpr).Args[0]
		return handleAssertionOnly(pass, config, expr, handler, actualArg, oldExpr, shouldReport)
	default:
		return true
	}
}

func isZero(pass *analysis.Pass, arg ast.Expr) bool {
	if val, ok := arg.(*ast.BasicLit); ok && val.Kind == token.INT && val.Value == "0" {
		return true
	}
	info, ok := pass.TypesInfo.Types[arg]
	if ok {
		if t, ok := info.Type.(*gotypes.Basic); ok && t.Kind() == gotypes.Int && info.Value != nil {
			if i, ok := constant.Int64Val(info.Value); ok && i == 0 {
				return true
			}
		}
	} else if val, ok := arg.(*ast.Ident); ok && val.Obj != nil && val.Obj.Kind == ast.Con {
		if spec, ok := val.Obj.Decl.(*ast.ValueSpec); ok {
			if len(spec.Values) == 1 {
				if value, ok := spec.Values[0].(*ast.BasicLit); ok && value.Kind == token.INT && value.Value == "0" {
					return true
				}
			}
		}
	}

	return false
}

// getActualArg checks that the function is an assertion's actual function and return the "actual" parameter. If the
// function is not assertion's actual function, return nil.
func getActualArg(actualExpr *ast.CallExpr, handler gomegahandler.Handler) ast.Expr {
	funcName, ok := handler.GetActualFuncName(actualExpr)
	if !ok {
		return nil
	}

	switch funcName {
	case expect, omega:
		return actualExpr.Args[0]
	case expectWithOffset:
		return actualExpr.Args[1]
	default:
		return nil
	}
}

// Replace the len function call by its parameter, to create a fix suggestion
func replaceLenActualArg(actualExpr *ast.CallExpr, handler gomegahandler.Handler) {
	name, ok := handler.GetActualFuncName(actualExpr)
	if !ok {
		return
	}

	switch name {
	case expect, omega:
		arg := actualExpr.Args[0]
		if isActualIsLenFunc(arg) {
			// replace the len function call by its parameter, to create a fix suggestion
			actualExpr.Args[0] = arg.(*ast.CallExpr).Args[0]
		}
	case expectWithOffset:
		arg := actualExpr.Args[1]
		if isActualIsLenFunc(arg) {
			// replace the len function call by its parameter, to create a fix suggestion
			actualExpr.Args[1] = arg.(*ast.CallExpr).Args[0]
		}
	}
}

// Replace the nil comparison with the compared object, to create a fix suggestion
func replaceNilActualArg(actualExpr *ast.CallExpr, handler gomegahandler.Handler, nilable ast.Expr) bool {
	actualFuncName, ok := handler.GetActualFuncName(actualExpr)
	if !ok {
		return false
	}

	switch actualFuncName {
	case expect, omega:
		actualExpr.Args[0] = nilable
		return true

	case expectWithOffset:
		actualExpr.Args[1] = nilable
		return true

	default:
		return false
	}
}

// For the BeNumerically matcher, we want to avoid the assertion of length to be > 0 or >= 1, or just == number
func handleBeNumerically(matcher *ast.CallExpr, pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, oldExp string) bool {
	opExp, ok1 := matcher.Args[0].(*ast.BasicLit)
	valExp, ok2 := matcher.Args[1].(*ast.BasicLit)

	if ok1 && ok2 {
		op := opExp.Value
		val := valExp.Value

		if (op == `">"` && val == "0") || (op == `">="` && val == "1") {
			reverseAssertionFuncLogic(exp)
			handler.ReplaceFunction(exp.Args[0].(*ast.CallExpr), ast.NewIdent(beEmpty))
			exp.Args[0].(*ast.CallExpr).Args = nil
			reportLengthAssertion(pass, exp, handler, oldExp)
			return false
		} else if op == `"=="` {
			chooseNumericMatcher(pass, exp, handler, valExp)
			reportLengthAssertion(pass, exp, handler, oldExp)

			return false
		} else if op == `"!="` {
			reverseAssertionFuncLogic(exp)
			chooseNumericMatcher(pass, exp, handler, valExp)
			reportLengthAssertion(pass, exp, handler, oldExp)

			return false
		}
	}
	return true
}

func chooseNumericMatcher(pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, valExp ast.Expr) {
	caller := exp.Args[0].(*ast.CallExpr)
	if isZero(pass, valExp) {
		handler.ReplaceFunction(caller, ast.NewIdent(beEmpty))
		exp.Args[0].(*ast.CallExpr).Args = nil
	} else {
		handler.ReplaceFunction(caller, ast.NewIdent(haveLen))
		exp.Args[0].(*ast.CallExpr).Args = []ast.Expr{valExp}
	}
}

func reverseAssertionFuncLogic(exp *ast.CallExpr) {
	assertionFunc := exp.Fun.(*ast.SelectorExpr).Sel
	assertionFunc.Name = reverseassertion.ChangeAssertionLogic(assertionFunc.Name)
}

func isNegativeAssertion(exp *ast.CallExpr) bool {
	assertionFunc := exp.Fun.(*ast.SelectorExpr).Sel
	return reverseassertion.IsNegativeLogic(assertionFunc.Name)
}

func handleEqualMatcher(matcher *ast.CallExpr, pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, oldExp string) {
	equalTo, ok := matcher.Args[0].(*ast.BasicLit)
	if ok {
		chooseNumericMatcher(pass, exp, handler, equalTo)
	} else {
		handler.ReplaceFunction(exp.Args[0].(*ast.CallExpr), ast.NewIdent(haveLen))
		exp.Args[0].(*ast.CallExpr).Args = []ast.Expr{matcher.Args[0]}
	}
	reportLengthAssertion(pass, exp, handler, oldExp)
}

func handleBeZero(pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, oldExp string) {
	exp.Args[0].(*ast.CallExpr).Args = nil
	handler.ReplaceFunction(exp.Args[0].(*ast.CallExpr), ast.NewIdent(beEmpty))
	reportLengthAssertion(pass, exp, handler, oldExp)
}

func handleEqualNilMatcher(matcher *ast.CallExpr, pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, nilable ast.Expr, notEqual bool, oldExp string) {
	equalTo, ok := matcher.Args[0].(*ast.Ident)
	if !ok {
		return
	}

	if equalTo.Name == "false" {
		reverseAssertionFuncLogic(exp)
	} else if equalTo.Name != "true" {
		return
	}

	newFuncName, isItError := handleNilComparisonErr(pass, exp, nilable)

	handler.ReplaceFunction(exp.Args[0].(*ast.CallExpr), ast.NewIdent(newFuncName))
	exp.Args[0].(*ast.CallExpr).Args = nil

	reportNilAssertion(pass, exp, handler, nilable, notEqual, oldExp, isItError)
}

func handleNilBeBoolMatcher(pass *analysis.Pass, exp *ast.CallExpr, handler gomegahandler.Handler, nilable ast.Expr, notEqual bool, oldExp string) {
	newFuncName, isItError := handleNilComparisonErr(pass, exp, nilable)
	handler.ReplaceFunction(exp.Args[0].(*ast.CallExpr), ast.NewIdent(newFuncName))
	exp.Args[0].(*ast.CallExpr).Args = nil

	reportNilAssertion(pass, exp, handler, nilable, notEqual, oldExp, isItError)
}

func handleNilComparisonErr(pass *analysis.Pass, exp *ast.CallExpr, nilable ast.Expr) (string, bool) {
	newFuncName := beNil
	isItError := isExprError(pass, nilable)
	if isItError {
		if _, ok := nilable.(*ast.CallExpr); ok {
			newFuncName = succeed
		} else {
			reverseAssertionFuncLogic(exp)
			newFuncName = haveOccurred
		}
	}

	return newFuncName, isItError
}

func isAssertionFunc(name string) bool {
	switch name {
	case "To", "ToNot", "NotTo", "Should", "ShouldNot":
		return true
	}
	return false
}

func reportLengthAssertion(pass *analysis.Pass, expr *ast.CallExpr, handler gomegahandler.Handler, oldExpr string) {
	actualExpr := handler.GetActualExpr(expr.Fun.(*ast.SelectorExpr))
	replaceLenActualArg(actualExpr, handler)

	report(pass, expr, wrongLengthWarningTemplate, oldExpr)
}

func reportNilAssertion(pass *analysis.Pass, expr *ast.CallExpr, handler gomegahandler.Handler, nilable ast.Expr, notEqual bool, oldExpr string, isItError bool) {
	actualExpr := handler.GetActualExpr(expr.Fun.(*ast.SelectorExpr))
	changed := replaceNilActualArg(actualExpr, handler, nilable)
	if !changed {
		return
	}

	if notEqual {
		reverseAssertionFuncLogic(expr)
	}
	template := wrongNilWarningTemplate
	if isItError {
		template = wrongErrWarningTemplate
	}

	report(pass, expr, template, oldExpr)
}

func report(pass *analysis.Pass, expr ast.Expr, messageTemplate, oldExpr string) {
	newExp := goFmt(pass.Fset, expr)
	pass.Report(analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: fmt.Sprintf(messageTemplate, newExp),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should replace %s with %s", oldExpr, newExp),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     expr.Pos(),
						End:     expr.End(),
						NewText: []byte(newExp),
					},
				},
			},
		},
	})
}

func reportNewName(pass *analysis.Pass, id *ast.Ident, newName string, messageTemplate, oldExpr string) {
	pass.Report(analysis.Diagnostic{
		Pos:     id.Pos(),
		Message: fmt.Sprintf(messageTemplate, newName),
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fmt.Sprintf("should replace %s with %s", oldExpr, newName),
				TextEdits: []analysis.TextEdit{
					{
						Pos:     id.Pos(),
						End:     id.End(),
						NewText: []byte(newName),
					},
				},
			},
		},
	})
}

func reportNoFix(pass *analysis.Pass, pos token.Pos, message string, args ...any) {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}

	pass.Report(analysis.Diagnostic{
		Pos:     pos,
		Message: message,
	})
}

func getNilableFromComparison(actualArg ast.Expr) (ast.Expr, token.Token) {
	bin, ok := actualArg.(*ast.BinaryExpr)
	if !ok {
		return nil, token.ILLEGAL
	}

	if bin.Op == token.EQL || bin.Op == token.NEQ {
		if isNil(bin.Y) {
			return bin.X, bin.Op
		} else if isNil(bin.X) {
			return bin.Y, bin.Op
		}
	}

	return nil, token.ILLEGAL
}

func isNil(expr ast.Expr) bool {
	nilObject, ok := expr.(*ast.Ident)
	return ok && nilObject.Name == "nil" && nilObject.Obj == nil
}

func isComparison(pass *analysis.Pass, actualArg ast.Expr) (ast.Expr, ast.Expr, token.Token, bool) {
	bin, ok := actualArg.(*ast.BinaryExpr)
	if !ok {
		return nil, nil, token.ILLEGAL, false
	}

	first, second, op := bin.X, bin.Y, bin.Op
	replace := false
	switch realFirst := first.(type) {
	case *ast.Ident: // check if const
		info, ok := pass.TypesInfo.Types[realFirst]
		if ok {
			if _, ok := info.Type.(*gotypes.Basic); ok && info.Value != nil {
				replace = true
			}
		}

	case *ast.BasicLit:
		replace = true
	}

	if replace {
		first, second = second, first
	}

	switch op {
	case token.EQL:
	case token.NEQ:
	case token.GTR, token.GEQ, token.LSS, token.LEQ:
		if replace {
			op = reverseassertion.ChangeCompareOperator(op)
		}
	default:
		return nil, nil, token.ILLEGAL, false
	}
	return first, second, op, true
}

func goFmt(fset *token.FileSet, x ast.Expr) string {
	var b bytes.Buffer
	_ = printer.Fprint(&b, fset, x)
	return b.String()
}

func isExprError(pass *analysis.Pass, expr ast.Expr) bool {
	actualArgType := pass.TypesInfo.TypeOf(expr)
	switch t := actualArgType.(type) {
	case *gotypes.Named:
		if interfaces.ImplementsError(actualArgType) {
			return true
		}
	case *gotypes.Tuple:
		if t.Len() > 0 {
			switch t0 := t.At(0).Type().(type) {
			case *gotypes.Named, *gotypes.Pointer:
				if interfaces.ImplementsError(t0) {
					return true
				}
			}
		}
	}
	return false
}

func isPointer(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	_, ok := t.(*gotypes.Pointer)
	return ok
}

func isInterface(pass *analysis.Pass, expr ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(expr)
	return gotypes.IsInterface(t)
}

func isNumeric(pass *analysis.Pass, node ast.Expr) bool {
	t := pass.TypesInfo.TypeOf(node)

	switch t.String() {
	case "int", "uint", "int8", "uint8", "int16", "uint16", "int32", "uint32", "int64", "uint64", "float32", "float64":
		return true
	}
	return false
}

func checkNoAssertion(pass *analysis.Pass, expr *ast.CallExpr, handler gomegahandler.Handler) {
	funcName, ok := handler.GetActualFuncName(expr)
	if ok {
		if isActualFunc(funcName) {
			reportNoFix(pass, expr.Pos(), missingAssertionMessage, funcName)
		} else if isActualAsyncFunc(funcName) {
			reportNoFix(pass, expr.Pos(), missingAsyncAssertionMessage, funcName)
		}
	}
}

func isActualFunc(name string) bool {
	return name == expect || name == expectWithOffset
}

func isActualAsyncFunc(name string) bool {
	switch name {
	case eventually, eventuallyWithOffset, consistently, consistentlyWithOffset:
		return true
	}
	return false
}
