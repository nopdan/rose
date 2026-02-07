package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/nopdan/rose/converter"
	"github.com/nopdan/rose/encoder"
	"github.com/nopdan/rose/filter"
	"github.com/nopdan/rose/format"
	"github.com/nopdan/rose/frontend"
	"github.com/nopdan/rose/model"
)

// FormatInfo 格式信息
type FormatInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      int    `json:"kind"`
	Ext       string `json:"ext"`
	CanImport bool   `json:"canImport"`
	CanExport bool   `json:"canExport"`
}

// StoredFile 存储的文件信息
type StoredFile struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

// ConvertRequest 转换请求
type ConvertRequest struct {
	FileID       string                        `json:"fileId"`
	InputFormat  string                        `json:"inputFormat"`
	OutputFormat string                        `json:"outputFormat"`
	InputCustom  *converter.CustomFormatConfig `json:"inputCustom,omitempty"`
	OutputCustom *converter.CustomFormatConfig `json:"outputCustom,omitempty"`
	Encoder      *EncoderConfig                `json:"encoder,omitempty"`
	Filter       *FilterConfig                 `json:"filter,omitempty"`
}

// EncoderConfig 编码器配置
type EncoderConfig struct {
	Type            string `json:"type"`            // "pinyin" or "wubi"
	Schema          string `json:"schema"`          // "86", "98", "06"
	CodeTableFileID string `json:"codeTableFileId"` // 自定义码表文件ID
	UseAABC         bool   `json:"useAABC"`         // 组词规则
}

// FilterConfig 过滤器配置
type FilterConfig struct {
	MinLength     int      `json:"minLength"`
	MaxLength     int      `json:"maxLength"`
	MinFrequency  int      `json:"minFrequency"`
	MaxFrequency  int      `json:"maxFrequency"`
	FilterEnglish bool     `json:"filterEnglish"`
	FilterNumber  bool     `json:"filterNumber"`
	CustomRules   []string `json:"customRules"`
}

var (
	dist        = frontend.Dist
	storedFiles = make(map[string]*StoredFile)
	fileCounter = 0
)

func Serve(port int) {
	distFS, _ := fs.Sub(dist, "dist")
	http.Handle("/", http.FileServer(http.FS(distFS)))

	http.HandleFunc("/api/formats", handleFormats)
	http.HandleFunc("/api/upload", handleUpload)
	http.HandleFunc("/api/convert", handleConvert)

	ln, actualPort, err := listenWithFallback(port, 20)
	if err != nil {
		log.Fatalf("Failed to bind port: %v", err)
	}
	addr := fmt.Sprintf("http://localhost:%d", actualPort)
	log.Println("Listening on " + addr)
	openBrowser(addr)
	if err := http.Serve(ln, nil); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

func listenWithFallback(port, maxAttempts int) (net.Listener, int, error) {
	for i := 0; i <= maxAttempts; i++ {
		curr := port + i
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", curr))
		if err == nil {
			return ln, curr, nil
		}
		if !isAddrInUse(err) {
			return nil, 0, err
		}
	}
	return nil, 0, fmt.Errorf("no available port from %d to %d", port, port+maxAttempts)
}

func isAddrInUse(err error) bool {
	if errors.Is(err, syscall.EADDRINUSE) {
		return true
	}
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if errors.Is(opErr.Err, syscall.EADDRINUSE) {
			return true
		}
		if sysErr, ok := opErr.Err.(*os.SyscallError); ok {
			return errors.Is(sysErr.Err, syscall.EADDRINUSE)
		}
	}
	return false
}

// handleFormats 获取所有支持的格式列表
func handleFormats(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}
	log.Println("GET /api/formats")
	writeJSON(w, getFormatList())
}

// handleUpload 处理文件上传
func handleUpload(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(1 << 32)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("POST /api/upload err: %v\n", err)
		http.Error(w, "文件上传失败", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("POST /api/upload %v", handler.Filename)
	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "文件读取失败", http.StatusInternalServerError)
		return
	}
	stored, err := saveUploadedFile(handler.Filename, fileData)
	if err != nil {
		http.Error(w, "文件保存失败", http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"id":       stored.ID,
		"filename": stored.Filename,
		"size":     stored.Size,
	})
}

// handleConvert 执行词库转换，将结果保存为临时文件并返回下载 ID
func handleConvert(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ConvertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "请求格式错误", http.StatusBadRequest)
		return
	}
	inputFormat := strings.TrimSpace(req.InputFormat)
	outputFormat := strings.TrimSpace(req.OutputFormat)
	if inputFormat == "" || outputFormat == "" {
		http.Error(w, "输入或输出格式缺失", http.StatusBadRequest)
		return
	}
	stored := findStoredFile(req.FileID)
	if stored == nil {
		http.Error(w, "文件不存在，请先上传", http.StatusBadRequest)
		return
	}
	log.Printf("POST /api/convert file=%s input=%s output=%s", stored.Filename, inputFormat, outputFormat)

	conv := converter.NewConverter()

	// 构建 encoder
	if req.Encoder != nil {
		if enc := buildEncoder(req.Encoder); enc != nil {
			conv.RegisterEncoder(enc)
		}
	}

	// 构建 filter
	applyFilters(conv, req.Filter)

	// 生成输出路径
	outputFilename := buildOutputFilename(stored.Filename, outputFormat)
	outputDir := "output"
	_ = os.MkdirAll(outputDir, 0o755)
	outputPath := uniquePath(filepath.Join(outputDir, outputFilename))

	result, err := conv.Convert(&converter.Job{
		Input: &converter.InputSpec{
			Source: model.NewFileSource(stored.Path),
			Path:   stored.Filename,
			Format: inputFormat,
			Custom: req.InputCustom,
		},
		Output: &converter.OutputSpec{
			Path:      outputPath,
			Format:    outputFormat,
			Custom:    req.OutputCustom,
			Overwrite: true,
		},
	})
	if err != nil {
		log.Printf("转换失败: %v", err)
		http.Error(w, fmt.Sprintf("转换失败: %v", err), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{
		"outputPath": outputPath,
		"stats": map[string]int{
			"inputEntries":  result.Stats.InputEntries,
			"outputEntries": result.Stats.OutputEntries,
			"filteredOut":   result.Stats.FilteredOut,
		},
	})
}

