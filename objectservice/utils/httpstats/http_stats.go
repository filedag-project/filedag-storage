package httpstats

import (
	"context"
	"github.com/filedag-project/filedag-storage/objectservice/consts"
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	logging "github.com/ipfs/go-log/v2"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	apiRecordTemplate = "api-record"
	storeDuration     = time.Second * 10
)

var log = logging.Logger("http-stats")

// HTTPStats holds statistics information about
// HTTP requests made by all clients
type HTTPStats struct {
	sync.RWMutex
	currentS3Requests  HTTPAPIStats
	currentIamRequests HTTPAPIStats
	totalIamRequests   HTTPAPIStats
	totalIamErrors     HTTPAPIStats
	totalIam4xxErrors  HTTPAPIStats
	totalIam5xxErrors  HTTPAPIStats
	totalIamCanceled   HTTPAPIStats
	totalS3Requests    HTTPAPIStats
	totalS3Errors      HTTPAPIStats
	totalS34xxErrors   HTTPAPIStats
	totalS35xxErrors   HTTPAPIStats
	totalS3Canceled    HTTPAPIStats
}
type APIStatsSys struct {
	Db         *uleveldb.ULevelDB
	HttpStats  *HTTPStats
	ObjectInfo *ObjectInfo
}
type ObjectInfo struct {
	PutObjCount map[string]uint64
	GetObjCount map[string]uint64
	PutObjBytes map[string]uint64
	GetObjBytes map[string]uint64
	sync.RWMutex
}

// NewHttpStatsSys - new an HttpStats  system
func NewHttpStatsSys(db *uleveldb.ULevelDB) *APIStatsSys {
	apiStatsSys := &APIStatsSys{Db: db, HttpStats: &HTTPStats{}, ObjectInfo: &ObjectInfo{}}
	apiStatsSys.load()
	return apiStatsSys
}
func (st *APIStatsSys) StoreApiLog(ctx context.Context) {
	tc := time.NewTicker(storeDuration)
	for {
		select {
		case <-ctx.Done():
			st.store()
		case <-tc.C:
			st.store()
		}

	}
}
func (st *APIStatsSys) RecordAPIHandler(api string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "admin/v1") || strings.Contains(r.URL.Path, "consoles/v1") {
			st.HttpStats.currentIamRequests.inc(api)
			defer st.HttpStats.currentS3Requests.dec(api)
		} else {
			st.HttpStats.currentS3Requests.inc(api)
			defer st.HttpStats.currentS3Requests.dec(api)
		}
		statsWriter := NewResponseRecorder(w)

		f.ServeHTTP(statsWriter, r)

		st.HttpStats.updateStats(api, r, statsWriter)
		st.updateObjInfo(api, r, statsWriter)
	}
}

// HTTPAPIStats holds statistics information about
// a given API in the requests.
type HTTPAPIStats struct {
	apiStats map[string]int
	sync.RWMutex
}

// inc increments the api stats counter.
func (stats *HTTPAPIStats) inc(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if stats.apiStats == nil {
		stats.apiStats = make(map[string]int)
	}
	stats.apiStats[api]++
}

// dec increments the api stats counter.
func (stats *HTTPAPIStats) dec(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if val, ok := stats.apiStats[api]; ok && val > 0 {
		stats.apiStats[api]--
	}
}

// inc increments the api stats counter.
func (obj *ObjectInfo) inc(put bool, filetype string, bytes uint64) {
	if obj == nil {
		return
	}
	obj.Lock()
	defer obj.Unlock()
	if put {
		if obj.PutObjCount == nil {
			obj.PutObjCount = make(map[string]uint64)
		}
		obj.PutObjCount["PutObjCount"+filetype]++
		if obj.PutObjBytes == nil {
			obj.PutObjBytes = make(map[string]uint64)
		}
		obj.PutObjBytes["PutObjBytes"+filetype] += bytes
	} else {
		if obj.GetObjCount == nil {
			obj.GetObjCount = make(map[string]uint64)
		}
		obj.GetObjCount["GetObjectHandler"+filetype]++
		if obj.GetObjBytes == nil {
			obj.GetObjBytes = make(map[string]uint64)
		}
		obj.GetObjBytes["GetObjBytes"+filetype] += bytes
	}

}

