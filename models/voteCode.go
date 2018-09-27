package models

type VoteCode string
type VoteCodeInfo struct {
	Filter
	ExpireDate string
}

func Decode(string) (*VoteCodeInfo, error) {
	ret := &VoteCodeInfo{}
	return ret, nil
}
