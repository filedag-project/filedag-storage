//// This file is part of MinIO Console Server
//// Copyright (c) 2021 MinIO, Inc.
////
//// This program is free software: you can redistribute it and/or modify
//// it under the terms of the GNU Affero General Public License as published by
//// the Free Software Foundation, either version 3 of the License, or
//// (at your option) any later version.
////
//// This program is distributed in the hope that it will be useful,
//// but WITHOUT ANY WARRANTY; without even the implied warranty of
//// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//// GNU Affero General Public License for more details.
////
//// You should have received a copy of the GNU Affero General Public License
//// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
package restapi

import (
	"context"
	"encoding/base64"
	"github.com/filedag-project/filedag-storage/http/console/madmin/object"
	"github.com/filedag-project/filedag-storage/http/console/models"
	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/user_api"
	"io"
	"path"
	"strconv"
	"time"
)

//
//import (
//	"archive/zip"
//	"bytes"
//	"context"
//	"encoding/base64"
//	"fmt"
//	"github.com/filedag-project/filedag-storage/http/console/madmin/tags"
//	"github.com/filedag-project/filedag-storage/http/console/models"
//	"github.com/filedag-project/filedag-storage/http/console/restapi/operations/user_api"
//	minio "github.com/filedag-project/filedag-storage/http/console/s3api"
//	"io"
//	"net/http"
//	"net/url"
//	"path"
//	"path/filepath"
//	"regexp"
//	"strconv"
//	"strings"
//	"time"
//
//	"errors"
//
//	"github.com/go-openapi/runtime"
//	"github.com/go-openapi/runtime/middleware"
//)
//
//// enum types
//const (
//	objectStorage = iota // MinIO and S3 compatible cloud storage
//	fileSystem           // POSIX compatible file systems
//)
//
//func registerObjectsHandlers(api *operations.ConsoleAPI) {
//	// list objects
//	api.UserAPIListObjectsHandler = user_api.ListObjectsHandlerFunc(func(params user_api.ListObjectsParams, session *models.Principal) middleware.Responder {
//		resp, err := getListObjectsResponse(session, params)
//		if err != nil {
//			return user_api.NewListObjectsDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewListObjectsOK().WithPayload(resp)
//	})
//	// delete object
//	api.UserAPIDeleteObjectHandler = user_api.DeleteObjectHandlerFunc(func(params user_api.DeleteObjectParams, session *models.Principal) middleware.Responder {
//		if err := getDeleteObjectResponse(session, params); err != nil {
//			return user_api.NewDeleteObjectDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewDeleteObjectOK()
//	})
//	// delete multiple objects
//	api.UserAPIDeleteMultipleObjectsHandler = user_api.DeleteMultipleObjectsHandlerFunc(func(params user_api.DeleteMultipleObjectsParams, session *models.Principal) middleware.Responder {
//		if err := getDeleteMultiplePathsResponse(session, params); err != nil {
//			return user_api.NewDeleteMultipleObjectsDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewDeleteMultipleObjectsOK()
//	})
//	// download object
//	api.UserAPIDownloadObjectHandler = user_api.DownloadObjectHandlerFunc(func(params user_api.DownloadObjectParams, session *models.Principal) middleware.Responder {
//		isFolder := false
//
//		var prefix string
//		if params.Prefix != "" {
//			encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//			decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//			if err != nil {
//				return user_api.NewDownloadObjectDefault(int(400)).WithPayload(prepareError(err))
//			}
//			prefix = string(decodedPrefix)
//		}
//
//		folders := strings.Split(prefix, "/")
//		if folders[len(folders)-1] == "" {
//			isFolder = true
//		}
//		var resp middleware.Responder
//		var err *models.Error
//
//		if isFolder {
//			resp, err = getDownloadFolderResponse(session, params)
//		} else {
//			resp, err = getDownloadObjectResponse(session, params)
//		}
//
//		if err != nil {
//			return user_api.NewDownloadObjectDefault(int(err.Code)).WithPayload(err)
//		}
//		return resp
//	})
//	// upload object
//	api.UserAPIPostBucketsBucketNameObjectsUploadHandler = user_api.PostBucketsBucketNameObjectsUploadHandlerFunc(func(params user_api.PostBucketsBucketNameObjectsUploadParams, session *models.Principal) middleware.Responder {
//		if err := getUploadObjectResponse(session, params); err != nil {
//			if strings.Contains(*err.DetailedMessage, "413") {
//				return user_api.NewPostBucketsBucketNameObjectsUploadDefault(413).WithPayload(err)
//			}
//			return user_api.NewPostBucketsBucketNameObjectsUploadDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPostBucketsBucketNameObjectsUploadOK()
//	})
//	// get share object url
//	api.UserAPIShareObjectHandler = user_api.ShareObjectHandlerFunc(func(params user_api.ShareObjectParams, session *models.Principal) middleware.Responder {
//		resp, err := getShareObjectResponse(session, params)
//		if err != nil {
//			return user_api.NewShareObjectDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewShareObjectOK().WithPayload(*resp)
//	})
//	// set object legalhold status
//	api.UserAPIPutObjectLegalHoldHandler = user_api.PutObjectLegalHoldHandlerFunc(func(params user_api.PutObjectLegalHoldParams, session *models.Principal) middleware.Responder {
//		if err := getSetObjectLegalHoldResponse(session, params); err != nil {
//			return user_api.NewPutObjectLegalHoldDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPutObjectLegalHoldOK()
//	})
//	// set object retention
//	api.UserAPIPutObjectRetentionHandler = user_api.PutObjectRetentionHandlerFunc(func(params user_api.PutObjectRetentionParams, session *models.Principal) middleware.Responder {
//		if err := getSetObjectRetentionResponse(session, params); err != nil {
//			return user_api.NewPutObjectRetentionDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPutObjectRetentionOK()
//	})
//	// delete object retention
//	api.UserAPIDeleteObjectRetentionHandler = user_api.DeleteObjectRetentionHandlerFunc(func(params user_api.DeleteObjectRetentionParams, session *models.Principal) middleware.Responder {
//		if err := deleteObjectRetentionResponse(session, params); err != nil {
//			return user_api.NewDeleteObjectRetentionDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewDeleteObjectRetentionOK()
//	})
//	// set tags in object
//	api.UserAPIPutObjectTagsHandler = user_api.PutObjectTagsHandlerFunc(func(params user_api.PutObjectTagsParams, session *models.Principal) middleware.Responder {
//		if err := getPutObjectTagsResponse(session, params); err != nil {
//			return user_api.NewPutObjectTagsDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPutObjectTagsOK()
//	})
//	//Restore file version
//	api.UserAPIPutObjectRestoreHandler = user_api.PutObjectRestoreHandlerFunc(func(params user_api.PutObjectRestoreParams, session *models.Principal) middleware.Responder {
//		if err := getPutObjectRestoreResponse(session, params); err != nil {
//			return user_api.NewPutObjectRestoreDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewPutObjectRestoreOK()
//	})
//	// Metadata in object
//	api.UserAPIGetObjectMetadataHandler = user_api.GetObjectMetadataHandlerFunc(func(params user_api.GetObjectMetadataParams, session *models.Principal) middleware.Responder {
//		resp, err := getObjectMetadataResponse(session, params)
//		if err != nil {
//			return user_api.NewGetObjectMetadataDefault(int(err.Code)).WithPayload(err)
//		}
//		return user_api.NewGetObjectMetadataOK().WithPayload(resp)
//	})
//}

