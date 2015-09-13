package routes

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func IndexIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "Index!")
}

func IndexStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
