/*
 * @Description:
 * @Author: mali
 * @Date: 2022-09-08 10:58:47
 * @LastEditTime: 2022-09-08 10:59:00
 * @LastEditors: VSCode
 * @Reference:
 */
package models

import "time"

// BaseModel 模型基类
type BaseModel struct {
	ID uint64 `gorm:"column:id;primaryKey;autoIncrement;" json:"id,omitempty"`
}

// CommonTimestampsField 时间戳
type CommonTimestampsField struct {
	CreatedAt time.Time `gorm:"column:created_at;index;" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;index;" json:"updated_at,omitempty"`
}
