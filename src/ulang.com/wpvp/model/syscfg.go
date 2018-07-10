package model

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func SysCfgIncUserId() (string, error) {
	session := sess.Copy()
	defer session.Close()
	cfg := &SysCfg{}
	c := session.DB(db_name).C(DOC_SYSCFG)
	err := c.Find(nil).One(cfg)
	if err != nil {
		return "", err
	}
	user_id := cfg.NextUserId
	err = c.UpdateId(cfg.Id, bson.M{"$inc": bson.M{"NextUserId": 1}})
	return fmt.Sprintf("U%d", user_id), err
}
func SysCfgIncRetrievalTask() (string, error) {
	session := sess.Copy()
	defer session.Close()
	cfg := &SysCfg{}
	c := session.DB(db_name).C(DOC_SYSCFG)
	c.Find(nil).One(cfg)
	retrievalTaskId := cfg.NextTaskId
	err := c.UpdateId(cfg.Id, bson.M{"$inc": bson.M{"NextTaskId": 1}})
	return fmt.Sprintf("RT%d", retrievalTaskId), err
}

func NewSysCfg() (*SysCfg, error) {

	session := sess.Copy()
	defer session.Close()
	var sysCfg []SysCfg
	c := session.DB(db_name).C(DOC_SYSCFG)
	err := c.Find(nil).All(&sysCfg)
	if err != nil {
		return nil, err
	}
	if len(sysCfg) == 1 {
		return &sysCfg[0], nil
	}
	sc := SysCfg{}
	sc.Id = bson.NewObjectId()
	sc.NextUserId = 1000001
	sc.NextTaskId = 1000001
	err = c.Insert(&sc)
	return &sc, err

}
