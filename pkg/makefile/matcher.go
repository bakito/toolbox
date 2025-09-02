package makefile

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/format"
	types2 "github.com/onsi/gomega/types"
	"github.com/pmezard/go-difflib/difflib"
)

func EqualDiff(expected string) types2.GomegaMatcher {
	return &EqualDiffMatcher{
		Expected: expected,
	}
}

type EqualDiffMatcher struct {
	Expected any
	diff     string
}

func (matcher *EqualDiffMatcher) Match(actual any) (success bool, err error) {
	if actual == nil && matcher.Expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized.")
	}
	if actualByteSlice, ok := actual.([]byte); ok {
		if expectedByteSlice, ok := matcher.Expected.([]byte); ok {
			diff, err := unifiedDiff(string(actualByteSlice), string(expectedByteSlice))
			if err != nil {
				return false, err
			}
			matcher.diff = diff
			return true, nil
		}
	}
	if actualString, ok := actual.(string); ok {
		if expectedString, ok := matcher.Expected.(string); ok {
			diff, err := unifiedDiff(actualString, expectedString)
			if err != nil {
				return false, err
			}
			matcher.diff = diff
			return true, nil
		}
	}
	return false, fmt.Errorf("expected %s to be of type string or []byte", reflect.TypeOf(actual))

}

func (matcher *EqualDiffMatcher) FailureMessage(actual any) (message string) {
	actualString, actualOK := actual.(string)
	expectedString, expectedOK := matcher.Expected.(string)
	if actualOK && expectedOK {
		return format.MessageWithDiff(actualString, "to equal", expectedString)
	}

	return format.Message(actual, "to equal", matcher.Expected)
}

func (matcher *EqualDiffMatcher) NegatedFailureMessage(actual any) (message string) {
	return format.Message(actual, "not to equal", matcher.Expected)
}

func unifiedDiff(a, b string) (string, error) {
	ud := difflib.UnifiedDiff{
		A:        difflib.SplitLines(a),
		B:        difflib.SplitLines(b),
		FromFile: "Expected",
		ToFile:   "Current",
		Context:  3, // lines of context, tweak as needed
	}
	return difflib.GetUnifiedDiffString(ud)
}
