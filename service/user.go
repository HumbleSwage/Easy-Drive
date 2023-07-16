package service

import (
	"context"
	"easy-drive/conf"
	"easy-drive/consts"
	"easy-drive/pkg/ctl"
	"easy-drive/pkg/e"
	"easy-drive/pkg/utils"
	"easy-drive/pkg/utils/commonUtil"
	"easy-drive/pkg/utils/fileUtil"
	"easy-drive/repositry/cache"
	"easy-drive/repositry/dao"
	"easy-drive/repositry/model"
	"easy-drive/types"
	"errors"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var UserSrv *UserService
var UserSrvOnce sync.Once

type UserService struct {
}

func GetUserSrv() *UserService {
	UserSrvOnce.Do(func() {
		UserSrv = &UserService{}
	})
	return UserSrv
}

func (us *UserService) UserRegister(ctx context.Context, req *types.UserRegisterReq) (resp interface{}, err error) {
	code := e.Success

	userDao := dao.NewUserDao(ctx)
	if userDao.IsUserNameExists(req.NickName) {
		code = e.UserAlreadyExistsError
		return ctl.RespError(code), nil
	}

	// 缓存key
	key := cache.VerificationCodeCacheKey(0, req.Email)

	// 校验邮箱
	client := cache.RedisClient
	val, err := client.Get(key).Result()
	if err != nil {
		utils.LogrusObj.Error("从redis获取校验邮箱验证码时发生错误:", err)
		return ctl.RespError(), nil
	}
	if !strings.EqualFold(val, req.EmailCode) {
		code = e.EmailCodeError
		return ctl.RespError(code), nil
	}

	// 查询系统设置
	systemSetting := dao.NewSystemDaoByDB(userDao.DB)
	setting, err := systemSetting.GetSystemSettingById(consts.UserInitSpaceId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("获取用户默认内存出错：", err)
			return ctl.RespError(), nil
		}
	}

	// 邮箱校验通过,注册用户到数据库
	userId := commonUtil.GenerateUserID() // 生成用户id
	user := &model.User{
		UserId:     userId,
		UserName:   req.NickName,
		NickName:   req.NickName,
		Email:      req.Email,
		Status:     1,
		TotalSpace: cast.ToInt64(setting.Text) * 1024 * 1024, // 内存单位为G
	}

	// 密码加密
	err = user.SetPassword(req.Password)
	if err != nil {
		utils.LogrusObj.Error("用户密码加密失败:", err)
		return ctl.RespError(), nil
	}
	// 用户数据入库
	err = userDao.AddUser(user)
	if err != nil {
		utils.LogrusObj.Error("创建用户时发生错误:", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(), err
}

func (us *UserService) UserLogin(ctx *gin.Context, req *types.UserLoginReq) (resp interface{}, err error) {
	code := e.Success
	fmt.Println(0)
	// 检验用户是否存在
	userDao := dao.NewUserDao(ctx)
	flag, err := userDao.IsEmailExists(req.Email)
	if err != nil {
		utils.LogrusObj.Error("检查邮箱是否存在出错:", err)
		return ctl.RespError(), nil
	}
	if !flag {
		// 用户不存在
		code = e.UserNotRegisterError
		return ctl.RespError(code), nil
	}
	fmt.Println(1)

	// 获取用户
	user, err := userDao.GetUserByEmail(req.Email)
	if err != nil {
		utils.LogrusObj.Error("通过邮箱获取用户错误:", err)
		return ctl.RespError(), nil
	}

	// 检验账户状态
	if user.Status == 0 {
		code = e.UseAccountDisable
		return ctl.RespError(code), nil
	}

	// 校验密码
	if req.Password != user.Password {
		code = e.UserPasswordError
		return ctl.RespError(code), nil
	}

	// 用户token:TODO:后面全部改成Token
	token, err := utils.GenerateToken(user.UserId, user.UserName, user.Authority)
	if err != nil {
		utils.LogrusObj.Error("生成用户token错误:", err)
		return ctl.RespError(), nil
	}

	// 获取session对象
	session := sessions.Default(ctx)

	// 存入用户的id
	session.Set(consts.UserInfo, user.UserId)

	// 更新用户存储
	fileDao := dao.NewFileDaoByDB(userDao.DB)
	useSpace, err := fileDao.SelectUserFileSpace(user, consts.NormalFile.Index())
	if err != nil {
		utils.LogrusObj.Error("根据用户查询已使用内存时出错:", err)
		return ctl.RespError(), nil
	}
	if user.UseSpace < useSpace {
		user.UseSpace = useSpace
		// 用户入库
		err = userDao.UpdateUser(user)
		if err != nil {
			utils.LogrusObj.Error("更新用户已使用内存时出错:", err)
			return ctl.RespError(), nil
		}
	}

	// 保存
	err = session.Save()
	if err != nil {
		utils.LogrusObj.Error("保存session出错:", err)
		code := e.StoreInSessionError
		return ctl.RespError(code), nil
	}

	// 获取当前工作目录
	workDir, _ := os.Getwd()

	// 设置配置文件路径
	configFile := filepath.Join(workDir, "conf/local/whiteList.yaml")
	viper.SetConfigFile(configFile)

	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Failed to read configuration file:", err)
		return
	}

	// 从配置文件中获取白名单数据
	whitelist := viper.GetStringSlice("whiteList")

	// 检查该用户email是否在白名单中
	admin := false
	for _, mailAddress := range whitelist {
		if mailAddress == user.Email {
			admin = true
			break
		}
	}
	// 绑定返回数据
	data := types.UserLoginResp{
		UserId:    user.UserId,
		NickName:  user.NickName,
		Authority: user.Authority,
		Token:     token,
		Admin:     admin,
	}

	return ctl.RespSuccessWithData(data), nil
}

func (us *UserService) UserResetPwd(ctx context.Context, req *types.UserResetPwdReq) (resp interface{}, err error) {
	code := e.Success

	// 检查是已否注册
	userDao := dao.NewUserDao(ctx)
	flag, err := userDao.IsEmailExists(req.Email)
	if err != nil {
		utils.LogrusObj.Error("检查邮箱是否存在出错:", err)
		return ctl.RespError(), nil
	}
	if !flag {
		// 用户不存在
		code = e.UserNotRegisterError
		return ctl.RespError(code), nil
	}

	// 获取缓存key
	key := cache.VerificationCodeCacheKey(1, req.Email)

	// 获取缓存value
	client := cache.RedisClient
	value, err := client.Get(key).Result()
	if err != nil {
		utils.LogrusObj.Error("从redis中获取重置密码邮箱验证码出错:", err)
		return ctl.RespError(), nil
	}

	// 校验验证码
	if !strings.EqualFold(value, req.EmailCode) {
		// 校验未通过
		code := e.EmailCodeError
		return ctl.RespError(code), nil
	}

	// 修改用户密码
	user, err := userDao.GetUserByEmail(req.Email)
	if err != nil {
		utils.LogrusObj.Error("根据邮箱获取用户出错:", err)
		return ctl.RespError(), nil
	}
	err = user.SetPassword(req.Password)
	if err != nil {
		utils.LogrusObj.Error("修改用户密码出错:", err)
		return ctl.RespError(), nil
	}
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更新用户密码入库出错:", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(), nil
}

func (us *UserService) GetUserAvatar(ctx context.Context, userId string) (resp interface{}, err error) {
	// TODO:这个地方先暂时这么写，后面不要用io流的形式，而是直接用网址的形式
	var content []byte

	// 首先获取用户头像地址
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.LogrusObj.Error("根据用户id获取用户出错:", err)
			utils.LogrusObj.Error("用户id:", userId)
			return ctl.RespError(), nil
		}
	}

	// 图片位置
	pConfig := conf.Conf.UploadPath

	// 头像
	var avatarPath string
	if user.Avatar == "" {
		// 读取默认头像
		avatarPath = pConfig.AvatarPath + "default_avatar.png"
	} else {
		// 如果头像不为空，则返回数据库地址中对应位置的头像
		avatarPath = pConfig.AvatarPath + user.Avatar
	}

	// 获取文件路径
	workDir, _ := os.Getwd()
	avatarPath = filepath.Join(workDir, avatarPath)
	fileInfo, err := os.Stat(avatarPath)
	if os.IsNotExist(err) {
		utils.LogrusObj.Error("头像不存在:", err)
		return ctl.RespError(), nil
	}

	// 打开文件
	file, err := os.OpenFile(avatarPath, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		utils.LogrusObj.Error("打开头像文件出错:", err)
		return ctl.RespError(), nil
	}

	// 创建缓冲流，并读入
	content = make([]byte, fileInfo.Size())
	if _, err := file.Read(content); err != nil {
		utils.LogrusObj.Error("文件对象写出错:", err)
		return ctl.RespError(), nil
	}

	// 返回文件字节流
	return content, err
}