// --- 辅助函数 ---

func getFormatList() []FormatInfo {
	formats := format.GlobalRegistry.List()
	items := make([]FormatInfo, 0, len(formats))
	for _, f := range formats {
		canImport, canExport := false, false
		if actual, ok := format.GlobalRegistry.Get(f.ID); ok {
			if _, ok := actual.(model.Importer); ok {
				canImport = true
			}
			if _, ok := actual.(model.Exporter); ok {
				canExport = true
			}
		}
		items = append(items, FormatInfo{
			ID:        f.ID,
			Name:      f.Name,
			Kind:      int(f.Type),
			Ext:       f.Extension,
			CanImport: canImport,
			CanExport: canExport,
		})
	}
	return items
}

// uniquePath 在文件名重复时添加 (2)、(3)… 后缀
func uniquePath(p string) string {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return p
	}
	ext := filepath.Ext(p)
	base := strings.TrimSuffix(p, ext)
	for i := 2; ; i++ {
		candidate := fmt.Sprintf("%s(%d)%s", base, i, ext)
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}

func buildEncoder(cfg *EncoderConfig) encoder.Encoder {
	if cfg == nil {
		return nil
	}
	switch cfg.Type {
	case "pinyin":
		return encoder.NewEncoder("pinyin", nil)
	case "wubi":
		params := map[string]any{
			"schema":  cfg.Schema,
			"useAABC": cfg.UseAABC,
		}
		if cfg.CodeTableFileID != "" {
			if f := findStoredFile(cfg.CodeTableFileID); f != nil {
				data, err := os.ReadFile(f.Path)
				if err == nil {
					params["codeTableData"] = data
					params["schema"] = "custom"
				}
			}
		}
		return encoder.NewEncoder("wubi", params)
	case "none":
		return &encoder.NullEncoder{}
	default:
		return nil
	}
}

func buildOutputFilename(inputFilename, outputFormat string) string {
	base := strings.TrimSuffix(inputFilename, filepath.Ext(inputFilename))
	name := base + "_" + outputFormat
	ext := getFormatExtension(outputFormat)
	if ext == "" {
		ext = ".txt"
	}
	name += ext
	return name
}

func getFormatExtension(formatID string) string {
	if f, ok := format.GlobalRegistry.Get(formatID); ok {
		if info := f.Info(); info != nil {
			return info.Extension
		}
	}
	return ""
}

func saveUploadedFile(filename string, data []byte) (*StoredFile, error) {
	uploadDir := filepath.Join(os.TempDir(), "rose_uploads")
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		return nil, err
	}
	fileCounter++
	fileID := fmt.Sprintf("f_%d", fileCounter)
	filePath := filepath.Join(uploadDir, fmt.Sprintf("%s_%s", fileID, filename))
	if err := os.WriteFile(filePath, data, 0o644); err != nil {
		return nil, err
	}
	stored := &StoredFile{
		ID: fileID, Filename: filename, Path: filePath, Size: int64(len(data)),
	}
	storedFiles[fileID] = stored
	return stored, nil
}

func findStoredFile(fileID string) *StoredFile {
	if fileID == "" {
		return nil
	}
	if f, ok := storedFiles[fileID]; ok {
		return f
	}
	return nil
}

func applyFilters(conv *converter.Converter, cfg *FilterConfig) {
	if cfg == nil {
		return
	}
	if cfg.MinLength > 0 || cfg.MaxLength > 0 {
		conv.AddFilter(filter.NewLengthFilter(cfg.MinLength, cfg.MaxLength))
	}
	if cfg.MinFrequency > 0 || cfg.MaxFrequency > 0 {
		conv.AddFilter(filter.NewFrequencyFilter(cfg.MinFrequency, cfg.MaxFrequency))
	}
	if cfg.FilterEnglish || cfg.FilterNumber {
		conv.AddFilter(filter.NewCharacterFilter(cfg.FilterEnglish, cfg.FilterNumber))
	}
	if len(cfg.CustomRules) > 0 {
		conv.AddFilter(filter.NewRegexFilter(cfg.CustomRules))
	}
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")
}

func openBrowser(url string) {
	var name string
	switch runtime.GOOS {
	case "windows":
		name = "explorer"
	case "linux":
		name = "xdg-open"
	default:
		name = "open"
	}
	cmd := exec.Command(name, url)
	cmd.Start()
}
