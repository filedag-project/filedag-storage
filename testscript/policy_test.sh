. ./common.sh

function test_policy_set() {
  echo 'test policy set'
  # creat bucket for test
  mc mb loc1/testapolicyset >/dev/null 2>&1
  mc mb loc1/testadminpolicyset >/dev/null 2>&1
  mc policy set public loc1/testadminpolicyset >/dev/null 2>&1
  mc cp test_cp.txt loc1/testapolicyset >/dev/null 2>&1 # text
  mc cp test_cp.txt loc1/testadminpolicyset >/dev/null 2>&1 # text
  echo ' 1)user testA set bucket policy'
  mc policy set public loc1/testapolicyset >/dev/null 2>&1
  test_except "test set public policy" 0 $?
  mc policy set upload loc1/testapolicyset >/dev/null 2>&1
  test_except "test set upload policy" 0 $?
  mc policy set download loc1/testapolicyset >/dev/null 2>&1
  test_except "test set download policy" 0 $?
  mc policy set none loc1/testapolicyset >/dev/null 2>&1
  test_except "test set none policy" 0 $?
  mc policy set public loc2/testapolicyset >/dev/null 2>&1
  test_except "test no permission set public policy" 1 $?
  mc policy set upload loc2/testapolicyset >/dev/null 2>&1
  test_except "test no permission set upload policy" 1 $?
  mc policy set download loc2/testapolicyset >/dev/null 2>&1
  test_except "test no permission set download policy" 1 $?
  mc policy set none loc2/testapolicyset >/dev/null 2>&1
  test_except "test no permission set none policy" 1 $?
  echo ' 2)user testA set object policy'
  mc policy set public loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object set public policy" 0 $?
  mc policy set upload loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object set upload policy" 0 $?
  mc policy set download loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object set download policy" 0 $?
  mc policy set none loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object set none policy" 0 $?
  mc policy set public loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object no permission set public policy" 1 $?
  mc policy set upload loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object no permission set upload policy" 1 $?
  mc policy set download loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object no permission set download policy" 1 $?
  mc policy set none loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test object no permission set none policy" 1 $?
}

function test_policy_get() {
  echo 'test policy get'
  # creat bucket for test
  mc mb loc1/testapolicyget >/dev/null 2>&1
  mc mb loc1/testadminpolicyget >/dev/null 2>&1
  mc policy set public loc1/testadminpolicyget >/dev/null 2>&1
  mc cp test_cp.txt loc1/testapolicyget >/dev/null 2>&1 # text
  mc cp test_cp.txt loc1/testadminpolicyget >/dev/null 2>&1 # text
  echo ' 1)user testA get bucket policy'
  mc policy get  loc1/testapolicyget >/dev/null 2>&1
  test_except "test get policy" 0 $?
  mc policy get  loc2/testapolicyget >/dev/null 2>&1
  test_except "test get policy no permission" 1 $?
  echo ' 2)user testA get object policy'
  mc policy get  loc1/testapolicyget/test_cp.txt >/dev/null 2>&1
  test_except "test object get policy" 0 $?
  mc policy get  loc2/testapolicyget/test_cp.txt >/dev/null 2>&1
  test_except "test object no permission get  policy" 1 $?
}
init
test_policy_set
test_policy_get
close
