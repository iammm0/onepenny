package attachment

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadAttachmentResponse 上传文件成功时的返回体
type UploadAttachmentResponse struct {
	Message string `json:"message"`
	FileURL string `json:"file_url"`
}

// ErrorResponse 通用错误返回体
type ErrorResponse struct {
	Error string `json:"error"`
}

// AttachmentController 处理附件相关请求
type AttachmentController struct{}

// NewAttachmentController 创建新的 AttachmentController
func NewAttachmentController() *AttachmentController {
	return &AttachmentController{}
}

// UploadAttachment godoc
// @Summary     上传单个文件
// @Description 接收 form-data 中的 file 字段，保存文件并返回可访问的 URL
// @Tags        attachment
// @Accept      multipart/form-data
// @Produce     json
// @Param       file formData  file                        true  "要上传的文件"
// @Success     200  {object}  UploadAttachmentResponse     "上传成功"
// @Failure     400  {object}  ErrorResponse                "请求中没有文件或参数错误"
// @Failure     500  {object}  ErrorResponse                "服务器内部错误，文件保存失败"
// @Router      /attachment [post]
func (ctl *AttachmentController) UploadAttachment(c *gin.Context) {
	// 从表单中获取文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "no file is received"})
		return
	}

	// 指定存储目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create upload directory"})
		return
	}

	// 使用 UUID 重命名文件
	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFilename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to save file"})
		return
	}

	// 构造访问 URL（假设已配置 r.Static("/uploads", "./uploads")）
	fileURL := fmt.Sprintf("/uploads/%s", newFilename)
	c.JSON(http.StatusOK, UploadAttachmentResponse{
		Message: "upload success",
		FileURL: fileURL,
	})
}
