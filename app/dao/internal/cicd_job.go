// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// CicdJobDao is the manager for logic model data accessing and custom defined data operations functions management.
type CicdJobDao struct {
	gmvc.M                // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	C      cicdJobColumns // C is the short type for Columns, which contains all the column names of Table for convenient usage.
	DB     gdb.DB         // DB is the raw underlying database management object.
	Table  string         // Table is the underlying table name of the DAO.
}

// CicdJobColumns defines and stores column names for table cicd_job.
type cicdJobColumns struct {
	Id          string //
	PipelineId  string //
	AgentId     string //
	Concurrency string //
	JobType     string //
	JobStatus   string //
	Script      string //
	Comment     string //
	Author      string //
	CreatedAt   string //
}

// NewCicdJobDao creates and returns a new DAO object for table data access.
func NewCicdJobDao() *CicdJobDao {
	columns := cicdJobColumns{
		Id:          "id",
		PipelineId:  "pipeline_id",
		AgentId:     "agent_id",
		Concurrency: "concurrency",
		JobType:     "job_type",
		JobStatus:   "job_status",
		Script:      "script",
		Comment:     "comment",
		Author:      "author",
		CreatedAt:   "created_at",
	}
	return &CicdJobDao{
		C:     columns,
		M:     g.DB("default").Model("cicd_job").Safe(),
		DB:    g.DB("default"),
		Table: "cicd_job",
	}
}
