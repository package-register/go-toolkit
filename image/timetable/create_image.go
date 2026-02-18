package timetable

import (
	"image/color"
	"math/rand"
	"regexp"
	"time"

	"github.com/fogleman/gg"
)

var (
	// colorList 颜色列表
	colorList = []color.RGBA{
		{173, 216, 230, 185}, // LightBlue (淡蓝色) - 略微鲜明
		{255, 228, 181, 185}, // LightSalmon (淡鲑鱼色) - 温暖柔和
		{144, 238, 144, 185}, // LightGreen (淡绿色) - 清新
		{255, 203, 164, 185}, // PeachPuff (桃色) - 温馨
		{221, 160, 221, 185}, // Plum (李子色) - 柔和紫色
		{255, 182, 193, 185}, // LightPink (淡粉色) - 甜美
		{176, 224, 230, 185}, // PowderBlue (粉蓝色) - 清爽
		{240, 230, 140, 185}, // Khaki (卡其色) - 自然
		{205, 133, 63, 185},  // Goldenrod (金菊黄) - 温和
		{255, 160, 122, 185}, // Salmon (鲑鱼色) - 比淡鲑鱼色鲜明
		{169, 169, 169, 185}, // DarkGray (深灰色) - 中性，用于区分
		{210, 180, 140, 185}, // Tan (棕褐色) - 大地色
		{173, 216, 230, 185}, // LightSkyBlue (淡天蓝色) - 比淡蓝色更蓝一点
		{255, 165, 0, 185},   // Orange (橙色) - 适度鲜明，注意使用
		{152, 251, 152, 185}, // LightSeaGreen (浅海绿) - 清新略深
	}
	weekPlaceData = []float64{250, 500, 750, 1000, 1250, 1500, 1750}
	hPlaceData    = map[string]float64{
		"第1节": 125, "第2节": 245,
		"第3节": 365, "第4节": 490,
		"第5节": 625, "第6节": 740,
		"第7节": 855, "第8节": 975,
	}
)

// IsExistAndColor 方法用于检查课程是否存在于颜色数据中
func (o *ImgOption) IsExistAndColor(colorData map[string]color.RGBA, lessonName string) color.RGBA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if _, exists := colorData[lessonName]; !exists {
		colorIndex := rng.Intn(len(colorList))
		colorData[lessonName] = colorList[colorIndex]
	}
	return colorData[lessonName]
}

func (o *ImgOption) processDayInfo(weekDay int, dayInfo map[string]any, colorData map[string]color.RGBA, context *gg.Context) {
	removeAllSpaces := func(str string) string {
		if str == "" {
			return ""
		}
		re := regexp.MustCompile(`\s+`)
		return re.ReplaceAllString(str, "")
	}

	for timeSlot, courseInfo := range dayInfo {
		if courseInfo == "noClass" {
			continue
		}

		lessonName := courseInfo.(map[string]string)["courseName"]
		teacher := courseInfo.(map[string]string)["teacher"]
		place := courseInfo.(map[string]string)["classroom"]
		lessonName = removeAllSpaces(lessonName)

		// 配置颜色
		colorInfo := o.IsExistAndColor(colorData, lessonName)

		// 绘制边框和矩形
		context.SetRGB(float64(colorInfo.R)/255.0, float64(colorInfo.G)/255.0, float64(colorInfo.B)/255.0)
		context.DrawRectangle(weekPlaceData[weekDay], hPlaceData[timeSlot], 245, 110)
		context.FillPreserve()
		context.Stroke()
		context.SetRGB(0, 0, 0)

		context.DrawString(lessonName, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+30)
		context.DrawString(place, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+65)
		context.DrawString(teacher, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+100)
	}
}

// CreateTplWithLocal 从本地文件创建模板
func (o *ImgOption) CreateTplWithLocal(imgTemplate string, cnameData map[int]map[string]any) *gg.Context {
	img, _ := gg.LoadImage(imgTemplate)
	context := gg.NewContextForImage(img)
	return o.create(context, cnameData)
}

// CreateTplWithCtx 从上下文创建模板
func (o *ImgOption) CreateTplWithCtx(ctx *gg.Context, cnameData map[int]map[string]any) *gg.Context {
	return o.create(ctx, cnameData)
}

func (o *ImgOption) create(context *gg.Context, cnameData map[int]map[string]any) *gg.Context {
	// 颜色数据 和 组合数据
	colorData := make(map[string]color.RGBA)
	_ = context.LoadFontFace(o.FontPath, 28)

	for weekDay, dayName := range cnameData {
		o.processDayInfo(weekDay, dayName, colorData, context)
	}
	return context
}
