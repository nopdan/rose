package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nopdan/rose/converter"
	"github.com/nopdan/rose/filter"
	"github.com/nopdan/rose/format"
	"github.com/nopdan/rose/model"
	"github.com/nopdan/rose/server"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "serve":
		serveWeb()
		return
	case "-l":
		listFormats()
		return
	case "-h":
		printUsage()
		return
	case "-v":
		printVersion()
		return
	}

	if len(os.Args) < 4 {
		fmt.Println("转换参数不足")
		printUsage()
		return
	}
	convertFile(os.Args[1:])
}

// printUsage 打印使用说明
func printUsage() {
	fmt.Println("蔷薇词库转换工具 - 重构版")
	fmt.Println()
	fmt.Println("用法:")
	fmt.Println("  rose serve [端口]                         # 启动 Web 服务 (默认端口 7800)")
	fmt.Println("  rose -l                                   # 列出支持的格式")
	fmt.Println("  rose <输入文件> <输入格式> <输出格式> [输出文件]  # 转换文件")
	fmt.Println()
	fmt.Println("示例:")
	fmt.Println("  rose serve          # 启动 Web 服务在 7800 端口")
	fmt.Println("  rose serve 8080     # 启动 Web 服务在 8080 端口")
	fmt.Println("  rose input.scel sogou_scel rime output.txt")
	fmt.Println()
	fmt.Println("其他命令:")
	fmt.Println("  rose -h  显示帮助信息")
	fmt.Println("  rose -v  显示版本信息")
}

// printVersion 打印版本信息
func printVersion() {
	fmt.Println("蔷薇词库转换工具 v2.0.0")
	fmt.Println("重构版 - 简洁架构设计")
	fmt.Println("作者: nopdan")
}

// listFormats 列出支持的格式
func listFormats() {
	formats := format.GlobalRegistry.List()

	fmt.Println("支持的格式:")
	fmt.Printf("%-20s %-15s %-10s %s\n", "ID", "类型", "可导出", "名称")
	fmt.Printf("%-20s %-15s %-10s %s\n", strings.Repeat("-", 20), strings.Repeat("-", 15), strings.Repeat("-", 10), strings.Repeat("-", 20))

	for _, format := range formats {
		canExport := "否"
		// 通过类型断言判断是否支持导出
		if actualFormat, ok := formatRegistryGet(format.ID); ok {
			if _, ok := actualFormat.(model.Exporter); ok {
				canExport = "是"
			}
		}
		fmt.Printf("%-20s %-15s %-10s %s\n", format.ID, format.Type.String(), canExport, format.Name)
	}
}

// convertFile 转换文件（新格式）
func convertFile(args []string) {
	inputFile := args[0]
	inputFormat := args[1]
	outputFormat := args[2]

	outputFile := ""
	if len(args) > 3 {
		outputFile = args[3]
	}

	outputFile = ensureOutputFilename(inputFile, outputFormat, outputFile)

	conv := converter.NewConverter()
	conv.AddFilter(filter.NewLengthFilter(1, 10))

	result, err := conv.Convert(&converter.Job{
		Input: &converter.InputSpec{
			Path:   inputFile,
			Format: inputFormat,
		},
		Output: &converter.OutputSpec{
			Path:      outputFile,
			Format:    outputFormat,
			Overwrite: true,
		},
		Encoder: nil,
	})
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
		return
	}

	fmt.Printf("转换完成，输出文件: %s\n", result.OutputFile)
	if result.Stats != nil {
		fmt.Printf("处理了 %d 个词条，输出 %d 个词条\n",
			result.Stats.InputEntries, result.Stats.OutputEntries)
	}
}

func ensureOutputFilename(inputFilename, outputFormat, outputFilename string) string {
	finalName := outputFilename
	if finalName == "" {
		dir := filepath.Dir(inputFilename)
		base := strings.TrimSuffix(filepath.Base(inputFilename), filepath.Ext(inputFilename))
		finalName = filepath.Join(dir, base+"_"+outputFormat)
	}

	if ext := getFormatExtension(outputFormat); ext != "" {
		if filepath.Ext(finalName) == "" {
			finalName += ext
		}
	}

	return finalName
}

func getFormatExtension(formatID string) string {
	if f, ok := format.GlobalRegistry.Get(formatID); ok {
		if info := f.Info(); info != nil {
			return info.Extension
		}
	}
	return ""
}

func formatRegistryGet(id string) (model.Format, bool) {
	return format.GlobalRegistry.Get(id)
}

// serveWeb 启动 Web 服务
func serveWeb() {
	port := 7800
	if len(os.Args) > 2 {
		if p, err := strconv.Atoi(os.Args[2]); err == nil {
			port = p
		}
	}
	server.Serve(port)
}