// Update statistics from http request and response data
func (st *HTTPStats) updateStats(api string, r *http.Request, w *ResponseRecorder) {
	if strings.Contains(r.URL.Path, "admin/v1") || strings.Contains(r.URL.Path, "consoles/v1") {
		st.totalIamRequests.inc(api)
		code := w.StatusCode
		switch {
		case code == 0:
		case code == 499:
			// 499 is a good error, shall be counted as canceled.
			st.totalIamCanceled.inc(api)
		case code >= http.StatusBadRequest:
			st.totalIamErrors.inc(api)
			if code >= http.StatusInternalServerError {
				st.totalIam5xxErrors.inc(api)
			} else {
				st.totalIam4xxErrors.inc(api)
			}
		}
	} else {
		st.totalS3Requests.inc(api)
		code := w.StatusCode
		switch {
		case code == 0:
		case code == 499:
			// 499 is a good error, shall be counted as canceled.
			st.totalS3Canceled.inc(api)
		case code >= http.StatusBadRequest:
			st.totalS3Errors.inc(api)
			if code >= http.StatusInternalServerError {
				st.totalS35xxErrors.inc(api)
			} else {
				st.totalS34xxErrors.inc(api)
			}
		}
	}
}

// Update statistics from http request and response data
func (st *APIStatsSys) updateObjInfo(api string, r *http.Request, w *ResponseRecorder) {
	if strings.Contains(r.URL.Path, "admin/v1") || strings.Contains(r.URL.Path, "consoles/v1") {
		return
	} else {
		fileType := r.Header.Get(consts.ContentType)
		if fileType == "" {
			fileType = "unknown"
		}
		st.ObjectInfo.inc(api == "PutObjectHandler", fileType, uint64(r.ContentLength))
	}
}

