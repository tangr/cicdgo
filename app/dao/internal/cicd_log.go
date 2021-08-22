// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// CicdLogDao is the manager for logic model data accessing and custom defined data operations functions management.
type CicdLogDao struct {
	gmvc.M                // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	C      cicdLogColumns // C is the short type for Columns, which contains all the column names of Table for convenient usage.
	DB     gdb.DB         // DB is the raw underlying database management object.
	Table  string         // Table is the underlying table name of the DAO.
}

// CicdLogColumns defines and stores column names for table cicd_log.
type cicdLogColumns struct {
	Id         string //
	JobId      string //
	TaskStatus string //
	Ipaddr     string //
	UpdatedAt  string //
	Output     string //
}

// NewCicdLogDao creates and returns a new DAO object for table data access.
func NewCicdLogDao() *CicdLogDao {
	columns := cicdLogColumns{
		Id:         "id",
		JobId:      "job_id",
		TaskStatus: "task_status",
		Ipaddr:     "ipaddr",
		UpdatedAt:  "updated_at",
		Output:     "output",
	}
	return &CicdLogDao{
		C:     columns,
		M:     g.DB("default").Model("cicd_log").Safe(),
		DB:    g.DB("default"),
		Table: "cicd_log",
	}
}
