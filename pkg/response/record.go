package response

import (
	"encoding/json"
	"time"

	"GoRestify/pkg/activity"
	"GoRestify/pkg/models"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_types"
)

// RecordCreateInstant make it simpler for calling the record
func (r *Response) RecordCreateInstant(ev pkg_types.Event, newData interface{}) {
	r.Record(ev, nil, newData)
}

// Record will send the activity for read/update/delete to the ActivityCh
func (r *Response) Record(ev pkg_types.Event, data ...interface{}) {
	r.initiateRecordCh(ev, data...)
}

func (r *Response) initiateRecordCh(ev pkg_types.Event, data ...interface{}) {

	if !r.Engine.ActivityActive {
		return
	}

	var userID uint
	recordType := r.findRecordType(data...)
	before, after := r.fillBeforeAfter(recordType, data...)

	if id, ok := r.Context.Get("USER_ID"); ok {
		userID = id.(uint)
	}

	var Username string
	if username, ok := r.Context.Get("USERNAME"); ok {
		Username = username.(string)
	}

	activity := models.Activity{
		Event:      ev.String(),
		OperatorID: userID,
		Username:   Username,
		IP:         r.Context.Request.Header.Get("X-User-IP"),
		URI:        r.Context.Request.RequestURI,
		Before:     string(before),
		After:      string(after),
	}

	r.Engine.ActivityCh <- activity

	_ = activity

}

// RecordType is and int used as an enum
type RecordType int

const (
	read RecordType = iota
	writeBefore
	writeAfter
	writeBoth
)

// findRecordType is helper function for finding the best way for recording data
func (r *Response) findRecordType(data ...interface{}) RecordType {
	switch len(data) {
	case 0:
		return read
	case 2:
		return writeBoth
	default:
		if data[0] == nil {
			return writeAfter
		}
	}

	return writeBefore
}

// fillBeforeAfter check if there is a need for entering before data or not
func (r *Response) fillBeforeAfter(recordType RecordType, data ...interface{}) (before, after []byte) {
	var err error
	if recordType == writeBefore || recordType == writeBoth {
		before, err = json.Marshal(data[0])
		pkg_log.CheckError(err, "error in encoding data to before-json")
	}
	if recordType == writeAfter || recordType == writeBoth {
		after, err = json.Marshal(data[1])
		pkg_log.CheckError(err, "error in encoding data to after-json")
	}

	return
}

// ActivityWatcher is used for watching activity channel
func ActivityWatcher() {
	var arr []models.Activity
	counter := 0
	var activityModel models.Activity

	// after n sec activity will be inserted
	tickTimer := time.Tick(10 * time.Second)

	for {
		select {
		case activityModel = <-pkg_config.Config.ActivityCh:
			counter++
			arr = append(arr, activityModel)
			if counter > 50 {
				batch, err := activity.CreateBatch(arr)
				if err != nil {
					pkg_log.Error("activity watcher error (counter): ", batch, err)
				}
				counter = 0
				arr = []models.Activity{}
			}
		case <-tickTimer:
			if len(arr) > 0 {
				batch, err := activity.CreateBatch(arr)
				if err != nil {
					pkg_log.Error("activity watcher error (ticker): ", batch, err)
				}
				counter = 0
				arr = []models.Activity{}
			}
		}
	}
}
