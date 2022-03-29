package restapi

import (
	"fmt"
	"testing"
)

func TestNewConsoleCredentials(t *testing.T) {
	got, err := NewConsoleCredentials("test1", "testsecretKey", "us-east-1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(got)
	tokens, err := got.Get()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(tokens)
	//type args struct {
	//	accessKey string
	//	secretKey string
	//	location  string
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	want    *credentials.Credentials
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		got, err := NewConsoleCredentials("test", "test", tt.args.location)
	//		if (err != nil) != tt.wantErr {
	//			t.Errorf("NewConsoleCredentials() error = %v, wantErr %v", err, tt.wantErr)
	//			return
	//		}
	//		if !reflect.DeepEqual(got, tt.want) {
	//			t.Errorf("NewConsoleCredentials() got = %v, want %v", got, tt.want)
	//		}
	//	})
	//}
}
