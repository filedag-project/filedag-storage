package restapi

import (
	"fmt"
	models "github.com/filedag-project/filedag-storage/cmd/console/model"
	"testing"
)

func Test_getListUsersResponse(t *testing.T) {
	session := &models.Principal{
		STSAccessKeyID:"W2W2JWQUI52SVMMJK5MV",
		STSSecretAccessKey: "vSZmxbVcx+lP3iUZgqqXT0PZTYH2HAEsuQ+9hLJT",
		STSSessionToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NLZXkiOiJXMlcySldRVUk1MlNWTU1KSzVNViIsImV4cCI6MTY0NjI5Mzk4OSwicGFyZW50IjoibWluaW9hZG1pbiJ9.kmC564DCOOiDUpl4FiAWDx0839tTtuuZmvxoN_tSxpnqOBt_W8zoZgASi-ag9jD29kUThnVjR4I92qIs-TTL9g",
		AccountAccessKey: "minioadmin",
		Hm: false,
	}
	got, got1 := getListUsersResponse(session)
	if got1 != nil {
		fmt.Println(got1)
	}
	fmt.Println(got)
}
