#!/bin/bash

. ./common.sh
# diff user creat bucket
function test_diff_user_creat_bucket() {
   echo 'diff user creat bucket:'
   mc mb loc1/test1 >/dev/null 2>&1
   test_except "user 'testA' creat  bucket" 0 $?
   mc mb loc2/test2 >/dev/null 2>&1
   test_except "user 'testB' creat  bucket" 0 $?
   mc mb loc/test3 >/dev/null 2>&1
   test_except "user 'admin' creat  bucket" 0 $?
}

# testA user creat diff buckets
function testA_user_creat_diff_bucket() {
  echo "user 'testA' creat diff buckets:"
  mc mb loc1/loc1creatbucket  >/dev/null 2>&1 # success
  test_except "user 'testA' creat normal bucket" 0 $?
  mc mb loc1/loc1creatbucketA  >/dev/null 2>&1  # fail
  test_except "user 'testA' creat illegal bucket loc1creatbucketA" 1 $?
  mc mb loc1/te >/dev/null 2>&1 # fail
  test_except "user 'testA' creat illegal bucket te " 1 $?
  mc mb loc1/loc1test- >/dev/null 2>&1 # fail
  test_except "user 'testA' creat illegal bucket loc1test-" 1 $?
  mc mb loc1/loc1test_ >/dev/null 2>&1 # fail
  test_except "user 'testA' creat illegal bucket loc1test_" 1 $?
  mc mb loc1/loc1test_1 >/dev/null 2>&1 # fail
  test_except "user 'testA' creat illegal bucket loc1test_1" 1 $?
}

# test user delete bucket
function test_user_del_bucket() {
   echo 'user delete bucket:'
   echo " 1)user 'testA'  delete bucket"
   # '1) creat bucket for test'
   mc mb loc1/loc1bucketexistempty >/dev/null 2>&1 # bucket exist but empty
   mc mb loc1/loc1bucketexistbutnotempty >/dev/null 2>&1 # bucket exist not empty
   mc cp test_cp.txt loc1/loc1bucketexistbutnotempty >/dev/null 2>&1
   mc mb loc1/loc1bucketexistbutnopermission >/dev/null 2>&1 # no permission
   # 2) test
   mc rb loc1/loc1bucketexistempty >/dev/null 2>&1 # success
   test_except "user 'testA' remove exist&empty bucket" 0 $?
   mc rb loc1/loc1bucketnotexist >/dev/null 2>&1 # fail
   test_except "user 'testA' remove non-existing bucket" 1 $?
   mc rb loc1/loc1bucketexistbutnotempty >/dev/null 2>&1 # fail
   test_except "user 'testA' remove bucket that is not empty" 1 $?
   mc rb loc2/loc1bucketexistbutnopermission >/dev/null 2>&1 # fail
   test_except "user 'testA' remove bucket that has no permissions" 1 $?

   # admin user delete bucket
   echo ' 2)the admin user delete bucket'
   # '1) creat bucket for test'
   mc mb loc/locbucketexistbutempty >/dev/null 2>&1
   mc mb loc/locbucketexistbutnotempty >/dev/null 2>&1
   mc cp test_cp.txt loc/locbucketexistbutnotempty >/dev/null 2>&1
   # '2) test'
   mc rb loc/locbucketexistbutempty >/dev/null 2>&1 # success
   test_except "user admin remove  exist&empty bucket" 0 $?
   mc rb loc/locbucketnotexist >/dev/null 2>&1 # fail
   test_except "user admin remove non-existing bucket" 1 $?
   mc rb loc/locbucketexistbutnotempty >/dev/null 2>&1 # fail
   test_except "user admin remove bucket that is not empty" 1 $?
}

# list bucket
function test_list_bucket() {
  echo 'list bucket:'
   echo " 1)user 'testA' list bucket"
   # "1) make bucket for test"
   mc mb loc1/locbucketexisthaveper >/dev/null 2>&1
   mc cp test_cp.txt loc1/locbucketexisthaveper >/dev/null 2>&1 # exist and have permission
   mc mb loc1/locbucketexistnothaveper >/dev/null 2>&1 # exist but have not permission
   # "2) test"
   mc ls loc1/locbucketexisthaveper >/dev/null 2>&1 #success
   test_except "user 'testA' ls an existing and authorized bucket" 0 $?
   mc ls loc1/locbucketnotexist >/dev/null 2>&1 # fail
   test_except "user 'testA' ls non-existing bucket" 1 $?
   mc ls loc2/locbucketexistnothaveper >/dev/null 2>&1 # fail
   test_except "user 'testB' ls bucket without permission" 1 $?
   echo " 2)the admin user list bucket"
   # "1) make bucket for test"
   mc mb loc/locbucketexist >/dev/null 2>&1
   mc cp test_cp.txt loc/locbucketexist >/dev/null 2>&1
   # "2) test"
   mc ls loc/locbucketexist >/dev/null 2>&1 #success
   test_except "user admin ls an existing bucket" 0 $?
   mc ls loc/locbucketnotexist >/dev/null 2>&1 # fail
   test_except "user admin ls a non-existing bucket" 1 $?
}

init
test_diff_user_creat_bucket
testA_user_creat_diff_bucket
test_user_del_bucket
test_list_bucket
close