func (us *UserService) GetUseSpace(ctx *gin.Context) (resp interface{}, err error) {
	code := e.Success

	// session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 获取用户空间
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error("根据用户id获取用户出错：", err)
		return ctl.RespError(), nil
	}

	// 定义返回数据
	data := &types.UserSpaceResp{
		UseSpace:   user.UseSpace,
		TotalSpace: user.TotalSpace,
	}
	return ctl.RespSuccessWithData(data, code), nil
}

func (us *UserService) UserLogout(ctx *gin.Context) (resp interface{}, err error) {
	code := e.Success

	// session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 记录用户登录时间
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error("通过用户id获取用户出错:", err)
		return ctl.RespError(), nil
	}
	currentTime := time.Now()
	user.LastLoginTime = &currentTime
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更新用户最后一次登录时间出错:", err)
		return ctl.RespError(), nil
	}

	// 删除session
	if u := session.Get(consts.UserInfo); u != nil {
		session.Delete(consts.UserInfo)
		err = session.Save()
		if err != nil {
			utils.LogrusObj.Error("用户退出删除session出错:", err)
			return nil, err
		}
	} else {
		code := e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 响应成功
	return ctl.RespSuccess(), err
}

func (us *UserService) UpdateUserAvatar(ctx *gin.Context, req *types.UpdateUserAvatarReq) (resp interface{}, err error) {
	code := e.Success

	// 获取session
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 获取用户
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error(err)
		return ctl.RespError(), nil
	}

	// 保存图片到本地
	path, err := fileUtil.UploadAvatarToLocalStatic(req.File, req.FileSize, userId)
	if err != nil {
		code = e.UpdateAvatarError
		utils.LogrusObj.Error("头像上传到本地失败:", err)
		return ctl.RespError(code), nil
	}

	// 更新用户头像：注意这个地方最后跟更新用户有关系
	user.Avatar = path
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更换头像时更新用户出错:", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(code), err
}

func (us *UserService) UpdatePassword(ctx *gin.Context, req *types.UpdatePasswordReq) (resp interface{}, err error) {
	code := e.Success

	// 获取session对象
	session := sessions.Default(ctx)

	// 获取用户id
	var userId string
	if u := session.Get(consts.UserInfo); u != nil {
		userId = u.(string)
	} else {
		code = e.UserSessionExpiration
		return ctl.RespError(code), nil
	}

	// 获取用户
	userDao := dao.NewUserDao(ctx.Request.Context())
	user, err := userDao.GetUserByUserId(userId)
	if err != nil {
		utils.LogrusObj.Error("根据用户id获取用户出错：", err)
		return ctl.RespError(), nil
	}

	// 更新用户密码
	err = user.SetPassword(req.Password)
	if err != nil {
		utils.LogrusObj.Error("用户密码加密失败:", err)
		return ctl.RespError(), nil
	}

	// 用户入库
	err = userDao.UpdateUser(user)
	if err != nil {
		utils.LogrusObj.Error("更新用户密码出错:", err)
		return ctl.RespError(), nil
	}

	return ctl.RespSuccess(code), err
}
