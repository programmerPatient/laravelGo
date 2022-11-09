/*
 * @Description:根据数据库表生成结构体
 * @Author: mali
 * @Date: 2022-09-15 09:33:58
 * @LastEditTime: 2022-09-22 10:46:24
 * @LastEditors: VSCode
 * @Reference:
 */
package generate

import (
	"io"
	"os"
	"strings"
	"sync"
	"unicode"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/laravelGo/core/console"
	"github.com/laravelGo/core/database"
)

type Field struct {
	Field      string `gorm:"column:Field"`
	Type       string `gorm:"column:Type"`
	Null       string `gorm:"column:Null"`
	Key        string `gorm:"column:Key"`
	Default    string `gorm:"column:Default"`
	Extra      string `gorm:"column:Extra"`
	Privileges string `gorm:"column:Privileges"`
	Comment    string `gorm:"column:Comment"`
}

type Table struct {
	Name    string `gorm:"column:Name"`
	Comment string `gorm:"column:Comment"`
}

var wg sync.WaitGroup

/**
 * @Author: mali
 * @Func:
 * @Description: 执行生成
 * @Param:
 * @Return:
 * @Example:
 * @param {string} database 需要读取的库
 * @param {...string} tableNames 读取的指定表 不传的话读取所有表
 */
func Generate(database string, tableNames ...string) {
	tables := getTables(database, tableNames...) //生成所有表信息
	wg.Add(len(tables))
	for _, table := range tables {
		fields := getFields(table.Name)
		go generateModel(table, fields)
	}
	wg.Wait()
	console.Success("执行完成")
}

//获取表信息
func getTables(databse string, tableNames ...string) []Table {

	query := database.DB.Debug()
	var tables []Table
	if len(tableNames) > 0 {
		table_string := ""
		for _, v := range tableNames {
			if table_string == "" {
				table_string += "'" + v + "'"
			} else {
				table_string += ",'" + v + "'"
			}

		}
		query.Raw(
			"SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" +
				databse +
				"' AND table_name IN (" + table_string + ")" +
				";").Find(&tables)
	} else {
		query.Raw(
			"SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES WHERE table_schema='" +
				databse +
				"';").Find(&tables)
	}

	return tables
}

//获取所有字段信息
func getFields(tableName string) []Field {

	query := database.DB.Debug()
	var fields []Field
	query.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

//生成Model
func generateModel(table Table, fields []Field) {
	defer wg.Done()
	content := "package " + table.Name + "\n\n"
	//是否引入time包
	is_import_time := false

	//结构体字符串
	struct_str := "type " + generator.CamelCase(table.Name) + " struct {\n"
	//生成字段
	for _, field := range fields {
		fieldName := generator.CamelCase(field.Field)
		fieldJson := getFieldJson(field)
		fieldGorm := getFieldGorm(field)
		fieldType := getFiledType(field)
		if fieldType == "time.Time" {
			is_import_time = true
		}
		fieldComment := getFieldComment(field)
		struct_str += "	" + fieldName + " " + fieldType + " `" + fieldGorm + " " + fieldJson + "` " + fieldComment + "\n"
	}
	struct_str += "}\n"

	if is_import_time {
		content += "import \"time\"\n\n"
	}
	//表注释
	if len(table.Comment) > 0 {
		content += "// " + table.Comment + "\n"
	}
	content += struct_str
	content += "func (entity *" + generator.CamelCase(table.Name) + ") TableName() string {\n"
	content += "	" + `return "` + table.Name + `"`
	content += "\n}"
	dir := "app/models/" + table.Name + "/"

	// os.MkdirAll 会确保父目录和子目录都会创建，第二个参数是目录权限，使用 0777
	os.MkdirAll(dir, os.ModePerm)

	filename := dir + table.Name + ".go"
	var f *os.File
	var err error
	if checkFileIsExist(filename) {
		console.Error(filename + " 已存在，需删除才能重新生成...")
		_, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0666) //打开文件
		if err != nil {
			panic(err)
		}
	}

	f, err = os.Create(filename)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	_, err = io.WriteString(f, content)
	if err != nil {
		console.ExitIf(err)
	} else {
		console.Success(generator.CamelCase(table.Name) + " 已生成...")
	}
}