// getListObjectsResponse returns a list of objects
func getListObjectsResponse(session *models.Principal, params models.ListObjectsParams) (*models.ListObjectsResponse, *models.Error) {
	var prefix string
	var recursive bool
	var withVersions bool
	var withMetadata bool
	if params.Prefix != nil {
		encodedPrefix := SanitizeEncodedPrefix(*params.Prefix)
		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
		if err != nil {
			return nil, prepareError(err)
		}
		prefix = string(decodedPrefix)
	}
	if params.Recursive != nil {
		recursive = *params.Recursive
	}
	if params.WithVersions != nil {
		withVersions = *params.WithVersions
	}
	if params.WithMetadata != nil {
		withMetadata = *params.WithMetadata
	}
	// bucket request needed to proceed
	if params.BucketName == "" {
		return nil, prepareError(errBucketNameNotInRequest)
	}
	mAdmin, err := NewMinioAdminClient(session)
	if err != nil {
		return nil, prepareError(err)
	}
	// create a minioClient interface implementation
	// defining the client to be used
	adminClient := AdminClient{Client: mAdmin}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	objs, err := listBucketObjects(ctx, adminClient, params.BucketName, prefix, recursive, withVersions, withMetadata)
	if err != nil {
		return nil, prepareError(err)
	}

	resp := &models.ListObjectsResponse{
		Objects: objs,
		Total:   int64(len(objs)),
	}
	return resp, nil
}

