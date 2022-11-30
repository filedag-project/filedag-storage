package httpstats

import (
	"github.com/filedag-project/filedag-storage/objectservice/uleveldb"
	logging "github.com/ipfs/go-log/v2"
	"net/http"
	"strings"
	"sync"
)

const apiRecordTemplate = "api-record"

var log = logging.Logger("http-stats")

// HTTPStats holds statistics information about
// HTTP requests made by all clients
type HTTPStats struct {
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
	Db        *uleveldb.ULevelDB
	HttpStats *HTTPStats
}

// NewHttpStatsSys - new an HttpStats  system
func NewHttpStatsSys(db *uleveldb.ULevelDB) *APIStatsSys {
	return &APIStatsSys{Db: db, HttpStats: NewHTTPStats(db)}
}
func (st *APIStatsSys) StoreApiLog() {
	err := st.Db.Put(apiRecordTemplate, st.HttpStats)
	if err != nil {
		log.Errorf("store api info err")
		return
	}
}
func (st *HTTPStats) RecordAPIHandler(api string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st.currentS3Requests.inc(api)
		defer st.currentS3Requests.dec(api)

		statsWriter := NewResponseRecorder(w)

		f.ServeHTTP(statsWriter, r)

		st.updateStats(api, r, statsWriter)
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

// Update statistics from http request and response data
func (st *HTTPStats) updateStats(api string, r *http.Request, w *ResponseRecorder) {
	// Ignore non S3 requests
	if strings.Contains(r.URL.Path, "admin/v1") {
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

// NewHTTPStats Prepare new HTTPStats structure
func NewHTTPStats(db *uleveldb.ULevelDB) *HTTPStats {
	var h HTTPStats
	db.Get(apiRecordTemplate, &h)
	return &h
}
