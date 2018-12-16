package base

const (
	// 最大分数
	MaxScore = 150

	// 最低分数
	MinScore = 0

	// Private Key
	Private = "private"
	Data    = "data"

	// time
	DateFormat     = "2006-01-02"
	DateTimeFormat = "2006-01-02 15:04:05"

	// enum
	// AccountTypeStudent => student
	AccountTypeStudent = 1
	// AccountTypeTeacher => teacher
	AccountTypeTeacher = 2

	// status code
	StatusValid    = 1 // Imply that this meta is available
	StatusArchived = 2 // Imply that this meta will be no longer in use, just exist for reference
	StatusAbandon  = 3 // Imply that this data is dropped by user
)
