package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wrpota/go-echo/configs"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/repository/mysql/create"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the files and databases required for the project",
	Run: func(cmd *cobra.Command, args []string) {
		initDatabase(cmd, args)
	},
}

func initDatabase(cmd *cobra.Command, args []string) {
	//判断是否已安装
	database := configs.Get().GetString("Database.Write.DataBase")
	var count int64
	variable.GormReadDb.Raw("SELECT count(*) as count FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", database).Count(&count)
	if count > 0 {
		fmt.Fprintf(os.Stderr, "execute %s args:%v error: %v\n", cmd.Name(), args, "The database has been installed")
		return
	}

	//初始化数据库
	if sql, err := create.GetSql(); err == nil {
		for _, v := range sql {
			variable.GormWriteDb.Exec(v)
		}
	}

}

// 初始化数据库
// 管理员

// 用户
// 权限
// 角色
// 角色权限
// 用户角色
// 定时任务
