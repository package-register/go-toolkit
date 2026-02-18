package timetable

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fogleman/gg"
)

// ImgOption 包含生成图片所需的各种选项
type ImgOption struct {
	FontName        string // 字体名称
	FontPath        string // 字体路径
	StoragePath     string // 存储路径
	TemplateImgPath string // 模板图片路径
	Size            struct {
		Width  int
		Height int
	}
	Image *gg.Context // 图片对象
}

type ImgOptionModifier func(*ImgOption)

// WithFontOrPath 根据参数类型设置字体名称或字体路径
func WithFontOrPath(fontOrPath string) ImgOptionModifier {
	return func(img *ImgOption) {
		if _, err := os.Stat(fontOrPath); os.IsNotExist(err) {
			img.FontName = fontOrPath
		} else {
			img.FontPath = fontOrPath
		}
	}
}

// WithSize 返回一个设置图片尺寸的 ImgOptionModifier
func WithSize(w, h int) ImgOptionModifier {
	return func(img *ImgOption) {
		img.Size.Width = w
		img.Size.Height = h
	}
}

// WithStoragePath 返回一个设置存储路径的 ImgOptionModifier
func WithStoragePath(path string) ImgOptionModifier {
	return func(img *ImgOption) {
		img.StoragePath = path
	}
}

// NewGenerator 根据提供的修改函数创建一个新的 ImgOption 实例
func NewGenerator(modifiers ...ImgOptionModifier) *ImgOption {
	// 默认设置
	op := &ImgOption{
		FontName:        "Deng",
		StoragePath:     "output",
		TemplateImgPath: "init_photo.png",
		Size: struct {
			Width  int
			Height int
		}{
			Width:  2000,
			Height: 1200,
		},
	}

	op.Image = gg.NewContext(
		op.Size.Width,
		op.Size.Height,
	)

	for _, modifier := range modifiers {
		modifier(op)
	}

	op.InitFont()
	return op
}

func ExistFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}

// InitFont 方法用于初始化字体
func (o *ImgOption) InitFont() {
	// 如果字体路径不为空，则加载指定路径的字体
	if o.FontPath != "" {
		err := o.Image.LoadFontFace(o.FontPath, 30)
		if err != nil {
			fmt.Println("Error loading font:", err)
		}
		return
	}

	// 根据平台设置字体路径
	fontPath := func(platform string) (fontPath string) {
		switch platform {
		case "darwin":
			fontPath = "/Library/Fonts/" + o.FontName + ".ttf"
		case "linux":
			fontPath = "/etc/fonts/" + o.FontName + ".ttf"
		case "windows":
			fontPath = "C:\\Windows\\Fonts\\" + o.FontName + ".ttf"
		default:
			fontPath = ""
		}
		o.FontPath = fontPath
		return fontPath
	}(runtime.GOOS)

	if !ExistFile(fontPath) {
		fmt.Printf("字体文件未找到! %s", fontPath)
	}

	err := o.Image.LoadFontFace(fontPath, 30)
	if err != nil {
		fmt.Println("Error loading font:", err)
	}
}
