package iamapi

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
)

type CommonResponse struct {
	ResponseMetadata struct {
		RequestId string `xml:"RequestId"`
	} `xml:"ResponseMetadata"`
}

type ListUsersResponse struct {
	CommonResponse
	XMLName         xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListUsersResponse"`
	ListUsersResult struct {
		Users       []*iam.User `xml:"Users>member"`
		IsTruncated bool        `xml:"IsTruncated"`
	} `xml:"ListUsersResult"`
}

type ListAccessKeysResponse struct {
	CommonResponse
	XMLName              xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListAccessKeysResponse"`
	ListAccessKeysResult struct {
		AccessKeyMetadata []*iam.AccessKeyMetadata `xml:"AccessKeyMetadata>member"`
		IsTruncated       bool                     `xml:"IsTruncated"`
	} `xml:"ListAccessKeysResult"`
}

type DeleteAccessKeyResponse struct {
	CommonResponse
	XMLName xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteAccessKeyResponse"`
}

type CreatePolicyResponse struct {
	CommonResponse
	XMLName            xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreatePolicyResponse"`
	CreatePolicyResult struct {
		Policy iam.Policy `xml:"Policy"`
	} `xml:"CreatePolicyResult"`
}

type CreateUserResponse struct {
	CommonResponse
	XMLName          xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateUserResponse"`
	CreateUserResult struct {
		User iam.User `xml:"User"`
	} `xml:"CreateUserResult"`
}

type DeleteUserResponse struct {
	CommonResponse
	XMLName xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ DeleteUserResponse"`
}

type PutUserPolicyResponse struct {
	CommonResponse
	XMLName xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ PutUserPolicyResponse"`
}

type GetUserPolicyResponse struct {
	CommonResponse
	XMLName             xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetUserPolicyResponse"`
	GetUserPolicyResult struct {
		UserName       string `xml:"UserName"`
		PolicyName     string `xml:"PolicyName"`
		PolicyDocument string `xml:"PolicyDocument"`
	} `xml:"GetUserPolicyResult"`
}

// ListUserPoliciesResponse
//<ListUserPoliciesResponse xmlns="https://iam.amazonaws.com/doc/2010-05-08/">
// <ListUserPoliciesResult>
//    <PolicyNames>
//       <member>AllAccessPolicy</member>
//       <member>KeyPolicy</member>
//    </PolicyNames>
//    <IsTruncated>false</IsTruncated>
// </ListUserPoliciesResult>
// <ResponseMetadata>
//    <RequestId>7a62c49f-347e-4fc4-9331-6e8eEXAMPLE</RequestId>
// </ResponseMetadata>
//</ListUserPoliciesResponse>
type ListUserPoliciesResponse struct {
	CommonResponse
	XMLName                xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ListUserPoliciesResponse"`
	ListUserPoliciesResult struct {
		PolicyNames struct {
			Member []string `xml:"Member"`
		} `xml:"PolicyNames"`
	}
}

// CreateGroupResponse CreateGroup Response
type CreateGroupResponse struct {
	CommonResponse
	XMLName           xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ CreateGroupResponse"`
	CreateGroupResult struct {
		G Group `xml:"Group"`
	} `xml:"CreateGroupResult"`
}

// GetGroupResponse GetGroup Response
type GetGroupResponse struct {
	CommonResponse
	XMLName     xml.Name    `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetGroupResponse"`
	GroupResult GroupResult `xml:"GroupResult"`
}
type GroupResult struct {
	G Group `xml:"Group"`
}
type Group struct {
	Path      string `xml:"Path"`
	GroupName string `xml:"GroupName"`
	GroupId   string `xml:"GroupId"`
	Arn       string `xml:"Arn"`
}

// ListGroupsResponse listGroup Response
type ListGroupsResponse struct {
	CommonResponse
	XMLName     xml.Name         `xml:"https://iam.amazonaws.com/doc/2010-05-08/ GetGroupResponse"`
	GroupResult ListGroupsResult `xml:"ListGroupsResult"`
}
type GroupMember struct {
	GM Group `xml:"Member"`
}
type ListGroupsResult struct {
	Groups []GroupMember `xml:"Groups"`
}
type ErrorResponse struct {
	CommonResponse
	XMLName xml.Name `xml:"https://iam.amazonaws.com/doc/2010-05-08/ ErrorResponse"`
	Error   struct {
		iam.ErrorDetails
		Type string `xml:"Type"`
	} `xml:"Error"`
}

func (r *CommonResponse) SetRequestId() {
	r.ResponseMetadata.RequestId = fmt.Sprintf("%d", time.Now().UnixNano())
}
