package iamapi

//const (
//	DefaultTestAccessKey = "test"
//	DefaultTestSecretKey = "test"
//)
//
//func TestIamApiServer_GetUserList(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/list-user"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestIamApiServer_AddUser(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/add-user?accessKey=test1&secretKey=test12345"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestIamApiServer_RemoveUser(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/remove-user?accessKey=test1"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u, 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//
//func TestIamApiServer_ChangePassword(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("newPassword", "test2222")
//	u := "http://127.0.0.1:9985/admin/v1/change-password?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", "test1", "test12345", t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//
//func TestIamApiServer_SetStatus(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("userName", "test1")
//	urlValues.Set("status", string(iam.AccountDisabled))
//	//urlValues.Set("status", string(iam.AccountEnabled))
//	u := "http://127.0.0.1:9985/admin/v1/update-accessKey_status?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", "test1", "test2222", t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//
//func TestIamApiServer_PutUserPolicy(t *testing.T) {
//	urlValues := make(url.Values)
//	policy := `{"Version":"2008-10-17","Statement":[{"Effect":"Allow","Sid":"1","Principal":{"AWS":["111122223333","444455556666"]},"Action":["s3:*"],"Resource":"arn:aws:s3:::test22/*"}]}`
//	urlValues.Set("policyDocument", policy)
//	urlValues.Set("userName", "test1")
//	urlValues.Set("policyName", "read2")
//	u := "http://127.0.0.1:9985/admin/v1/put-user-policy?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestIamApiServer_GetUserPolicy(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("userName", "test1")
//	urlValues.Set("policyName", "read")
//	u := "http://127.0.0.1:9985/admin/v1/get-user-policy?"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestIamApiServer_RemoveUserPolicy(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("userName", "test1")
//	urlValues.Set("policyName", "read")
//	u := "http://127.0.0.1:9985/admin/v1/remove-user-policy?"
//	req := testsign.MustNewSignedV4Request(http.MethodPost, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//func TestIamApiServer_ListUserPolicy(t *testing.T) {
//	urlValues := make(url.Values)
//	urlValues.Set("userName", "test1")
//	u := "http://127.0.0.1:9985/admin/v1/list-user-policy?"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+urlValues.Encode(), 0, nil, "s3", DefaultTestAccessKey, DefaultTestSecretKey, t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
//
//func TestIamApiServer_GetUserInfo(t *testing.T) {
//	urlValues := make(url.Values)
//	user := "test1"
//	urlValues.Set("userName", user)
//	u := "http://127.0.0.1:9985/admin/v1/user-info"
//	req := testsign.MustNewSignedV4Request(http.MethodGet, u+"?"+urlValues.Encode(), 0, nil, "s3", "test1", "test2222", t)
//
//	//req.Header.Set("Content-Type", "text/plain")
//	client := &http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//	body, err := ioutil.ReadAll(res.Body)
//
//	fmt.Println(res)
//	fmt.Println(string(body))
//}
