package iamapi

//func TestIamApiServer_ListGroups(t *testing.T) {
//	u := "http://127.0.0.1:9985/admin/v1/list-groups?pathPrefix"
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
