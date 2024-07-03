package handler

import (
	"context"
	"fmt"
	"github.com/baojunGit/goAdmin/account_srv/proto/pb"
	"github.com/baojunGit/goAdmin/account_web/req"
	"github.com/baojunGit/goAdmin/account_web/res"
	"github.com/baojunGit/goAdmin/exception"
	"github.com/baojunGit/goAdmin/jwt_op"
	"github.com/baojunGit/goAdmin/log"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"strconv"
	"time"
)

func HandleError(err error) string {
	if err != nil {
		switch err.Error() {
		case exception.AccountExists:
			return exception.AccountExists
		case exception.AccountNotFound:
			return exception.AccountNotFound
		case exception.SaltError:
			return exception.SaltError
		default:
			return exception.InternalError
		}
	}
	return ""
}

func AccountListHandler(c *gin.Context) {
	pageNoStr := c.DefaultQuery("pageNo", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "3")
	//log.Logger.Info("AccountListHandler调试通过")
	//insecure.NewCredentials() 创建了一个不安全的凭证实例，这意味着不会对传输层进行加密或认证，这种方法在进行本地测试或在内部网络中不需要传输安全的情况下非常有用
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s := fmt.Sprintf("AccountListHandler-GRPC拨号失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	pageNo, _ := strconv.ParseInt(pageNoStr, 10, 32)
	pageSize, _ := strconv.ParseInt(pageSizeStr, 10, 32)
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountList(context.Background(), &pb.PagingRequest{
		PageNo:   uint32(pageNo),
		PageSize: uint32(pageSize),
	})
	if err != nil {
		s := fmt.Sprintf("GetAccountList调用失败:%s", err.Error())
		log.Logger.Info(s)
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	var resList []res.Account4Res
	for _, item := range r.AccountList {
		resList = append(resList, pb2res(item))
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "ok",
		"total": r.Total,
		"data":  resList,
	})
}

func pb2res(accountRes *pb.AccountRes) res.Account4Res {
	return res.Account4Res{
		Mobile:   accountRes.Mobile,
		NickName: accountRes.Nickname,
		Gender:   accountRes.Gender,
	}
}

func LoginByPasswordHandler(c *gin.Context) {
	var loginByPassword req.LoginByPassword
	err := c.ShouldBind(&loginByPassword)
	if err != nil {
		log.Logger.Error("LoginByPassword出错：" + err.Error())
		c.JSON(http.StatusOK, gin.H{
			"msg": "解析参数错误",
		})
		return
	}
	//TODO 校验手机号码格式
	// loginByPassword.Mobile不匹配正则表达式，就报错
	conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Logger.Error("LoginByPassword 拨号出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	client := pb.NewAccountServiceClient(conn)
	r, err := client.GetAccountByMobile(context.Background(), &pb.MobileRequest{Mobile: loginByPassword.Mobile})
	if err != nil {
		log.Logger.Error("GRPC GetAccountByMobile 出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	cheRes, err := client.CheckPassword(context.Background(), &pb.CheckPasswordRequest{
		Password:       loginByPassword.Password,
		HashedPassword: r.Password,
		AccountId:      uint32(r.Id),
	})
	if err != nil {
		log.Logger.Error("GRPC CheckPassword 出错：" + err.Error())
		e := HandleError(err)
		c.JSON(http.StatusOK, gin.H{
			"msg": e,
		})
		return
	}
	checkResult := "登录失败"
	if cheRes.Result {
		checkResult = "登录成功"
		j := jwt_op.NewJWT()
		now := time.Now()
		claims := jwt_op.CustomClaims{
			StandardClaims: jwt.StandardClaims{
				NotBefore: now.Unix(),
				ExpiresAt: now.Add(time.Hour * 24 * 30).Unix(),
			},
			ID:          r.Id,
			NickName:    r.Nickname,
			AuthorityId: int32(r.Role),
		}
		token, err := j.GenerateJWT(claims)
		if err != nil {
			log.Logger.Error("GRPC GenerateJWT 出错：" + err.Error())
			e := HandleError(err)
			c.JSON(http.StatusOK, gin.H{
				"msg": e,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":    "",
			"result": checkResult,
			"token":  token,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":    "",
		"result": checkResult,
		"token":  "",
	})
}

// 32ae35dee2ab4601b44008a7b4153823d7da3da3edadd596b1aac176fc26767a
