// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// CicdPipelineDao is the manager for logic model data accessing and custom defined data operations functions management.
type CicdPipelineDao struct {
	gmvc.M                     // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	C      cicdPipelineColumns // C is the short type for Columns, which contains all the column names of Table for convenient usage.
	DB     gdb.DB              // DB is the raw underlying database management object.
	Table  string              // Table is the underlying table name of the DAO.
}

// CicdPipelineColumns defines and stores column names for table cicd_pipeline.
type cicdPipelineColumns struct {
	Id           string //
	PipelineName string //
	GroupId      string //
	AgentId      string //
	Body         string //
	Author       string //
	UpdatedAt    string //
}

// NewCicdPipelineDao creates and returns a new DAO object for table data access.
func NewCicdPipelineDao() *CicdPipelineDao {
	columns := cicdPipelineColumns{
		Id:           "id",
		PipelineName: "pipeline_name",
		GroupId:      "group_id",
		AgentId:      "agent_id",
		Body:         "body",
		Author:       "author",
		UpdatedAt:    "updated_at",
	}
	return &CicdPipelineDao{
		C:     columns,
		M:     g.DB("default").Model("cicd_pipeline").Safe(),
		DB:    g.DB("default"),
		Table: "cicd_pipeline",
	}
}
