package models

import "github.com/pkg/errors"

type VoteCode string
type VoteCodeInfo struct {
	Filter
	ExpireDate string
}

// Verify 校验信息是否正确
func (vi *VoteCodeInfo) Verify() bool {
	if vi == nil {
		return false
	}
	return true
}

// 从投票码中解析出该票所在的班级
func Decode(voteCode string) (*VoteCodeInfo, error) {
	if voteCode != "aronic" {
		return nil, errors.New("invalid vote code")
	}
	ret := &VoteCodeInfo{}
	ret.Grade = 1
	ret.Index = 3
	return ret, nil
}
