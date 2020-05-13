package server

import (
	"reflect"
	"testing"
	"time"
)

func Test_crawlUserInfo(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name     string
		args     args
		wantInfo *userInfo
		wantErr  bool
	}{
		{
			"开放",
			args{"MORnlight"},
			&userInfo{"mornlight", true, true},
			false,
		},
		{
			"开放",
			args{"livid"},
			&userInfo{"Livid", true, true},
			false,
		},
		{
			"全部隐藏",
			args{"gBIn"},
			&userInfo{"gbin", false, true},
			false,
		},
		{
			"登录后可见",
			args{"morethansean"},
			&userInfo{"morethansean", false, true},
			false,
		},
		{
			"不存在的用户",
			args{"mornlightmornlight"},
			&userInfo{"", false, false},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := crawlUserInfo(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("crawlUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("crawlUserInfo() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
}

func Test_getUserInfo(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name     string
		args     args
		wantInfo *userInfo
		wantErr  bool
	}{
		{
			"开放",
			args{" MORnlight "},
			&userInfo{"mornlight", true, true},
			false,
		},
		{
			"开放",
			args{"livid"},
			&userInfo{"Livid", true, true},
			false,
		},
		{
			"全部隐藏",
			args{"gBIn"},
			&userInfo{"gbin", false, true},
			false,
		},
		{
			"登录后可见",
			args{"morethansean"},
			&userInfo{"morethansean", false, true},
			false,
		},
		{
			"不存在的用户",
			args{"mornlightmornLIGHT"},
			&userInfo{"", false, false},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := getUserInfo(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("getUserInfo() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}

	start := time.Now()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInfo, err := getUserInfo(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotInfo, tt.wantInfo) {
				t.Errorf("getUserInfo() gotInfo = %v, want %v", gotInfo, tt.wantInfo)
			}
		})
	}
	t.Run("已缓存用户耗时", func(t *testing.T) {
		elapsed := time.Since(start)
		if elapsed.Milliseconds() > 10 {
			t.Error("getUserInfo() used too many time when caches exist")
		}
	})
}

func Test_parseItemsInParam(t *testing.T) {
	type args struct {
		param string
	}
	tests := []struct {
		name         string
		args         args
		wantItems    []string
		wantExcluded bool
	}{
		{
			name: "one",
			args: args{
				param: "create",
			},
			wantItems: []string{"create"},
			wantExcluded: false,
		},
		{
			name: "exclude",
			args: args{
				param: "-create, go ",
			},
			wantItems: []string{"create", "go"},
			wantExcluded: true,
		},
		{
			name: "exclude multiple",
			args: args{
				param: "-create, go , -jobs",
			},
			wantItems: []string{"create", "go", "jobs"},
			wantExcluded: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, gotExcluded := parseItemsInParam(tt.args.param)
			if !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("parseItemsInParam() gotItems = %v, want %v", gotItems, tt.wantItems)
			}
			if gotExcluded != tt.wantExcluded {
				t.Errorf("parseItemsInParam() gotExcluded = %v, want %v", gotExcluded, tt.wantExcluded)
			}
		})
	}
}