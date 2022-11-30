package httpstatss

import (
	"net/http"
	"strings"
	"sync"
)

func (st *HTTPStats) RecordAPIHandler(api string, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		st.currentS3Requests.Inc(api)
		defer st.currentS3Requests.Dec(api)

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

// Inc increments the api stats counter.
func (stats *HTTPAPIStats) Inc(api string) {
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

// Dec increments the api stats counter.
func (stats *HTTPAPIStats) Dec(api string) {
	if stats == nil {
		return
	}
	stats.Lock()
	defer stats.Unlock()
	if val, ok := stats.apiStats[api]; ok && val > 0 {
		stats.apiStats[api]--
	}
}

// HTTPStats holds statistics information about
// HTTP requests made by all clients
type HTTPStats struct {
	s3RequestsInQueue       int32 // ref: https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	_                       int32 // For 64 bits alignment
	s3RequestsIncoming      uint64
	rejectedRequestsAuth    uint64
	rejectedRequestsTime    uint64
	rejectedRequestsHeader  uint64
	rejectedRequestsInvalid uint64
	currentS3Requests       HTTPAPIStats
	totalIamRequests        HTTPAPIStats
	totalIamErrors          HTTPAPIStats
	totalIam4xxErrors       HTTPAPIStats
	totalIam5xxErrors       HTTPAPIStats
	totalIamCanceled        HTTPAPIStats
	totalS3Requests         HTTPAPIStats
	totalS3Errors           HTTPAPIStats
	totalS34xxErrors        HTTPAPIStats
	totalS35xxErrors        HTTPAPIStats
	totalS3Canceled         HTTPAPIStats
}

// Update statistics from http request and response data
func (st *HTTPStats) updateStats(api string, r *http.Request, w *ResponseRecorder) {
	// Ignore non S3 requests
	if strings.Contains(r.URL.Path, "admin/v1") {
		st.totalIamRequests.Inc(api)
		code := w.StatusCode
		switch {
		case code == 0:
		case code == 499:
			// 499 is a good error, shall be counted as canceled.
			st.totalIamCanceled.Inc(api)
		case code >= http.StatusBadRequest:
			st.totalIamErrors.Inc(api)
			if code >= http.StatusInternalServerError {
				st.totalIam5xxErrors.Inc(api)
			} else {
				st.totalIam4xxErrors.Inc(api)
			}
		}
	} else {
		st.totalS3Requests.Inc(api)
		code := w.StatusCode
		switch {
		case code == 0:
		case code == 499:
			// 499 is a good error, shall be counted as canceled.
			st.totalS3Canceled.Inc(api)
		case code >= http.StatusBadRequest:
			st.totalS3Errors.Inc(api)
			if code >= http.StatusInternalServerError {
				st.totalS35xxErrors.Inc(api)
			} else {
				st.totalS34xxErrors.Inc(api)
			}
		}
	}
}

// NewHTTPStats Prepare new HTTPStats structure
func NewHTTPStats() *HTTPStats {
	return &HTTPStats{}
}