// listBucketObjects gets an array of objects in a bucket
func listBucketObjects(ctx context.Context, client AdminClient, bucketName string, prefix string, recursive, withVersions bool, withMetadata bool) ([]*object.Object, error) {
	var objects []*object.Object
	//opts := minio.ListObjectsOptions{
	//	Prefix:       prefix,
	//	Recursive:    recursive,
	//	WithVersions: withVersions,
	//	WithMetadata: withMetadata,
	//}
	//if withMetadata {
	//	opts.MaxKeys = 1
	//}
	listObjects, err := client.listObject(ctx, bucketName)
	if err != nil {
		return nil, err
	}
	for _, lsObj := range listObjects {
		objects = append(objects, &lsObj)
	}
	return objects, nil
}

//type httpRange struct {
//	Start  int64
//	Length int64
//}
//
//// Example:
////   "Content-Range": "bytes 100-200/1000"
////   "Content-Range": "bytes 100-200/*"
//func getRange(start, end, total int64) string {
//	// unknown total: -1
//	if total == -1 {
//		return fmt.Sprintf("bytes %d-%d/*", start, end)
//	}
//
//	return fmt.Sprintf("bytes %d-%d/%d", start, end, total)
//}
//
//// Example:
////   "Range": "bytes=100-200"
////   "Range": "bytes=-50"
////   "Range": "bytes=150-"
////   "Range": "bytes=0-0,-1"
//func parseRange(s string, size int64) ([]httpRange, error) {
//	if s == "" {
//		return nil, nil // header not present
//	}
//	const b = "bytes="
//	if !strings.HasPrefix(s, b) {
//		return nil, errors.New("invalid range")
//	}
//	var ranges []httpRange
//	for _, ra := range strings.Split(s[len(b):], ",") {
//		ra = strings.TrimSpace(ra)
//		if ra == "" {
//			continue
//		}
//		i := strings.Index(ra, "-")
//		if i < 0 {
//			return nil, errors.New("invalid range")
//		}
//		start, end := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
//		var r httpRange
//		if start == "" {
//			// If no start is specified, end specifies the
//			// range start relative to the end of the file.
//			i, err := strconv.ParseInt(end, 10, 64)
//			if err != nil {
//				return nil, errors.New("invalid range")
//			}
//			if i > size {
//				i = size
//			}
//			r.Start = size - i
//			r.Length = size - r.Start
//		} else {
//			i, err := strconv.ParseInt(start, 10, 64)
//			if err != nil || i >= size || i < 0 {
//				return nil, errors.New("invalid range")
//			}
//			r.Start = i
//			if end == "" {
//				// If no end is specified, range extends to end of the file.
//				r.Length = size - r.Start
//			} else {
//				i, err := strconv.ParseInt(end, 10, 64)
//				if err != nil || r.Start > i {
//					return nil, errors.New("invalid range")
//				}
//				if i >= size {
//					i = size - 1
//				}
//				r.Length = i - r.Start + 1
//			}
//		}
//		ranges = append(ranges, r)
//	}
//	return ranges, nil
//}

