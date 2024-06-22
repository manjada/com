package mjd

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/manjada/com/log"
	"reflect"
	"testing"
)

func TestElastic_Create(t *testing.T) {
	type fields struct {
		Index   string
		Refresh string
	}
	type args struct {
		document interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "Test create", fields: struct {
			Index   string
			Refresh string
		}{Index: "ocky_dev", Refresh: "true"}, args: args{document: struct {
			Id   string
			Nama string
		}{
			Id:   "121212112121",
			Nama: "ocky",
		}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := NewElastic()
			if err != nil {
				log.Error(err)
			}
			e = &Elastic{
				Index:   tt.fields.Index,
				Refresh: tt.fields.Refresh,
			}

			if res, err := e.Create(tt.args.document); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				log.Info("data ", res.Result)
			}
		})
	}
}

func TestElastic_GetByDocId(t *testing.T) {
	type fields struct {
		Index   string
		Refresh string
	}
	type args struct {
		documentId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *get.Response
		wantErr bool
	}{
		{name: "TEST_GET", fields: struct {
			Index   string
			Refresh string
		}{Index: "ocky_dev", Refresh: ""}, args: args{documentId: "DUM9RY0B4aBs_QQ0-H_h"}, want: nil, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elastic{
				Index:   tt.fields.Index,
				Refresh: tt.fields.Refresh,
			}
			e, err := NewElastic()
			got, err := e.GetByDocId(tt.args.documentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByDocId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByDocId() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElastic_IsDocumentExist(t *testing.T) {
	type fields struct {
		Index   string
		Refresh string
	}
	type args struct {
		documentId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elastic{
				Index:   tt.fields.Index,
				Refresh: tt.fields.Refresh,
			}
			got, err := e.IsDocumentExist(tt.args.documentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDocumentExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDocumentExist() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElastic_Search(t *testing.T) {
	type fields struct {
		Index   string
		Refresh string
	}
	type args struct {
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *search.Response
		wantErr bool
	}{
		{name: "SEARCH_DATA", fields: struct {
			Index   string
			Refresh string
		}{Index: "ocky_dev"}, args: struct{ data map[string]interface{} }{data: map[string]interface{}{}}, want: &search.Response{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Elastic{
				Index:   tt.fields.Index,
				Refresh: tt.fields.Refresh,
			}
			e, err := NewElastic()
			got, err := e.Search(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
