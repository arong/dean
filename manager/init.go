package manager

import (
	"os"

	"github.com/arong/dean/models"
	"github.com/astaxie/beego/logs"
)

// Init init da config
func Init(conf *models.DBConfig) {
	// allocate memory
	Ma.Init(conf)

	// data warm up
	err := Ma.LoadAllData()
	if err != nil {
		logs.Error("init failed", "err", err)
		os.Exit(-1)
	}

	Sm.store = &Ma
	Tm.store = &Ma
	Um.store = &Ma
}
