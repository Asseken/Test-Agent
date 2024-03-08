package user

import (
	"fmt"
	"os/user"
)

type UserInfo struct {
	Username string `json:"username"`
	UID      string `json:"uid"`
	GID      string `json:"gid"`
	Name     string `json:"name"`
	HomeDir  string `json:"home_dir"`
}

func Userinfo() UserInfo {
	var userInfo UserInfo

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("无法获取当前用户信息:", err)
		return userInfo
	}

	userInfo = UserInfo{
		Username: currentUser.Username,
		UID:      currentUser.Uid,
		GID:      currentUser.Gid,
		Name:     currentUser.Name,
		HomeDir:  currentUser.HomeDir,
	}

	return userInfo
}
