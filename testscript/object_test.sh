#!/bin/bash

. ./common.sh

function test_object_upload() {
    echo "test upload file"
    # creat bucket for test
    mc mb loc1/testaupload >/dev/null 2>&1
    mc mb loc1/testadminupload >/dev/null 2>&1
    mc anonymous set upload loc1/testadminupload >/dev/null 2>&1
    echo " 1)user 'testA' upload file"
    mc cp test_cp.txt loc1/testaupload >/dev/null 2>&1 # text
    test_except "test upload txt" 0 $?
    mc cp test_upload.jpg loc1/testaupload  >/dev/null 2>&1 # image
    test_except "test upload image" 0 $?
    mc cp test_upload.mp4 loc1/testaupload >/dev/null 2>&1 # video
    test_except "test upload video" 0 $?
    mc cp --disable-multipart test_upload_copy1.jpg loc1/testaupload >/dev/null 2>&1 # single upload
    test_except "test single upload" 0 $?
    mc cp test_upload_copy2.jpg loc1/testaupload >/dev/null 2>&1 # mul upload
    test_except "test multipart upload file" 0 $?
    mc cp testfile/test_upload.jpg loc1/testaupload >/dev/null 2>&1 # same name but content not same
    test_except "test upload file that has same name but content diff" 0 $?
    mc cp test_cp.txt loc2/testaupload >/dev/null 2>&1 # text
    test_except "test upload file to a bucket without permission" 1 $?
    echo " 2)user admin upload file"
    mc cp test_cp.txt loc/testadminupload >/dev/null 2>&1 # text
    test_except "user admin upload txt" 0 $?
    mc cp test_upload.jpg loc/testadminupload  >/dev/null 2>&1 # image
    test_except "user admin upload image" 0 $?
    mc cp test_upload.mp4 loc/testadminupload >/dev/null 2>&1 # video
    test_except "user admin upload video" 0 $?
    mc cp --disable-multipart test_upload_copy1.jpg loc/testadminupload >/dev/null 2>&1 # single upload
    test_except "user admin single upload" 0 $?
    mc cp test_upload_copy2.jpg loc/testadminupload >/dev/null 2>&1 # mul upload
    test_except "user admin multipart upload file" 0 $?
    mc cp testfile/test_upload.jpg loc/testadminupload >/dev/null 2>&1 # same name but content not same
    test_except "user admin upload file that has same name but content diff" 0 $?
}
function test_object_download() {
  echo "test object download"
  # creat bucket for test
  mc mb loc1/testadownload >/dev/null 2>&1
  mc mb loc1/testadmindownload >/dev/null 2>&1
  mc anonymous set download loc1/testadmindownload >/dev/null 2>&1
  mc cp test_cp.txt loc1/testadownload >/dev/null 2>&1 # text
  mc cp test_upload.jpg loc1/testadownload  >/dev/null 2>&1 # image
  mc cp test_upload.mp4 loc1/testadownload >/dev/null 2>&1 # video
  mc cp test_cp.txt loc1/testadmindownload >/dev/null 2>&1 # text
  mc cp test_upload.jpg loc1/testadmindownload  >/dev/null 2>&1 # image
  mc cp test_upload.mp4 loc1/testadmindownload >/dev/null 2>&1 # video
  echo " 1)user 'testA' download file"
  mc cat loc1/testadownload/test_cp.txt >/dev/null 2>&1 # text
  test_except "test download txt" 0 $?
  mc cat loc1/testadownload/test_upload.jpg  >/dev/null 2>&1 # image
  test_except "test download image" 0 $?
  mc cat loc1/testadownload/test_upload.mp4 >/dev/null 2>&1 # video
  test_except "test download video" 0 $?
  mc cat loc2/testadownload/test_cp.txt >/dev/null 2>&1 # no permission
  test_except "test download the file without permission" 1 $?
  mc cat loc1/testadownload/test_cp_no_exist.txt >/dev/null 2>&1 # no exist
  test_except "test download non-existing file" 1 $?
  echo " 2)user admin download file"
  mc cat loc/testadmindownload/test_cp.txt >/dev/null 2>&1 # text
  test_except "user admin download txt" 0 $?
  mc cat loc/testadmindownload/test_upload.jpg  >/dev/null 2>&1 # image
  test_except "user admin download image" 0 $?
  mc cat loc/testadmindownload/test_upload.mp4 >/dev/null 2>&1 # video
  test_except "user admin download video" 0 $?
  mc cat loc/testadownload/test_cp_no_exist.txt >/dev/null 2>&1 # no exist
  test_except "user admin download non-existing file" 1 $?
}
function test_object_delete() {
    echo "test object delete"
    # creat bucket for test
    mc mb loc1/testadel >/dev/null 2>&1
    mc mb loc1/testadmindel >/dev/null 2>&1
    mc anonymous set public loc1/testadmindel >/dev/null 2>&1
    mc cp test_cp.txt loc1/testadel >/dev/null 2>&1 # text
    mc cp test_cp.txt loc1/testadmindel >/dev/null 2>&1 # text
    echo " 1)user 'testA' delete file"
    mc cat loc1/testadel/test_cp.txt >/dev/null 2>&1 # exist
    test_except "test delete file" 0 $?
    mc cat loc1/testadel/test_cp_no_exist.txt >/dev/null 2>&1 # no exist
    test_except "test delete non-existing file" 1 $?
    mc cat test_cp.txt loc2/testadel/test_cp.txt >/dev/null 2>&1 # no permission
    test_except "test delete file without permission" 1 $?
    echo " 2)admin delete file"
    mc cat loc/testadmindel/test_cp.txt >/dev/null 2>&1 # text
    test_except "user admin delete file" 0 $?
    mc cat loc1/testadmindel/test_cp_no_exist.txt >/dev/null 2>&1 # no exist
    test_except "user admin delete  non-existing file" 1 $?
}
init
test_object_upload
test_object_download
test_object_delete
close
