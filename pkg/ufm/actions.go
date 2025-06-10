package ufm

import (
//	"bytes"
	"errors"
	"time"
	"fmt"
	"net/http"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"io"
	//"strconv"
	"strings"
	// "os"
	"context"
)

const ActionsPath = "/ufmRestV2/actions"
const JobsPath = "/ufmRestV2/jobs"

// get cables info action
// payload example:
//{"params":{"port_id":"abc123_16"},"action":"get_cables_info","object_ids":["abc123"],"object_type":"System","description":"","identifier":"id"}
func (u *UfmClient) GetCablesInfo(portId string) (cableInfo string, err error) {
	//submit action 
	payload := `{"params":{"port_id":""},"action":"get_cables_info","object_ids":[],"object_type":"System","description":"","identifier":"id"}`
	payload, _ = sjson.Set(payload, "params.port_id", portId)
	objectId := strings.Split(portId, "_")[0]
	payload, _ = sjson.Set(payload, "object_ids.0", objectId)
	//fmt.Println("payload:", payload)
	jobId, err := u.Action(payload)
	if err != nil {
		return
	}
	return u.WaitJobResult(jobId, 10)
}

// submit action, get a job status path
func (u *UfmClient) Action(payload string) (jobId string, err error) {
	path := ActionsPath
	resp, err := u.Post(path, "", payload)
	if err != nil {
		return
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusAccepted {
		// Extract a specific header (e.g., "Content-Type")
		jobStatusPath := resp.Header.Get("Location")
		jobStatusPathSplit := strings.Split(jobStatusPath, "/")
		jobId = jobStatusPathSplit[len(jobStatusPathSplit)-1]
	} else {
		err = errors.New(fmt.Sprintf("Error submitting action: %s", string(bodyBytes)))
	}
	return
		
}

// get job status loop (wait for outcome, with timeout)
func (u *UfmClient) WaitJobResult(jobId string, deadline int) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(deadline)*time.Second)	
	defer cancel()
	ticker := time.NewTicker(1* time.Second)
	defer ticker.Stop()
	getJob := func(jobId string)(ret string, err error) {
		resp, err := u.Get(JobsPath, []string{"parent_id="+jobId})
		if err != nil {
			return
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			err = errors.New(fmt.Sprintf("Error retrieving job %s: %s", jobId, string(bodyBytes)))
		} else {
			//fmt.Println("jobdata:", string(bodyBytes))
			ret = string(bodyBytes)
		}
		return
	}
	var job string
	for {
		select {
		case <- ctx.Done():
			err = errors.New(fmt.Sprintf("Gave up trying to get result for %s after %d seconds\n", jobId, deadline))
			return
		case <- ticker.C:
			//fmt.Println("Checking job status:", jobId)
			job, err = getJob(jobId)
			if err != nil {
				return
			} else if gjson.Parse(job).Array()[0].Get("Status").String() == "Completed" {
				result = gjson.Parse(job).String()
				return
			}
		}
	}		
}

