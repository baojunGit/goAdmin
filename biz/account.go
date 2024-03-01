package biz

import (
	"context"
	"crypto/md5"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/baojunGit/goAdmin/exception"
	"github.com/baojunGit/goAdmin/initialize"
	"github.com/baojunGit/goAdmin/model"
	"github.com/baojunGit/goAdmin/proto/pb"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type AccountServer struct {
	// mustEmbedUnimplementedAccountServiceServer 该方法是一个占位符没有实际的逻辑，只是为了满足编译器对接口的要求
	pb.UnimplementedAccountServiceServer
}

// Paginate 分页方法
func Paginate(pageNo, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageNo == 0 {
			pageNo = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		// mysql分页
		offset := (pageNo - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func Model2Pb(account model.Account) *pb.AccountRes {
	accountRes := &pb.AccountRes{
		Id:       int32(account.ID),
		Mobile:   account.Mobile,
		Password: account.Password,
		Nickname: account.NickName,
		Gender:   account.Gender,
		Role:     uint32(account.Role),
	}
	return accountRes
}

// GetAccountList ：a *AccountServer是一个接收者，表示这个方法是与类型为 AccountServer 的指针关联的
func (a *AccountServer) GetAccountList(ctx context.Context, req *pb.PagingRequest) (*pb.AccountListRes, error) {
	var accountList []model.Account
	//result := initialize.DB.Find(&accountList) 可以用来测试
	result := initialize.DB.Scopes(Paginate(int(req.PageNo), int(req.PageSize))).Find(&accountList)
	if result.Error != nil {
		return nil, result.Error
	}
	accountListRes := &pb.AccountListRes{}
	accountListRes.Total = int32(result.RowsAffected)
	for _, account := range accountList {
		accountRes := Model2Pb(account)
		accountListRes.AccountList = append(accountListRes.AccountList, accountRes)
	}
	return accountListRes, nil
}

func (a *AccountServer) GetAccountByMobile(ctx context.Context, req *pb.MobileRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := initialize.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 0 {
		return nil, errors.New(exception.AccountNotFound)
	}
	res := Model2Pb(account)
	return res, nil
}

func (a *AccountServer) GetAccountById(ctx context.Context, req *pb.IdRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := initialize.DB.First(&account, req.Id)
	if result.RowsAffected == 0 {
		return nil, errors.New(exception.AccountNotFound)
	}
	res := Model2Pb(account)
	return res, nil
}
func (a *AccountServer) AddAccount(ctx context.Context, req *pb.AddAccountRequest) (*pb.AccountRes, error) {
	var account model.Account
	result := initialize.DB.Where(&model.Account{Mobile: req.Mobile}).First(&account)
	if result.RowsAffected == 1 {
		return nil, errors.New(exception.AccountExists)
	}
	account.Mobile = req.Mobile
	account.NickName = req.Nickname
	account.Role = 1
	options := password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: md5.New,
	}
	salt, encodePwd := password.Encode(req.Password, &options)
	account.Salt = salt
	account.Password = encodePwd
	r := initialize.DB.Create(&account)
	if r.Error != nil {
		return nil, errors.New(exception.InternalError)
	}
	accountRes := Model2Pb(account)
	return accountRes, nil
}
func (a *AccountServer) UpdateAccount(ctx context.Context, req *pb.UpdateAccountRequest) (*pb.UpdateAccountRes, error) {
	return &pb.UpdateAccountRes{Result: true}, nil
}
func (a *AccountServer) CheckPassword(ctx context.Context, req *pb.CheckPasswordRequest) (*pb.CheckPasswordRes, error) {
	return &pb.CheckPasswordRes{Result: true}, nil
}

// mustEmbedUnimplementedAccountServiceServer 该方法是一个占位符没有实际的逻辑，只是为了满足编译器对接口的要求
//func (a *AccountServer) mustEmbedUnimplementedAccountServiceServer() {
//	// 实现 mustEmbedUnimplementedAccountServiceServer 方法的逻辑
//}