//func getDownloadObjectResponse(session *models.Principal, params user_api.DownloadObjectParams) (middleware.Responder, *models.Error) {
//	ctx := context.Background()
//	var prefix string
//	mClient, err := NewMinioAdminClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	client := AdminClient{Client: mClient}
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	resp, err := mClient.GetObject(ctx, params.BucketName, prefix, minio.GetObjectOptions{})
//	if err != nil {
//		return nil, prepareError(err)
//	}
//
//	return middleware.ResponderFunc(func(rw http.ResponseWriter, _ runtime.Producer) {
//		defer resp.Close()
//
//		isPreview := params.Preview != nil && *params.Preview
//		// indicate it's a download / inline content to the browser, and the size of the object
//		var filename string
//		prefixElements := strings.Split(prefix, "/")
//		if len(prefixElements) > 0 {
//			if prefixElements[len(prefixElements)-1] == "" {
//				filename = prefixElements[len(prefixElements)-2]
//			} else {
//				filename = prefixElements[len(prefixElements)-1]
//			}
//		}
//		escapedName := url.PathEscape(filename)
//
//		// indicate object size & content type
//		stat, err := resp.Stat()
//		if err != nil {
//			LogError("Failed to get Stat() response from server for %s: %v", prefix, err)
//			return
//		}
//
//		// if we are getting a Range Request (video) handle that specially
//		ranges, err := parseRange(params.HTTPRequest.Header.Get("Range"), stat.Size)
//		if err != nil {
//			LogError("Unable to parse range header input %s: %v", params.HTTPRequest.Header.Get("Range"), err)
//			return
//		}
//
//		if isPreview {
//			rw.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", escapedName))
//			rw.Header().Set("X-Frame-Options", "SAMEORIGIN")
//			rw.Header().Set("X-XSS-Protection", "1")
//		} else {
//			rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", escapedName))
//		}
//
//		rw.Header().Set("Last-Modified", stat.LastModified.UTC().Format(http.TimeFormat))
//
//		contentType := stat.ContentType
//		if isPreview {
//			// In case content type was uploaded as octet-stream, we double verify content type
//			if stat.ContentType == "application/octet-stream" {
//				contentType = mimedb.TypeByExtension(filepath.Ext(escapedName))
//			}
//		}
//		rw.Header().Set("Content-Type", contentType)
//		length := stat.Size
//		if len(ranges) > 0 {
//			start := ranges[0].Start
//			length = ranges[0].Length
//
//			_, err = resp.Seek(start, io.SeekStart)
//			if err != nil {
//				LogError("Unable to seek at offset %d: %v", start, err)
//				return
//			}
//
//			rw.Header().Set("Accept-Ranges", "bytes")
//			rw.Header().Set("Access-Control-Allow-Origin", "*")
//			rw.Header().Set("Content-Range", getRange(start, start+length-1, stat.Size))
//			rw.WriteHeader(http.StatusPartialContent)
//		}
//
//		rw.Header().Set("Content-Length", fmt.Sprintf("%d", length))
//		_, err = io.Copy(rw, io.LimitReader(resp, length))
//		if err != nil {
//			LogError("Unable to write all data to client: %v", err)
//			return
//		}
//	}), nil
//}
//func getDownloadFolderResponse(session *models.Principal, params user_api.DownloadObjectParams) (middleware.Responder, *models.Error) {
//	ctx := context.Background()
//	var prefix string
//	mClient, err := newMinioClient(session)
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	folders := strings.Split(prefix, "/")
//
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	minioClient := minioClient{client: mClient}
//	objects, err := listBucketObjects(ctx, minioClient, params.BucketName, prefix, true, false, false)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	w := new(bytes.Buffer)
//	zipw := zip.NewWriter(w)
//	var folder string
//	if len(folders) > 1 {
//		folder = folders[len(folders)-2]
//	}
//	for i := 0; i < len(objects); i++ {
//		name := folder + objects[i].Name[len(prefix)-1:]
//		object, err := mClient.GetObject(ctx, params.BucketName, objects[i].Name, minio.GetObjectOptions{})
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		f, err := zipw.Create(name)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		buf := new(bytes.Buffer)
//		buf.ReadFrom(object)
//		f.Write(buf.Bytes())
//	}
//	zipw.Close()
//	resp := io.NopCloser(bytes.NewReader(w.Bytes()))
//
//	return middleware.ResponderFunc(func(rw http.ResponseWriter, _ runtime.Producer) {
//		defer resp.Close()
//
//		// indicate it's a download / inline content to the browser, and the size of the object
//		var prefixPath string
//		var filename string
//		if params.Prefix != "" {
//			encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//			decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//			if err != nil {
//				LogError("Unable to parse encoded prefix %s: %v", encodedPrefix, err)
//				return
//			}
//
//			prefixPath = string(decodedPrefix)
//		}
//		prefixElements := strings.Split(prefixPath, "/")
//		if len(prefixElements) > 0 {
//			if prefixElements[len(prefixElements)-1] == "" {
//				filename = prefixElements[len(prefixElements)-2]
//			} else {
//				filename = prefixElements[len(prefixElements)-1]
//			}
//		}
//		escapedName := url.PathEscape(filename)
//
//		rw.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", escapedName))
//		rw.Header().Set("Content-Type", "application/zip")
//
//		// Copy the stream
//		_, err := io.Copy(rw, resp)
//		if err != nil {
//			LogError("Unable to write all the requested data: %v", err)
//		}
//	}), nil
//}
//
//// getDeleteObjectResponse returns whether there was an error on deletion of object
//func getDeleteObjectResponse(session *models.Principal, params user_api.DeleteObjectParams) *models.Error {
//	ctx := context.Background()
//	var prefix string
//	if params.Path != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Path)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	s3Client, err := newS3BucketClient(session, params.BucketName, prefix)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a mc S3Client interface implementation
//	// defining the client to be used
//	mcClient := mcClient{client: s3Client}
//	var rec bool
//	var version string
//	var allVersions bool
//	if params.Recursive != nil {
//		rec = *params.Recursive
//	}
//	if params.VersionID != nil {
//		version = *params.VersionID
//	}
//	if params.AllVersions != nil {
//		allVersions = *params.AllVersions
//	}
//	minClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	client2 := minioClient{client: minClient}
//	err = deleteObjects(ctx, mcClient, client2, params.BucketName, prefix, version, rec, allVersions)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//// getDeleteMultiplePathsResponse returns whether there was an error on deletion of any object
//func getDeleteMultiplePathsResponse(session *models.Principal, params user_api.DeleteMultipleObjectsParams) *models.Error {
//	ctx := context.Background()
//	var version string
//	var allVersions bool
//	if params.AllVersions != nil {
//		allVersions = *params.AllVersions
//	}
//	minClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	client2 := minioClient{client: minClient}
//	for i := 0; i < len(params.Files); i++ {
//		if params.Files[i].VersionID != "" {
//			version = params.Files[i].VersionID
//		}
//		prefix := params.Files[i].Path
//		s3Client, err := newS3BucketClient(session, params.BucketName, prefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		// create a mc S3Client interface implementation
//		// defining the client to be used
//		mcClient := mcClient{client: s3Client}
//		err = deleteObjects(ctx, mcClient, client2, params.BucketName, params.Files[i].Path, version, params.Files[i].Recursive, allVersions)
//		if err != nil {
//			return prepareError(err)
//		}
//	}
//	return nil
//}
//
//// deleteObjects deletes either a single object or multiple objects based on recursive flag
//func deleteObjects(ctx context.Context, client MCClient, client2 MinioClient, bucket string, path string, versionID string, recursive bool, allVersions bool) error {
//	if allVersions {
//		if recursive {
//			if err := deleteMultipleObjects(ctx, client, recursive, true); err != nil {
//				return err
//			}
//		} else {
//			objects, err := listBucketObjects(ctx, client2, bucket, path, recursive, true, false)
//			if err != nil {
//				return err
//			}
//			for i := range objects {
//				if err := deleteSingleObject(ctx, client, bucket, path, objects[i].VersionID); err != nil {
//					return err
//				}
//			}
//		}
//	}
//	if recursive {
//		if err := deleteMultipleObjects(ctx, client, recursive, false); err != nil {
//			return err
//		}
//	} else {
//		if err := deleteSingleObject(ctx, client, bucket, path, versionID); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//// deleteMultipleObjects uses listing before removal, it can list recursively or not,
////   Use cases:
////      * Remove objects recursively
//func deleteMultipleObjects(ctx context.Context, client MCClient, recursive bool, allVersions bool) error {
//	isRemoveBucket := false
//	isIncomplete := false
//	isBypass := false
//	listOpts := mc.ListOptions{Recursive: recursive, Incomplete: isIncomplete, ShowDir: mc.DirNone, WithOlderVersions: allVersions, WithDeleteMarkers: allVersions}
//	// TODO: support older Versions
//	contentCh := make(chan *mc.ClientContent, 1)
//
//	resultCh := client.remove(ctx, isIncomplete, isRemoveBucket, isBypass, contentCh)
//OUTER_LOOP:
//	for content := range client.list(ctx, listOpts) {
//		if content.Err != nil {
//			switch content.Err.ToGoError().(type) {
//			// ignore same as mc
//			case mc.PathInsufficientPermission:
//				// Ignore Permission error.
//				continue
//			}
//			close(contentCh)
//			return content.Err.Cause
//		}
//		sent := false
//		for !sent {
//			select {
//			case contentCh <- content:
//				sent = true
//			case result := <-resultCh:
//				if result.Err != nil {
//					switch result.Err.ToGoError().(type) {
//					// ignore same as mc
//					case mc.PathInsufficientPermission:
//						// Ignore Permission error.
//						continue
//					}
//					close(contentCh)
//					return result.Err.Cause
//				}
//				break OUTER_LOOP
//			}
//		}
//	}
//	close(contentCh)
//	for result := range resultCh {
//		if result.Err != nil {
//			switch result.Err.ToGoError().(type) {
//			// ignore same as mc
//			case mc.PathInsufficientPermission:
//				// Ignore Permission error.
//				continue
//			}
//			return result.Err.Cause
//		}
//	}
//	return nil
//
//}
//
//func deleteSingleObject(ctx context.Context, client MCClient, bucket, object string, versionID string) error {
//	targetURL := fmt.Sprintf("%s/%s", bucket, object)
//	contentCh := make(chan *mc.ClientContent, 1)
//	contentCh <- &mc.ClientContent{URL: *newClientURL(targetURL), VersionID: versionID}
//	close(contentCh)
//
//	isRemoveBucket := false
//	isIncomplete := false
//	isBypass := false
//
//	resultCh := client.remove(ctx, isIncomplete, isRemoveBucket, isBypass, contentCh)
//	for result := range resultCh {
//		if result.Err != nil {
//			switch result.Err.ToGoError().(type) {
//			// ignore same as mc
//			case mc.PathInsufficientPermission:
//				// Ignore Permission error.
//				continue
//			}
//			return result.Err.Cause
//		}
//	}
//	return nil
//}

