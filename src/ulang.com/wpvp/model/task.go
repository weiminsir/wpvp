package model

import (
	"gopkg.in/mgo.v2/bson"
)

func NewRetrieveTask(task *RetrievalTask) error {
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	err := c.Insert(&task)
	return err
}
func (task *RetrievalTask) Delete() error {
	session := sess.Copy()
	defer session.Close()
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	cond := bson.M{}
	cond["RetrievalTaskId"] = task.TaskId
	return c.Remove(cond)
}
func (task *RetrievalTask) Update() error {
	session := sess.Copy()

	defer session.Close()
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	cond := bson.M{}
	cond["TaskId"] = task.TaskId
	return c.Update(cond, task)
}
func FindRetrievalTask(taskId string) (*RetrievalTask, bool) {
	session := sess.Copy()

	defer session.Close()
	task := RetrievalTask{}
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	cond := bson.M{}
	cond["TaskId"] = taskId
	err := c.Find(cond).One(&task)
	if err != nil {
		return nil, false
	}
	return &task, true
}
func GetRetrievalTask(cond *bson.M, sort *string) ([]RetrievalTask, error) {
	session := sess.Copy()
	defer session.Close()
	rts := []RetrievalTask{}
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	iter := c.Find(cond)
	if sort != nil {
		iter.Sort(*sort)
	}
	err := iter.All(&rts)
	if err != nil {
		return nil, err
	}
	return rts, err
}
func ListRetrievalTask(cond *bson.M, sort *string, pageStart, pageSize int) ([]RetrievalTask, int, error) {
	session := sess.Copy()
	defer session.Close()
	rts := []RetrievalTask{}
	c := session.DB(db_name).C(DOC_RETRIEVAL_TASK)
	iter := c.Find(cond)
	if sort != nil {
		iter.Sort(*sort)
	}
	cnt, _ := iter.Count()

	err := iter.Skip(pageStart * pageSize).Limit(pageSize).All(&rts)

	if err != nil {
		return nil, 0, err
	}
	return rts, cnt, err
}
