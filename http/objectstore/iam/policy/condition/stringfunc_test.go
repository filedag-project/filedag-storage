package condition

import (
	"encoding/base64"
	"reflect"
	"testing"

	"github.com/filedag-project/filedag-storage/http/objectstore/iam/set"
)

func TestStringEqualsFuncEvaluate(t *testing.T) {
	case1Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringEqualsFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("eu-west-1"), NewStringValue("ap-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}
	case3Function, err := newStringEqualsFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue(S3LocationConstraint.ToKey().VarName())), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case4Function, err := newStringEqualsFunc(NewKey(ExistingObjectTag, "security"), NewValueSet(NewStringValue("public")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{}, false},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{"test2", case2Function, map[string][]string{}, false},

		{"test3", case3Function, map[string][]string{"LocationConstraint": {"us-west-1"}}, true},

		{"test4", case4Function, map[string][]string{"ExistingObjectTag/security": {"public"}}, true},
		{"test4", case4Function, map[string][]string{"ExistingObjectTag/security": {"private"}}, false},
		{"test4", case4Function, map[string][]string{"ExistingObjectTag/project": {"foo"}}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringNotEqualsFuncEvaluate(t *testing.T) {
	case1Function, err := newStringNotEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringNotEqualsFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("eu-west-1"), NewStringValue("ap-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{}, true},
		{"test1", case1Function, map[string][]string{"delimiter": {"/"}}, true},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{"test2", case2Function, map[string][]string{}, true},
		{"test2", case2Function, map[string][]string{"delimiter": {"/"}}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringEqualsIgnoreCaseFuncEvaluate(t *testing.T) {
	case1Function, err := newStringEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringEqualsIgnoreCaseFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("EU-WEST-1"), NewStringValue("AP-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{}, false},
		{"test1", case1Function, map[string][]string{"delimiter": {"/"}}, false},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{"test2", case2Function, map[string][]string{}, false},
		{"test2", case2Function, map[string][]string{"delimiter": {"/"}}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringNotEqualsIgnoreCaseFuncEvaluate(t *testing.T) {
	case1Function, err := newStringNotEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringNotEqualsIgnoreCaseFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("EU-WEST-1"), NewStringValue("AP-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{}, true},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{"test2", case2Function, map[string][]string{}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestBinaryEqualsFuncEvaluate(t *testing.T) {
	case1Function, err := newBinaryEqualsFunc(
		S3XAmzCopySource.ToKey(),
		NewValueSet(NewStringValue(base64.StdEncoding.EncodeToString([]byte("mybucket/myobject")))),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newBinaryEqualsFunc(
		S3LocationConstraint.ToKey(),
		NewValueSet(
			NewStringValue(base64.StdEncoding.EncodeToString([]byte("eu-west-1"))),
			NewStringValue(base64.StdEncoding.EncodeToString([]byte("ap-southeast-1"))),
		),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{}, false},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{"test2", case2Function, map[string][]string{}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringLikeFuncEvaluate(t *testing.T) {
	case1Function, err := newStringLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringLikeFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("eu-west-*"), NewStringValue("ap-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{}, false},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-2"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{"test2", case2Function, map[string][]string{}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringNotLikeFuncEvaluate(t *testing.T) {
	case1Function, err := newStringNotLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringNotLikeFunc(S3LocationConstraint.ToKey(), NewValueSet(NewStringValue("eu-west-*"), NewStringValue("ap-southeast-1")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{"test1", case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{"test1", case1Function, map[string][]string{}, true},

		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"eu-west-2"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{"test2", case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{"test2", case2Function, map[string][]string{}, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.evaluate(testCase.values)

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringFuncKey(t *testing.T) {
	case1Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		function       CondFunction
		expectedResult Key
	}{
		{case1Function, S3XAmzCopySource.ToKey()},
	}

	for i, testCase := range testCases {
		result := testCase.function.key()

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
	}
}

func TestStringFuncName(t *testing.T) {
	case1Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Function, err := newStringNotEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case3Function, err := newStringEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case4Function, err := newStringNotEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case5Function, err := newBinaryEqualsFunc(
		S3XAmzCopySource.ToKey(),
		NewValueSet(NewStringValue(base64.StdEncoding.EncodeToString([]byte("mybucket/myobject")))),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case6Function, err := newStringLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case7Function, err := newStringNotLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	testCases := []struct {
		testname       string
		function       CondFunction
		expectedResult name
	}{
		{"test1", case1Function, name{name: stringEquals}},
		{"test2", case2Function, name{name: stringNotEquals}},
		{"test3", case3Function, name{name: stringEqualsIgnoreCase}},
		{"test4", case4Function, name{name: stringNotEqualsIgnoreCase}},
		{"test5", case5Function, name{name: binaryEquals}},
		{"test6", case6Function, name{name: stringLike}},
		{"test7", case7Function, name{name: stringNotLike}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.name()

			if result != testCase.expectedResult {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringEqualsFuncToMap(t *testing.T) {
	case1Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case1Result := map[Key]ValueSet{
		S3XAmzCopySource.ToKey(): NewValueSet(NewStringValue("mybucket/myobject")),
	}

	case2Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(),
		NewValueSet(
			NewStringValue("mybucket/myobject"),
			NewStringValue("yourbucket/myobject"),
		),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Result := map[Key]ValueSet{
		S3XAmzCopySource.ToKey(): NewValueSet(
			NewStringValue("mybucket/myobject"),
			NewStringValue("yourbucket/myobject"),
		),
	}

	case3Function, err := newStringEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case3Result := map[Key]ValueSet{
		S3XAmzCopySource.ToKey(): NewValueSet(NewStringValue("mybucket/MYOBJECT")),
	}

	case4Function, err := newBinaryEqualsFunc(
		S3XAmzCopySource.ToKey(),
		NewValueSet(NewStringValue(base64.StdEncoding.EncodeToString([]byte("mybucket/myobject")))),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case4Result := map[Key]ValueSet{
		S3XAmzCopySource.ToKey(): NewValueSet(NewStringValue(base64.StdEncoding.EncodeToString([]byte("mybucket/myobject")))),
	}

	testCases := []struct {
		testname       string
		f              CondFunction
		expectedResult map[Key]ValueSet
	}{
		{"test1", case1Function, case1Result},
		{"test2", case2Function, case2Result},
		{"test3", case3Function, case3Result},
		{"test4", case4Function, case4Result},
		{"test5", &stringFunc{}, nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.f.toMap()

			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("case %v: result: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}

func TestStringFuncClone(t *testing.T) {
	case1Function, err := newStringEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case1Result := &stringFunc{
		n:          name{name: stringEquals},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/myobject"),
		ignoreCase: false,
		base64:     false,
		negate:     false,
	}

	case2Function, err := newStringNotEqualsFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case2Result := &stringFunc{
		n:          name{name: stringNotEquals},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/myobject"),
		ignoreCase: false,
		base64:     false,
		negate:     true,
	}

	case3Function, err := newStringEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case3Result := &stringFunc{
		n:          name{name: stringEqualsIgnoreCase},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/MYOBJECT"),
		ignoreCase: true,
		base64:     false,
		negate:     false,
	}

	case4Function, err := newStringNotEqualsIgnoreCaseFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/MYOBJECT")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case4Result := &stringFunc{
		n:          name{name: stringNotEqualsIgnoreCase},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/MYOBJECT"),
		ignoreCase: true,
		base64:     false,
		negate:     true,
	}

	case5Function, err := newBinaryEqualsFunc(
		S3XAmzCopySource.ToKey(),
		NewValueSet(NewStringValue(base64.StdEncoding.EncodeToString([]byte("mybucket/myobject")))),
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case5Result := &stringFunc{
		n:          name{name: binaryEquals},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/myobject"),
		ignoreCase: false,
		base64:     true,
		negate:     false,
	}

	case6Function, err := newStringLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case6Result := &stringLikeFunc{stringFunc{
		n:          name{name: stringLike},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/myobject"),
		ignoreCase: false,
		base64:     false,
		negate:     false,
	}}

	case7Function, err := newStringNotLikeFunc(S3XAmzCopySource.ToKey(), NewValueSet(NewStringValue("mybucket/myobject")), "")
	if err != nil {
		t.Fatalf("unexpected error. %v\n", err)
	}

	case7Result := &stringLikeFunc{stringFunc{
		n:          name{name: stringNotLike},
		k:          S3XAmzCopySource.ToKey(),
		values:     set.CreateStringSet("mybucket/myobject"),
		ignoreCase: false,
		base64:     false,
		negate:     true,
	}}

	testCases := []struct {
		testname       string
		function       CondFunction
		expectedResult CondFunction
	}{
		{"test1", case1Function, case1Result},
		{"test2", case2Function, case2Result},
		{"test3", case3Function, case3Result},
		{"test4", case4Function, case4Result},
		{"test5", case5Function, case5Result},
		{"test6", case6Function, case6Result},
		{"test7", case7Function, case7Result},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testname, func(t *testing.T) {
			result := testCase.function.clone()

			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Fatalf("case %v: expected: %v, got: %v\n", testCase.testname, testCase.expectedResult, result)
			}
		})
	}
}