func getUploadObjectResponse(session *models.Principal, params user_api.PostBucketsBucketNameObjectsUploadParams) *models.Error {
	ctx := context.Background()
	mClient, err := NewMinioAdminClient(session)
	if err != nil {
		return prepareError(err)
	}
	client := AdminClient{Client: mClient}
	if err := uploadFiles(ctx, client, params); err != nil {
		return prepareError(err, ErrorGeneric)
	}
	return nil
}

// uploadFiles gets files from http.Request form and uploads them to MinIO
func uploadFiles(ctx context.Context, client AdminClient, params user_api.PostBucketsBucketNameObjectsUploadParams) error {
	var prefix string
	if params.Prefix != nil {
		encodedPrefix := SanitizeEncodedPrefix(*params.Prefix)
		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
		if err != nil {
			return err
		}
		prefix = string(decodedPrefix)
	}

	// parse a request body as multipart/form-data.
	// 32 << 20 is default max memory
	mr, err := params.HTTPRequest.MultipartReader()
	if err != nil {
		return err
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		size, err := strconv.ParseInt(p.FormName(), 10, 64)
		if err != nil {
			return err
		}

		//contentType := p.Header.Get("content-type")

		err = client.putObject(ctx, params.BucketName, path.Join(prefix, p.FileName()), p, size)

		if err != nil {
			return err
		}
	}

	return nil
}