//获取字段类型
func getFiledType(field Field) string {
	if strings.Contains(field.Type, "integer") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	} else if strings.Contains(field.Type, "mediumint") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	} else if strings.Contains(field.Type, "bit") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	} else if strings.Contains(field.Type, "year") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	} else if strings.Contains(field.Type, "smallint") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint16"
		} else {
			return "int16"
		}
	} else if strings.Contains(field.Type, "tinyint") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint8"
		} else {
			return "int8"
		}
	} else if strings.Contains(field.Type, "bigint") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint64"
		} else {
			return "int64"
		}
	} else if strings.Contains(field.Type, "int") {
		if strings.Contains(field.Type, "unsigned") {
			return "uint32"
		} else {
			return "int32"
		}
	} else if strings.Contains(field.Type, "decimal") {
		return "float64"
	} else if strings.Contains(field.Type, "double") {
		return "float32"
	} else if strings.Contains(field.Type, "float") {
		return "float32"
	} else if strings.Contains(field.Type, "real") {
		return "float64"
	} else if strings.Contains(field.Type, "numeric") {
		return "float32"
	} else if strings.Contains(field.Type, "timestamp") {
		return "time.Time"
	} else if strings.Contains(field.Type, "datetime") {
		return "time.Time"
	} else if strings.Contains(field.Type, "time") {
		return "time.Time"
	} else if strings.Contains(field.Type, "date") {
		return "time.Time"
	} else if strings.Contains(field.Type, "bool") {
		return "bool"
	} else {
		return "string"
	}
}

//获取字段json描述
func getFieldJson(field Field) string {
	return `json:"` + Lcfirst(generator.CamelCase(field.Field)) + `,omitempty"`
}

// 首字母小写
func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

//获取字段gorm描述
func getFieldGorm(field Field) string {
	fieldContext := `gorm:"column:` + field.Field

	if field.Key == "PRI" {
		fieldContext = fieldContext + `;primaryKey`
	}
	if field.Key == "UNI" {
		fieldContext = fieldContext + `;unique`
	}
	if field.Extra == "auto_increment" {
		fieldContext = fieldContext + `;autoIncrement`
	}
	if field.Null == "NO" {
		fieldContext = fieldContext + `;not null`
	}
	return fieldContext + `"`
}

//获取字段说明
func getFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		//return "// " + field.Comment
		return "//" + strings.Replace(strings.Replace(field.Comment, "\r", "\\r", -1), "\n", "\\n", -1)
	}
	return ""
}

//检查文件是否存在
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// //生成JSON
// func generateJSON(table Table, fields []Field, ch chan int) {
// 	content := "package models\n\n"

// 	content += "type " + generator.CamelCase(table.Name) + "Reply struct {\n"
// 	//生成字段
// 	for _, field := range fields {
// 		fieldName := generator.CamelCase(field.Field)
// 		fieldJson := getFieldJson(field)
// 		fieldType := getFiledType(field)
// 		content += "	" + fieldName + " " + fieldType + " `" + fieldJson + "` " + "\n"
// 	}
// 	content += "}\n"

// 	filename := "app/models/" + table.Name + "_reply.go"
// 	if checkFileIsExist(filename) {
// 		fmt.Println(generator.CamelCase(table.Name) + " 已存在，需删除才能重新生成...")
// 		return
// 	}
// 	f, err := os.Create(filename)
// 	if err != nil {
// 		panic(err)
// 	}

// 	_, err = io.WriteString(f, content)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println(table.Name + "_reply.go 已生成...")
// 	}
// 	defer f.Close()

// 	ch <- 0
// 	close(ch)
// }