func (st *APIStatsSys) store() {
	st.HttpStats.Lock()
	defer st.HttpStats.Unlock()
	err := st.Db.Put(apiRecordTemplate+"totalIamRequests", st.HttpStats.totalIamRequests.apiStats)
	if err != nil {
		log.Errorf("store totalIamRequests err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalIamErrors", st.HttpStats.totalIamErrors.apiStats)
	if err != nil {
		log.Errorf("store totalIamErrors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalIam4xxErrors", st.HttpStats.totalIam4xxErrors.apiStats)
	if err != nil {
		log.Errorf("store totalIam4xxErrors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalIam5xxErrors", st.HttpStats.totalIam5xxErrors.apiStats)
	if err != nil {
		log.Errorf("store totalIam5xxErrors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalIamCanceled", st.HttpStats.totalIamCanceled.apiStats)
	if err != nil {
		log.Errorf("store totalIamCanceled err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalS3Requests", st.HttpStats.totalS3Requests.apiStats)
	if err != nil {
		log.Errorf("store totalS3Requests err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalS3Errors", st.HttpStats.totalS3Errors.apiStats)
	if err != nil {
		log.Errorf("store totalS3Errors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalS34xxErrors", st.HttpStats.totalS34xxErrors.apiStats)
	if err != nil {
		log.Errorf("store totalS34xxErrors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalS35xxErrors", st.HttpStats.totalS35xxErrors.apiStats)
	if err != nil {
		log.Errorf("store totalS35xxErrors err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"totalS3Canceled", st.HttpStats.totalS3Canceled.apiStats)
	if err != nil {
		log.Errorf("store totalS3Canceled err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"PutObjCount", &st.ObjectInfo.PutObjCount)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("store PutObjCount err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"GetObjCount", &st.ObjectInfo.GetObjCount)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("store GetObjCount err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"PutObjBytes", &st.ObjectInfo.PutObjBytes)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("store PutObjBytes err%v", err)
	}
	err = st.Db.Put(apiRecordTemplate+"GetObjBytes", &st.ObjectInfo.GetObjBytes)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("store GetObjBytes err%v", err)
	}
}
func (st *APIStatsSys) load() {
	st.HttpStats.RLock()
	defer st.HttpStats.RUnlock()
	err := st.Db.Get(apiRecordTemplate+"totalIamRequests", &st.HttpStats.totalIamRequests.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalIamRequests err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalIamErrors", &st.HttpStats.totalIamErrors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalIamErrors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalIam4xxErrors", &st.HttpStats.totalIam4xxErrors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalIam4xxErrors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalIam5xxErrors", &st.HttpStats.totalIam5xxErrors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalIam5xxErrors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalIamCanceled", &st.HttpStats.totalIamCanceled.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalIamCanceled err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalS3Requests", &st.HttpStats.totalS3Requests.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalS3Requests err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalS3Errors", &st.HttpStats.totalS3Errors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalS3Errors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalS34xxErrors", &st.HttpStats.totalS34xxErrors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalS34xxErrors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalS35xxErrors", &st.HttpStats.totalS35xxErrors.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalS35xxErrors err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"totalS3Canceled", &st.HttpStats.totalS3Canceled.apiStats)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load totalS3Canceled err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"PutObjCount", &st.ObjectInfo.PutObjCount)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load PutObjCount err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"GetObjCount", &st.ObjectInfo.GetObjCount)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load GetObjCount err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"PutObjBytes", &st.ObjectInfo.PutObjBytes)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load PutObjBytes err%v", err)
	}
	err = st.Db.Get(apiRecordTemplate+"GetObjBytes", &st.ObjectInfo.GetObjBytes)
	if err != nil && err != leveldb.ErrNotFound {
		log.Errorf("load GetObjBytes err%v", err)
	}
}

type ApiStats struct {
	ApiStats map[string]int
}
type ObjInfo struct {
}

// StatsResp holds statistics information about
// HTTP requests made by all clients
type StatsResp struct {
	CurrentS3Requests  ApiStats          `json:"current_s3_requests"`
	CurrentIamRequests ApiStats          `json:"current_iam_requests"`
	TotalIamRequests   ApiStats          `json:"total_iam_requests"`
	TotalIamErrors     ApiStats          `json:"total_iam_errors"`
	TotalIam4xxErrors  ApiStats          `json:"total_iam_4xx_errors"`
	TotalIam5xxErrors  ApiStats          `json:"total_iam_5xx_errors"`
	TotalIamCanceled   ApiStats          `json:"total_iam_canceled"`
	TotalS3Requests    ApiStats          `json:"total_s3_requests"`
	TotalS3Errors      ApiStats          `json:"total_s3_errors"`
	TotalS34xxErrors   ApiStats          `json:"total_s3_4xx_errors"`
	TotalS35xxErrors   ApiStats          `json:"total_s3_5xx_errors"`
	TotalS3Canceled    ApiStats          `json:"total_s3_canceled"`
	PutObjCount        map[string]uint64 `json:"put_obj_count"`
	GetObjCount        map[string]uint64 `json:"get_obj_count"`
	PutObjBytes        map[string]uint64 `json:"put_obj_bytes"`
	GetObjBytes        map[string]uint64 `json:"get_obj_bytes"`
}

func (st *APIStatsSys) GetCurrentStats(ctx context.Context) (StatsResp, error) {
	st.HttpStats.RLock()
	defer st.HttpStats.RUnlock()
	var statsResp = StatsResp{
		CurrentS3Requests:  ApiStats{ApiStats: st.HttpStats.currentS3Requests.apiStats},
		CurrentIamRequests: ApiStats{ApiStats: st.HttpStats.currentIamRequests.apiStats},
		TotalIamRequests:   ApiStats{ApiStats: st.HttpStats.totalIamRequests.apiStats},
		TotalIamErrors:     ApiStats{ApiStats: st.HttpStats.totalIamErrors.apiStats},
		TotalIam4xxErrors:  ApiStats{ApiStats: st.HttpStats.totalIam4xxErrors.apiStats},
		TotalIam5xxErrors:  ApiStats{ApiStats: st.HttpStats.totalIam5xxErrors.apiStats},
		TotalIamCanceled:   ApiStats{ApiStats: st.HttpStats.totalIamCanceled.apiStats},
		TotalS3Requests:    ApiStats{ApiStats: st.HttpStats.totalS3Requests.apiStats},
		TotalS3Errors:      ApiStats{ApiStats: st.HttpStats.totalS3Errors.apiStats},
		TotalS34xxErrors:   ApiStats{ApiStats: st.HttpStats.totalS34xxErrors.apiStats},
		TotalS35xxErrors:   ApiStats{ApiStats: st.HttpStats.totalS35xxErrors.apiStats},
		TotalS3Canceled:    ApiStats{ApiStats: st.HttpStats.totalS3Canceled.apiStats},
		PutObjCount:        st.ObjectInfo.PutObjCount,
		GetObjCount:        st.ObjectInfo.GetObjCount,
		PutObjBytes:        st.ObjectInfo.PutObjBytes,
		GetObjBytes:        st.ObjectInfo.GetObjBytes,
	}
	return statsResp, nil
}
