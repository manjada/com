package auth

import (
	"github.com/manjada/com/dto"
	"reflect"
	"testing"
)

func TestCreateAuth(t *testing.T) {
	type args struct {
		userId string
		td     *dto.TokenDetails
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateAuth(tt.args.userId, tt.args.td); (err != nil) != tt.wantErr {
				t.Errorf("CreateAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateToken(t *testing.T) {
	type args struct {
		user dto.UserToken
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.TokenDetails
		wantErr bool
	}{
		{name: "Create Token", args: struct{ user dto.UserToken }{user: dto.UserToken{}}, want: &dto.TokenDetails{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
