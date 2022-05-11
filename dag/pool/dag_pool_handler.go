package pool

import (
	"github.com/filedag-project/filedag-storage/dag/pool/userpolicy"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-merkledag"
	"io/ioutil"
	"net/http"
)

func (d *DagPoolServer) Add(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	secretKey := vars["secretKey"]
	if !d.dp.Iam.CheckUserPolicy(user, secretKey, userpolicy.OnlyWrite) {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	all, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	sum, err := d.dp.CidBuilder.Sum(all)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = d.dp.Add(r.Context(), merkledag.NewRawNode(all))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Write([]byte(sum.String()))
}

func (d *DagPoolServer) Get(w http.ResponseWriter, r *http.Request) {

}

func (d *DagPoolServer) Delete(w http.ResponseWriter, r *http.Request) {

}
