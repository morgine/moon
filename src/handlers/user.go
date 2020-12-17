package handlers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morgine/moon/pkg/cache"
	"github.com/morgine/moon/pkg/google_authenticator"
	"github.com/morgine/moon/src/errors"
	"github.com/morgine/moon/src/models"
	"github.com/morgine/moon/src/validators"
	"github.com/morgine/pkg/crypt/aes"
	"github.com/morgine/pkg/session"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"time"
)

var Now = time.Now

type User struct {
	m    *models.Model
	opts *Options
}

type Options struct {
	DB           *gorm.DB                    // 数据库 ORM
	CacheClient  cache.Client                // 数据缓存客户端
	Session      session.Storage             // token 存储器
	AuthExpires  int64                       // 会话过期时间
	AesCryptKey  []byte                      // 16 位字符串
	QRCodeConfig google_authenticator.Config // 谷歌验证器配置文件
}

func NewUser(opts *Options) (*User, error) {
	err := opts.DB.AutoMigrate(&User{})
	if err != nil {
		return nil, err
	}
	return &User{
		m: &models.Model{
			DB:  opts.DB,
			GAC: google_authenticator.NewClient(opts.QRCodeConfig),
			UserValidator: validators.NewUser(
				regexp.MustCompile("^[a-z0-9]{8,16}$"), // 用户名验证器
				regexp.MustCompile("^[\\w]{8,16}$"),    // 密码验证器
			),
			RecommendersCache: cache.NewRecommenders(cache.WithPrefixClient("recommenders_", opts.CacheClient)),
		},
		opts: opts,
	}, nil
}

// 注册账号，并绑定推荐人
func (h *User) Register() gin.HandlerFunc {
	type params struct {
		Username      string
		Password      string
		RecommenderID int // 推荐人 ID
	}
	return func(ctx *gin.Context) {
		ps := &params{}
		err := ctx.Bind(ps)
		if err != nil {
			SendError(ctx, err)
		} else {
			_, err = h.m.RegisterUser(ps.Username, ps.Password, ps.RecommenderID)
			if err != nil {
				SendError(ctx, err)
			} else {
				SendMessage(ctx, errors.StatusOK, "注册成功")
			}
		}
	}
}

func (h *User) getToken(ctx *gin.Context) string {
	return ctx.Request.Header.Get("Authorization")
}

// 用户鉴权
func (h *User) Auth(ctx *gin.Context) {
	token := h.getToken(ctx)
	if len(token) == 0 {
		return
	}
	user, err := h.decryptToken(token)
	if err != nil {
		SendError(ctx, err)
	} else {
		ok, err := h.opts.Session.CheckAndRefreshToken(user, token, h.opts.AuthExpires)
		if err != nil {
			SendError(ctx, err)
		} else {
			if !ok {
				SendError(ctx, errors.UserUnauthorized)
			} else {
				userID, err := strconv.Atoi(user)
				if err != nil {
					SendError(ctx, err)
				} else {
					ctx.Set("auth_user_id", userID)
				}
			}
		}
	}
}

// 获得登陆用户ID，需要在用户鉴权之后才有效
func (h *User) GetLoginUser(ctx *gin.Context) (userID int, ok bool) {
	user, ok := ctx.Get("auth_user_id")
	if ok {
		return user.(int), true
	} else {
		return 0, false
	}
}

// 获得登陆用户的账户信息
func (h *User) GetInfo(ctx *gin.Context) {
	userID, ok := h.GetLoginUser(ctx)
	if ok {
		user, err := h.m.GetUserByID(userID)
		if err != nil {
			SendError(ctx, err)
		} else {
			SendJSON(ctx, user)
		}
	}
}

// Login 登陆账号
func (h *User) Login() gin.HandlerFunc {
	type params struct {
		Username string
		Password string
	}
	return func(ctx *gin.Context) {
		ps := &params{}
		err := ctx.Bind(ps)
		if err != nil {
			SendError(ctx, err)
		} else {
			user, err := h.m.LoginUser(ps.Username, ps.Password)
			if err != nil {
				SendError(ctx, err)
			} else {
				uid := strconv.Itoa(user.ID)
				token, err := h.encryptToken(uid)
				if err != nil {
					SendError(ctx, err)
				} else {
					err = h.opts.Session.SaveToken(uid, token, h.opts.AuthExpires)
					if err != nil {
						SendError(ctx, err)
					} else {
						SendJSON(ctx, token)
					}
				}
			}
		}
	}
}

// 获得谷歌验证器二维码地址
func (h *User) GetGoogleAuthenticatorQRCodeUrl() gin.HandlerFunc {
	type params struct {
		With   int
		Height int
	}
	return func(ctx *gin.Context) {
		userID, ok := h.GetLoginUser(ctx)
		if ok {
			ps := &params{}
			err := ctx.Bind(ps)
			if err != nil {
				SendError(ctx, err)
			} else {
				imgUrl, err := h.m.GetGoogleAuthenticatorQRCodeUrl(userID, ps.With, ps.Height)
				if err != nil {
					SendError(ctx, err)
				} else {
					SendJSON(ctx, imgUrl)
				}
			}
		}
	}
}

// 绑定谷歌验证器
func (h *User) BindGoogle() gin.HandlerFunc {
	type params struct {
		GoogleCode string
	}
	return func(ctx *gin.Context) {
		userID, ok := h.GetLoginUser(ctx)
		if ok {
			ps := &params{}
			err := ctx.Bind(ps)
			if err != nil {
				SendError(ctx, err)
			} else {
				err := h.m.BindGoogleAuth(userID, ps.GoogleCode)
				if err != nil {
					SendError(ctx, err)
				} else {
					SendMessage(ctx, errors.StatusOK, "已绑定")
				}
			}
		}
	}
}

// ResetPassword 重置密码
func (h *User) ResetPassword() gin.HandlerFunc {
	type params struct {
		NewPassword string
		GoogleCode  string
	}
	return func(ctx *gin.Context) {
		userID, ok := h.GetLoginUser(ctx)
		if ok {
			ps := &params{}
			err := ctx.Bind(ps)
			if err != nil {
				SendError(ctx, err)
			} else {
				err := h.m.ResetPassword(userID, ps.GoogleCode, ps.NewPassword)
				if err != nil {
					SendError(ctx, err)
				} else {
					err = h.opts.Session.RemoveUser(strconv.Itoa(userID))
					if err != nil {
						SendError(ctx, err)
					} else {
						SendMessage(ctx, errors.StatusOK, "已重置")
					}
				}
			}
		}
	}
}

// Logout 退出登陆
func (h *User) Logout(adminID int, token string) error {
	return h.opts.Session.RemoveToken(strconv.Itoa(adminID), token)
}

// token 加密
func (h *User) encryptToken(adminID string) (token string, err error) {
	return aes.AesCBCEncrypt([]byte(fmt.Sprintf("%s:%10d", adminID, Now().UnixNano())), h.opts.AesCryptKey)
}

// token 解密
func (h *User) decryptToken(token string) (adminID string, err error) {
	data, err := aes.AesCBCDecrypt(token, h.opts.AesCryptKey)
	if err != nil {
		return "", err
	} else {
		sepIdx := bytes.Index(data, []byte(":"))
		return string(data[:sepIdx]), nil
	}
}
