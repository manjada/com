package mjd

import (
	"github.com/manjada/com/dto"
	"github.com/manjada/com/model"
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
		user model.User
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.TokenDetails
		wantErr bool
	}{
		{name: "Create Token", args: struct{ user model.User }{user: model.User{
			TransactionModel: model.TransactionModel{Id: "test123"},
			Roles: []model.Role{{TransactionModel: model.TransactionModel{}, Name: "CLIENT", Code: "CLT"}},
		}}, want: &dto.TokenDetails{}, wantErr: false},
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
