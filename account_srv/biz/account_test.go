package biz

import (
	"context"
	"fmt"
	"github.com/baojunGit/goAdmin/account_srv/initialize"
	"github.com/baojunGit/goAdmin/account_srv/proto/pb"
	"testing"
)

func init() {
	initialize.InitDB("../config")
}

func TestAccountServer_AddAccount(t *testing.T) {
	accountServer := AccountServer{}
	for i := 0; i < 5; i++ {
		s := fmt.Sprintf("1300000000%d", i)
		res, err := accountServer.AddAccount(context.Background(), &pb.AddAccountRequest{
			Mobile:   s,
			Password: s,
			Nickname: s,
			Gender:   "male",
		})
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(res.Id)
	}
}

func TestAccountServer_UpdateAccount(t *testing.T) {
	accountServer := AccountServer{}
	req := pb.UpdateAccountRequest{
		Id:       1,
		Mobile:   "18006996859",
		Password: "18006996859",
		Nickname: "18006996859",
		Gender:   "female",
		Role:     2,
	}
	res, err := accountServer.UpdateAccount(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)
}

func TestAccountServer_GetAccountList(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   1,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}
	res, err = accountServer.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   2,
		PageSize: 3,
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, account := range res.AccountList {
		fmt.Println(account.Id)
	}
}

func TestAccountServer_GetAccountByMobile(t *testing.T) {
	mobile := "13000000000"
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: mobile})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Id)
}

func TestAccountServer_GetAccountById(t *testing.T) {
	id := 3
	accountServer := AccountServer{}
	res, err := accountServer.GetAccountById(context.Background(), &pb.IdRequest{Id: uint32(id)})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Mobile)
}

func TestAccountServer_CheckPassword(t *testing.T) {
	accountServer := AccountServer{}
	res, err := accountServer.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       "13000000004",
		HashedPassword: "199a7413cf3e7620928f68c960b4d261d87ff30a43d2de96c838b8d1e05e4747",
		AccountId:      5,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)
	res, err = accountServer.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       "13000000001",
		HashedPassword: "199a7413cf3e7620928f68c960b4d261d87ff30a43d2de96c838b8d1e05e4747",
		AccountId:      5,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res.Result)
}
