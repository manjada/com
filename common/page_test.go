package common

import (
	"github.com/manjada/com/dto"
	"github.com/manjada/com/web"
	"github.com/manjada/com/web/mock"
	"testing"
)

func TestBuildPage(t *testing.T) {
	type args struct {
		c web.Context
	}
	tests := []struct {
		name string
		args args
		want dto.Page
	}{
		{
			name: "Test with all parameters",
			args: args{
				c: &mock.MockFiberCtx{
					QueriesMap: map[string]string{
						"page":    "2",
						"size":    "10",
						"search":  "test",
						"filter1": "value2",
						"filter2": "value3",
					},
				},
			},
			want: dto.Page{
				Page:   2,
				Size:   10,
				Search: "test",
				Filter: map[string]string{
					"filter1": "value1",
					"filter2": "value2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPage(tt.args.c)
			t.Logf("BuildPage() = %v, want %v", got, tt.want)
		})
	}
}
