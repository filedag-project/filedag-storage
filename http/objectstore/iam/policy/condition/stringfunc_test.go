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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{case1Function, map[string][]string{}, false},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{case2Function, map[string][]string{}, false},

		{case3Function, map[string][]string{"LocationConstraint": {"us-west-1"}}, true},

		{case4Function, map[string][]string{"ExistingObjectTag/security": {"public"}}, true},
		{case4Function, map[string][]string{"ExistingObjectTag/security": {"private"}}, false},
		{case4Function, map[string][]string{"ExistingObjectTag/project": {"foo"}}, false},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{case1Function, map[string][]string{}, true},
		{case1Function, map[string][]string{"delimiter": {"/"}}, true},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{case2Function, map[string][]string{}, true},
		{case2Function, map[string][]string{"delimiter": {"/"}}, true},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{case1Function, map[string][]string{}, false},
		{case1Function, map[string][]string{"delimiter": {"/"}}, false},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{case2Function, map[string][]string{}, false},
		{case2Function, map[string][]string{"delimiter": {"/"}}, false},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{case1Function, map[string][]string{}, true},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{case2Function, map[string][]string{}, true},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{case1Function, map[string][]string{}, false},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{case2Function, map[string][]string{}, false},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, true},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, false},
		{case1Function, map[string][]string{}, false},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-2"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, true},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, false},
		{case2Function, map[string][]string{}, false},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		values         map[string][]string
		expectedResult bool
	}{
		{case1Function, map[string][]string{"x-amz-copy-source": {"mybucket/myobject"}}, false},
		{case1Function, map[string][]string{"x-amz-copy-source": {"yourbucket/myobject"}}, true},
		{case1Function, map[string][]string{}, true},

		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"eu-west-2"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"ap-southeast-1"}}, false},
		{case2Function, map[string][]string{"LocationConstraint": {"us-east-1"}}, true},
		{case2Function, map[string][]string{}, true},
	}

	for i, testCase := range testCases {
		result := testCase.function.evaluate(testCase.values)

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		expectedResult name
	}{
		{case1Function, name{name: stringEquals}},
		{case2Function, name{name: stringNotEquals}},
		{case3Function, name{name: stringEqualsIgnoreCase}},
		{case4Function, name{name: stringNotEqualsIgnoreCase}},
		{case5Function, name{name: binaryEquals}},
		{case6Function, name{name: stringLike}},
		{case7Function, name{name: stringNotLike}},
	}

	for i, testCase := range testCases {
		result := testCase.function.name()

		if result != testCase.expectedResult {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		f              CondFunction
		expectedResult map[Key]ValueSet
	}{
		{case1Function, case1Result},
		{case2Function, case2Result},
		{case3Function, case3Result},
		{case4Function, case4Result},
		{&stringFunc{}, nil},
	}

	for i, testCase := range testCases {
		result := testCase.f.toMap()

		if !reflect.DeepEqual(result, testCase.expectedResult) {
			t.Fatalf("case %v: result: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
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
		function       CondFunction
		expectedResult CondFunction
	}{
		{case1Function, case1Result},
		{case2Function, case2Result},
		{case3Function, case3Result},
		{case4Function, case4Result},
		{case5Function, case5Result},
		{case6Function, case6Result},
		{case7Function, case7Result},
	}

	for i, testCase := range testCases {
		result := testCase.function.clone()

		if !reflect.DeepEqual(result, testCase.expectedResult) {
			t.Fatalf("case %v: expected: %v, got: %v\n", i+1, testCase.expectedResult, result)
		}
	}
}