//// getShareObjectResponse returns a share object url
//func getShareObjectResponse(session *models.Principal, params user_api.ShareObjectParams) (*string, *models.Error) {
//	ctx := context.Background()
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	s3Client, err := newS3BucketClient(session, params.BucketName, prefix)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	// create a mc S3Client interface implementation
//	// defining the client to be used
//	mcClient := mcClient{client: s3Client}
//	var expireDuration string
//	if params.Expires != nil {
//		expireDuration = *params.Expires
//	}
//	url, err := getShareObjectURL(ctx, mcClient, params.VersionID, expireDuration)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	return url, nil
//}
//
//func getShareObjectURL(ctx context.Context, client MCClient, versionID string, duration string) (url *string, err error) {
//	// default duration 7d if not defined
//	if strings.TrimSpace(duration) == "" {
//		duration = "168h"
//	}
//
//	expiresDuration, err := time.ParseDuration(duration)
//	if err != nil {
//		return nil, err
//	}
//	objURL, pErr := client.shareDownload(ctx, versionID, expiresDuration)
//	if pErr != nil {
//		return nil, pErr.Cause
//	}
//	return &objURL, nil
//}
//
//func getSetObjectLegalHoldResponse(session *models.Principal, params user_api.PutObjectLegalHoldParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = setObjectLegalHold(ctx, minioClient, params.BucketName, prefix, params.VersionID, *params.Body.Status)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func setObjectLegalHold(ctx context.Context, client MinioClient, bucketName, prefix, versionID string, status models.ObjectLegalHoldStatus) error {
//	var lstatus minio.LegalHoldStatus
//	if status == models.ObjectLegalHoldStatusEnabled {
//		lstatus = minio.LegalHoldEnabled
//	} else {
//		lstatus = minio.LegalHoldDisabled
//	}
//	return client.putObjectLegalHold(ctx, bucketName, prefix, minio.PutObjectLegalHoldOptions{VersionID: versionID, Status: &lstatus})
//}
//
//func getSetObjectRetentionResponse(session *models.Principal, params user_api.PutObjectRetentionParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = setObjectRetention(ctx, minioClient, params.BucketName, params.VersionID, prefix, params.Body)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func setObjectRetention(ctx context.Context, client MinioClient, bucketName, versionID, prefix string, retentionOps *models.PutObjectRetentionRequest) error {
//	if retentionOps == nil {
//		return errors.New("object retention options can't be nil")
//	}
//	if retentionOps.Expires == nil {
//		return errors.New("object retention expires can't be nil")
//	}
//
//	var mode minio.RetentionMode
//	if *retentionOps.Mode == models.ObjectRetentionModeGovernance {
//		mode = minio.Governance
//	} else {
//		mode = minio.Compliance
//	}
//	retentionUntilDate, err := time.Parse(time.RFC3339, *retentionOps.Expires)
//	if err != nil {
//		return err
//	}
//	opts := minio.PutObjectRetentionOptions{
//		GovernanceBypass: retentionOps.GovernanceBypass,
//		RetainUntilDate:  &retentionUntilDate,
//		Mode:             &mode,
//		VersionID:        versionID,
//	}
//	return client.putObjectRetention(ctx, bucketName, prefix, opts)
//}
//
//func deleteObjectRetentionResponse(session *models.Principal, params user_api.DeleteObjectRetentionParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = deleteObjectRetention(ctx, minioClient, params.BucketName, prefix, params.VersionID)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func deleteObjectRetention(ctx context.Context, client MinioClient, bucketName, prefix, versionID string) error {
//	opts := minio.PutObjectRetentionOptions{
//		GovernanceBypass: true,
//		VersionID:        versionID,
//	}
//
//	return client.putObjectRetention(ctx, bucketName, prefix, opts)
//}
//
//func getPutObjectTagsResponse(session *models.Principal, params user_api.PutObjectTagsParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//	err = putObjectTags(ctx, minioClient, params.BucketName, prefix, params.VersionID, params.Body.Tags)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func putObjectTags(ctx context.Context, client MinioClient, bucketName, prefix, versionID string, tagMap map[string]string) error {
//	opt := minio.PutObjectTaggingOptions{
//		VersionID: versionID,
//	}
//	otags, err := tags.MapToObjectTags(tagMap)
//	if err != nil {
//		return err
//	}
//	return client.putObjectTagging(ctx, bucketName, prefix, otags, opt)
//}
//
//// Restore Object Version
//func getPutObjectRestoreResponse(session *models.Principal, params user_api.PutObjectRestoreParams) *models.Error {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//
//	var prefix string
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	err = restoreObject(ctx, minioClient, params.BucketName, prefix, params.VersionID)
//	if err != nil {
//		return prepareError(err)
//	}
//	return nil
//}
//
//func restoreObject(ctx context.Context, client MinioClient, bucketName, prefix, versionID string) error {
//	// Select required version
//	srcOpts := minio.CopySrcOptions{
//		Bucket:    bucketName,
//		Object:    prefix,
//		VersionID: versionID,
//	}
//
//	// Destination object, same as current bucket
//	replaceMetadata := make(map[string]string)
//	replaceMetadata["copy-source"] = versionID
//
//	dstOpts := minio.CopyDestOptions{
//		Bucket:       bucketName,
//		Object:       prefix,
//		UserMetadata: replaceMetadata,
//	}
//
//	// Copy object call
//	_, err := client.copyObject(ctx, dstOpts, srcOpts)
//
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//// Metadata Response from minio-go API
//func getObjectMetadataResponse(session *models.Principal, params user_api.GetObjectMetadataParams) (*models.Metadata, *models.Error) {
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
//	defer cancel()
//	mClient, err := newMinioClient(session)
//	if err != nil {
//		return nil, prepareError(err)
//	}
//	// create a minioClient interface implementation
//	// defining the client to be used
//	minioClient := minioClient{client: mClient}
//	var prefix string
//
//	if params.Prefix != "" {
//		encodedPrefix := SanitizeEncodedPrefix(params.Prefix)
//		decodedPrefix, err := base64.StdEncoding.DecodeString(encodedPrefix)
//		if err != nil {
//			return nil, prepareError(err)
//		}
//		prefix = string(decodedPrefix)
//	}
//
//	objectInfo, err := getObjectInfo(ctx, minioClient, params.BucketName, prefix)
//
//	if err != nil {
//		return nil, prepareError(err)
//	}
//
//	metadata := &models.Metadata{ObjectMetadata: objectInfo.Metadata}
//
//	return metadata, nil
//}
//
//func getObjectInfo(ctx context.Context, client MinioClient, bucketName, prefix string) (minio.ObjectInfo, error) {
//	objectData, err := client.statObject(ctx, bucketName, prefix, minio.GetObjectOptions{})
//
//	if err != nil {
//		return minio.ObjectInfo{}, err
//	}
//
//	return objectData, nil
//}
//
//// newClientURL returns an abstracted URL for filesystems and object storage.
//func newClientURL(urlStr string) *mc.ClientURL {
//	scheme, rest := getScheme(urlStr)
//	if strings.HasPrefix(rest, "//") {
//		// if rest has '//' prefix, skip them
//		var authority string
//		authority, rest = splitSpecial(rest[2:], "/", false)
//		if rest == "" {
//			rest = "/"
//		}
//		host := getHost(authority)
//		if host != "" && (scheme == "http" || scheme == "https") {
//			return &mc.ClientURL{
//				Scheme:          scheme,
//				Type:            objectStorage,
//				Host:            host,
//				Path:            rest,
//				SchemeSeparator: "://",
//				Separator:       '/',
//			}
//		}
//	}
//	return &mc.ClientURL{
//		Type:      fileSystem,
//		Path:      rest,
//		Separator: filepath.Separator,
//	}
//}
//
//// Maybe rawurl is of the form scheme:path. (Scheme must be [a-zA-Z][a-zA-Z0-9+-.]*)
//// If so, return scheme, path; else return "", rawurl.
//func getScheme(rawurl string) (scheme, path string) {
//	urlSplits := strings.Split(rawurl, "://")
//	if len(urlSplits) == 2 {
//		scheme, uri := urlSplits[0], "//"+urlSplits[1]
//		// ignore numbers in scheme
//		validScheme := regexp.MustCompile("^[a-zA-Z]+$")
//		if uri != "" {
//			if validScheme.MatchString(scheme) {
//				return scheme, uri
//			}
//		}
//	}
//	return "", rawurl
//}
//
//// Assuming s is of the form [s delimiter s].
//// If so, return s, [delimiter]s or return s, s if cutdelimiter == true
//// If no delimiter found return s, "".
//func splitSpecial(s string, delimiter string, cutdelimiter bool) (string, string) {
//	i := strings.Index(s, delimiter)
//	if i < 0 {
//		// if delimiter not found return as is.
//		return s, ""
//	}
//	// if delimiter should be removed, remove it.
//	if cutdelimiter {
//		return s[0:i], s[i+len(delimiter):]
//	}
//	// return split strings with delimiter
//	return s[0:i], s[i:]
//}
//
//// getHost - extract host from authority string, we do not support ftp style username@ yet.
//func getHost(authority string) (host string) {
//	i := strings.LastIndex(authority, "@")
//	if i >= 0 {
//		// TODO support, username@password style userinfo, useful for ftp support.
//		return
//	}
//	return authority
//}