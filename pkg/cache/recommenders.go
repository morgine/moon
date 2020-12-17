package cache

import "strconv"

type Recommenders struct {
	client Client
}

func NewRecommenders(client Client) *Recommenders {
	return &Recommenders{client: client}
}

// 获得推荐人缓存，recommenderID 小于 0 则表示缓存不存在，否则表示缓存存在
func (rs *Recommenders) Get(userID int) (recommenderID int, err error) {
	data, err := rs.client.Get(strconv.Itoa(userID))
	if err != nil {
		return 0, err
	}
	if len(data) > 0 {
		recommenderID, err = strconv.Atoi(string(data))
		if err != nil {
			return 0, err
		} else {
			return recommenderID, nil
		}
	} else {
		return -1, nil
	}
}

// 设置推荐人缓存
func (rs Recommenders) Set(userID, recommenderID int) error {
	return rs.client.Set(strconv.Itoa(userID), []byte(strconv.Itoa(recommenderID)), 0)
}
