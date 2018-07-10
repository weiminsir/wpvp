package model

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"crypto/sha256"
	"github.com/save/juju/errors"
)

const (
	DOC_SYSCFG         = "sysconfig"
	DOC_SYSUSER        = "sysuser"
	DOC_AUTHORITY      = "user_authority"
	DOC_RETRIEVAL_TASK = "task"
)
const (
	AUTHORITY_SUPER    = 1
	AUTHORITY_ADVANCED = 2
	AUTHORITY_COMMON   = 3
)

func ListSysUsers(cond *bson.M) ([]SysUser, error) {
	session := sess.Copy()
	defer session.Close()
	users := []SysUser{}
	c := session.DB(db_name).C(DOC_SYSUSER)
	err := c.Find(cond).All(&users)

	if err != nil {
		return nil, err
	}
	return users, err
}

func FindSysUser(user_name string) (*SysUser, bool) {
	session := sess.Copy()
	defer session.Close()
	user := SysUser{}
	c := session.DB(db_name).C(DOC_SYSUSER)
	cond := bson.M{}
	cond["Username"] = user_name
	err := c.Find(cond).One(&user)
	if err != nil {
		return nil, false
	}
	return &user, true
}

func NewSysUser(user_name, password string) (*SysUser, error) {
	session := sess.Copy()
	defer session.Close()
	user := SysUser{}
	user.Username = user_name
	user.UserId, _ = SysCfgIncUserId()
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	user.Password = hashedPassword
	c := session.DB(db_name).C(DOC_SYSUSER)
	err := c.Insert(&user)
	return &user, err
}

func (user *SysUser) SetPassword(password string) error {
	session := sess.Copy()
	defer session.Close()
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	user.Password = hashedPassword
	c := session.DB(db_name).C(DOC_SYSUSER)
	cond := bson.M{}
	cond["UserId"] = user.UserId
	return c.Update(cond, user)
}

func (user *SysUser) Validate(password string) bool {
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	return user.Password == hashedPassword
}

func (user *SysUser) Update() error {
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_SYSUSER)

	cond := bson.M{}
	cond["UserId"] = user.UserId
	return c.Update(cond, user)
}

func (user *SysUser) Delete(UserId string) error {
	if user.UserId == UserId {
		return errors.New("permission denied...can not delete self ")
	}
	sysauth, err := FindAuthority(user.UserId)
	if err != nil {
		return err
	}
	userauth, _ := FindAuthority(UserId)
	if err != nil {
		return err
	}

	if sysauth.Role > AUTHORITY_ADVANCED || sysauth.Role >= userauth.Role {
		return errors.New("permission denied... ")
	}
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_SYSUSER)
	cond := bson.M{}
	cond["UserId"] = UserId
	return c.Remove(cond)
}

// TODO 此处到时候更改为管理员设置权限
func NewAuthority(UserId string) (*Authority, error) {
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_AUTHORITY)
	auth := &Authority{}
	auth.UserId = UserId
	auth.Role = AUTHORITY_COMMON
	return auth, c.Insert(auth)
}
func FindAuthority(UserId string) (*Authority, error) {
	session := sess.Copy()
	defer session.Close()
	auth := Authority{}
	c := session.DB(db_name).C(DOC_AUTHORITY)
	cond := bson.M{}
	cond["UserId"] = UserId
	err := c.Find(cond).One(auth)
	if err != nil {
		return nil, err
	}
	return &auth, nil
}
func (sysUser *SysUser) UpdateAuthority(UserId string) error {

	auth1, err := FindAuthority(sysUser.UserId)
	if err != nil {
		return errors.New("permission denied...")
	}
	if auth1.Role > AUTHORITY_ADVANCED {
		return errors.New("permission denied...")
	}
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_AUTHORITY)

	auth := &Authority{}
	cond := bson.M{}
	cond["UserId"] = UserId
	return c.Update(cond, auth)
}
