package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lxdumb/calcd/calc"
	"github.com/MadAppGang/httplog"
)

func calcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer w.Write([]byte("\n"))
	if r.Method != http.MethodPost {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	var req map[string]string
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	res, calcerr := calc.Calc(req["expression"])
	if calcerr != nil {
		if calcerr.Error() == "expression is not valid" {
			w.WriteHeader(422)
			w.Write([]byte("{\n    \"error\": \"expression is not valid\"\n}"))
			return
		}
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	w.WriteHeader(200)
	resmap := make(map[string]string)
	resb := strconv.FormatFloat(res, 'f', -1, 64)
	resmap["result"] = resb
	resjson, err2 := json.MarshalIndent(resmap, "", "    ")
	if err2 != nil {
		w.WriteHeader(500)
		w.Write([]byte("{\n    \"error\": \"Internal server error\"\n}"))
		return
	}
	w.Write(resjson)
}

var calclogHandler http.Handler = http.HandlerFunc(calcHandler)

func startserv() {
	shortLoggedHandler := httplog.LoggerWithFormatter(
		httplog.ShortLogFormatter,
	)
	http.Handle("/api/v1/calculate", shortLoggedHandler(calclogHandler))
	http.ListenAndServe(":8080", nil)
}
