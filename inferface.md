## Control Pannel

### 一、存储池概况
- 展示指标：桶（bucket），对象（object），使用空间大小（usage），服务（servers（服务挂载目录（drives））），widgets
- 接口：/api/v1/admin/info   
- 类型：Get
- 函数：AdminAPIAdminInfoHandler() -- getUsageWidgetsForDeployment()
- 返回值：AdminInfoResponse{}
  ```
    type AdminInfoResponse struct {
	    // buckets
	    Buckets int64 `json:"buckets,omitempty"`
	    // objects
	    Objects int64 `json:"objects,omitempty"`
	    // prometheus not ready
	    PrometheusNotReady bool `json:"prometheusNotReady,omitempty"`
	    // servers
	    Servers []*ServerProperties `json:"servers"`
	    // usage
	    Usage int64 `json:"usage,omitempty"`
	    // widgets
	    Widgets []*Widget `json:"widgets"`
    }
  ```
### 二、桶操作界面
#### 1、桶列表查询
- 展示指标：创建时间，桶名称，桶包含对象数量，用户对桶的权限（read，write），桶内对象使用空间（size）
- 接口：/api/v1/buckets
- 类型：Get
- 函数：UserAPIListBucketsHandler() -- getListBucketsResponse() -- getAccountBuckets()
- 返回值：Bucket{}
``` 
    type Bucket struct {
        // access
        Access *BucketAccess `json:"access,omitempty"`
        // creation date
        CreationDate string `json:"creation_date,omitempty"`
        // definition
        Definition string `json:"definition,omitempty"`
        // details
        Details *BucketDetails `json:"details,omitempty"`
        // name
        // Required: true
        // Min Length: 3
        Name *string `json:"name"`
        // objects
        Objects int64 `json:"objects,omitempty"`
        // rw access
        RwAccess *BucketRwAccess `json:"rw_access,omitempty"`
        // size
        Size int64 `json:"size,omitempty"`
    }
```
#### 2、查询桶详情
- 展示指标：创建时间，桶名称，桶包含对象数量，用户对桶的权限（read，write），桶内对象使用空间（size）
- 接口：/api/v1/buckets/{bucket-name}
- 类型：Get
- 函数：getBucketInfoResponse() -- getBucketInfo()
- 返回值：Bucket{}
#### 3、查询桶内对象列表
- 展示指标：创建时间，桶名称，桶包含对象数量，用户对桶的权限（read，write），桶内对象使用空间（size）
- 接口：/api/v1/buckets/{bucket_name}/objects
- 类型：Get
- 函数：UserAPIListObjectsHandler() -- getListObjectsResponse() -- listBucketObjects()
- 返回值：BucketObject{}
```
type BucketObject struct {
	// content type
	ContentType string `json:"content_type,omitempty"`
	// expiration
	Expiration string `json:"expiration,omitempty"`
	// expiration rule id
	ExpirationRuleID string `json:"expiration_rule_id,omitempty"`
	// is delete marker
	IsDeleteMarker bool `json:"is_delete_marker,omitempty"`
	// is latest
	IsLatest bool `json:"is_latest,omitempty"`
	// last modified
	LastModified string `json:"last_modified,omitempty"`
	// legal hold status
	LegalHoldStatus string `json:"legal_hold_status,omitempty"`
	// metadata
	Metadata map[string]string `json:"metadata,omitempty"`
	// name
	Name string `json:"name,omitempty"`
	// retention mode
	RetentionMode string `json:"retention_mode,omitempty"`
	// retention until date
	RetentionUntilDate string `json:"retention_until_date,omitempty"`
	// size
	Size int64 `json:"size,omitempty"`
	// tags
	Tags map[string]string `json:"tags,omitempty"`
	// user metadata
	UserMetadata map[string]string `json:"user_metadata,omitempty"`
	// user tags
	UserTags map[string]string `json:"user_tags,omitempty"`
	// version id
	VersionID string `json:"version_id,omitempty"`
}
```
#### 4、删除桶
- 展示指标：
- 接口：/api/v1/buckets/{bucket-name}
- 类型：Delete
- 函数：UserAPIDeleteBucketHandler() -- getDeleteBucketResponse() -- removeBucket()
- 返回值：
### 三、用户操作
#### 1、用户列表查询
- 展示指标：创建时间，桶名称，桶包含对象数量，用户对桶的权限（read，write），桶内对象使用空间（size）
- 接口：/api/v1/users
- 类型：Get
- 函数：AdminAPIListUsersHandler() -- getListUsersResponse() -- listUsers()
- 返回值：User{}
``` 
type User struct {
	// access key
	AccessKey string `json:"accessKey,omitempty"`
	// has policy
	HasPolicy bool `json:"hasPolicy,omitempty"`
	// member of
	MemberOf []string `json:"memberOf"`
	// policy
	Policy []string `json:"policy"`
	// status
	Status string `json:"status,omitempty"`
}
```