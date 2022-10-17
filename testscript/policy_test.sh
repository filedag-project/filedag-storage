#!/bin/bash

. ./common.sh

function test_policy_set() {
  echo 'test setting policy '
  # creat bucket for test
  mc mb loc1/testapolicyset >/dev/null 2>&1
  mc mb loc1/testadminpolicyset >/dev/null 2>&1
  mc policy set public loc1/testadminpolicyset >/dev/null 2>&1
  mc cp test_cp.txt loc1/testapolicyset >/dev/null 2>&1 # text
  mc cp test_cp.txt loc1/testadminpolicyset >/dev/null 2>&1 # text
  echo "' 1)user 'testA' set bucket policy'"
  mc policy set public loc1/testapolicyset >/dev/null 2>&1
  test_except "test setting public policy to bucket" 0 $?
  mc policy set upload loc1/testapolicyset >/dev/null 2>&1
  test_except "test setting upload policy to bucket" 0 $?
  mc policy set download loc1/testapolicyset >/dev/null 2>&1
  test_except "test setting download policy to bucket" 0 $?
  mc policy set none loc1/testapolicyset >/dev/null 2>&1
  test_except "test setting none policy to bucket" 0 $?
  mc policy set public loc2/testapolicyset >/dev/null 2>&1
  test_except "test setting public policy to bucket without permission" 1 $?
  mc policy set upload loc2/testapolicyset >/dev/null 2>&1
  test_except "test setting upload policy to bucket without permission" 1 $?
  mc policy set download loc2/testapolicyset >/dev/null 2>&1
  test_except "test setting download policy to bucket without permission" 1 $?
  mc policy set none loc2/testapolicyset >/dev/null 2>&1
  test_except "test setting none policy to bucket without permission" 1 $?
  echo ' 2)user testA set object policy'
  mc policy set public loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'public' policy to object" 0 $?
  mc policy set upload loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'upload' policy to object" 0 $?
  mc policy set download loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'download' policy to object" 0 $?
  mc policy set none loc1/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'none' policy to object" 0 $?
  mc policy set public loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'public' policy to object without permission" 1 $?
  mc policy set upload loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'upload' policy to object without permission" 1 $?
  mc policy set download loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'download' policy to object without permission" 1 $?
  mc policy set none loc2/testapolicyset/test_cp.txt >/dev/null 2>&1
  test_except "test setting 'none' policy to object without permission" 1 $?
}

function test_policy_get() {
  echo 'test getting policy'
  # creat bucket for test
  mc mb loc1/testapolicyget >/dev/null 2>&1
  mc mb loc1/testadminpolicyget >/dev/null 2>&1
  mc policy set public loc1/testadminpolicyget >/dev/null 2>&1
  mc cp test_cp.txt loc1/testapolicyget >/dev/null 2>&1 # text
  mc cp test_cp.txt loc1/testadminpolicyget >/dev/null 2>&1 # text
  echo "' 1)user 'testA' getting bucket policy'"
  mc policy get  loc1/testapolicyget >/dev/null 2>&1
  test_except "test getting policy" 0 $?
  mc policy get  loc2/testapolicyget >/dev/null 2>&1
  test_except "test getting a bucket policy without permission" 1 $?
  echo ' 2)user testA getting object policy'
  mc policy get  loc1/testapolicyget/test_cp.txt >/dev/null 2>&1
  test_except "test getting a object policy" 0 $?
  mc policy get  loc2/testapolicyget/test_cp.txt >/dev/null 2>&1
  test_except "test getting a object policy without permission" 1 $?
}
init
test_policy_set
test_policy_get
close
