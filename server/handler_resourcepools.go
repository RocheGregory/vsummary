package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gbolo/vsummary/common"
	"gopkg.in/go-playground/validator.v9"
)

func handlerResourcepool(w http.ResponseWriter, req *http.Request) {

	// log time on debug
	defer common.ExecutionTime(time.Now(), "handle")

	// read in body
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("error reading request body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Body.Close()

	// decode json body
	var reqStruct []common.ResourcePool
	err = json.Unmarshal(reqBody, &reqStruct)
	if err != nil {
		log.Errorf("failed to decode body: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate
	// TODO: fail if any obj is invalid?
	validate := validator.New()
	for _, o := range reqStruct {
		err = validate.Struct(o)
		if err != nil {
			log.Errorf("failed to validate body: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// insert to backend
	err = backend.InsertResourcepools(reqStruct)
	if err != nil {
		log.Errorf("failed to insert resourcepool: %s", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}
