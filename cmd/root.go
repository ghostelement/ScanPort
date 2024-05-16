package cmd

import (
	"ScanPort/svc/scan"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var defConfigName = "host.yml"

func init() {
	// 接受参数
	rootCmd.PersistentFlags().String("version", "", "版本")
}

// 根命令
var rootCmd = &cobra.Command{
	Use:     "scanport",
	Short:   "Scan host port with yaml configuration.",
	Example: "scanport host.yml",
	Version: "0.0.1",
	Args:    cobra.MaximumNArgs(1), // 最多接收一个参数
	RunE: func(cmd *cobra.Command, args []string) error {
		var profile string
		// 没有参数传递时，使用默认配置文件
		if len(args) == 0 {
			//检查当前目录是否存在配置文件 host.yml
			_, err := os.Stat(defConfigName)
			if err != nil {
				if os.IsNotExist(err) {
					return fmt.Errorf(color.RedString("host.yml does not exist"))
				}
			}
			profile = defConfigName
		} else {
			profile = args[0]
		}

		//读取配置文件
		config, err := scan.Config(profile)
		if err != nil {
			return err
		}

		//验证job任务字段是否合规并执行任务
		err = config.Validate()
		if err != nil {
			fmt.Println(color.RedString("Error: "), err)
		} else {
			//执行任务
			config.Scan()
		}

		return nil
	},
}

// 添加一个名为 "version" 的标志（对应 `-v` 和 `--version`）
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of adp",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("adp version: %s\n", rootCmd.Version)
	},
}

// Execute 将所有子命令添加到root命令并适当设置标志。
// 这由 main.main() 调用。它只需要对 rootCmd 调用一次。
func Execute() {
	// 将 "version" 子命令绑定到 rootCmd 上，使其可以通过 `-v` 或 `--version` 调用
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
	rootCmd.SetVersionTemplate("")
	rootCmd.AddCommand(versionCmd)
	// 执行根命令，并检查执行过程中是否发生了错误。
	cobra.CheckErr(rootCmd.Execute())
}
