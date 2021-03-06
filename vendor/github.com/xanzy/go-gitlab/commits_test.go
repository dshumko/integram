package gitlab

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGetCommitStatuses(t *testing.T) {
	mux, server, client := setup()
	defer teardown(server)

	mux.HandleFunc("/projects/1/repository/commits/b0b3a907f41409829b307a28b82fdbd552ee5a27/statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"ref":   "master",
			"stage": "test",
			"name":  "ci/jenkins",
			"all":   "true",
		})
		fmt.Fprint(w, `[{"id":1}]`)
	})

	opt := &GetCommitStatusesOptions{"master", "test", "ci/jenkins", true}
	statuses, _, err := client.Commits.GetCommitStatuses("1", "b0b3a907f41409829b307a28b82fdbd552ee5a27", opt)

	if err != nil {
		t.Errorf("Commits.GetCommitStatuses returned error: %v", err)
	}

	want := []*CommitStatus{{ID: 1}}
	if !reflect.DeepEqual(want, statuses) {
		t.Errorf("Commits.GetCommitStatuses returned %+v, want %+v", statuses, want)
	}
}

func TestSetCommitStatus(t *testing.T) {
	mux, server, client := setup()
	defer teardown(server)

	mux.HandleFunc("/projects/1/statuses/b0b3a907f41409829b307a28b82fdbd552ee5a27", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		testJsonBody(t, r, values{
			"state":       "running",
			"ref":         "master",
			"name":        "ci/jenkins",
			"target_url":  "http://abc",
			"description": "build",
		})
		fmt.Fprint(w, `{"id":1}`)
	})

	opt := &SetCommitStatusOptions{Running, "master", "ci/jenkins", "", "http://abc", "build"}
	status, _, err := client.Commits.SetCommitStatus("1", "b0b3a907f41409829b307a28b82fdbd552ee5a27", opt)

	if err != nil {
		t.Errorf("Commits.SetCommitStatus returned error: %v", err)
	}

	want := &CommitStatus{ID: 1}
	if !reflect.DeepEqual(want, status) {
		t.Errorf("Commits.SetCommitStatus returned %+v, want %+v", status, want)
	}
}
