package timetable

import (
	"strconv"

	"github.com/fogleman/gg"
)

type TemplateOption struct {
	WeekTime   map[int]string
	CreateTime string
	Cname      string
}

func (o *ImgOption) drawClassName(cname string, createTime string) {
	o.Image.SetRGB(0, 0, 0)
	o.Image.DrawString(cname, 870, 55)
	o.Image.DrawString(createTime, 870, 90)
}

func (o *ImgOption) drawSelection() {
	for i := range 8 {
		o.Image.DrawString("第 "+strconv.Itoa(i+1)+" 节", 70, float64(160+125*i))
	}
}

func (o *ImgOption) drawTime(timeData map[int]string) {
	if timeData == nil {
		timeData = map[int]string{
			1: "08:00—08:45",
			2: "08:55—09:40",
			3: "10:10—10:55",
			4: "11:05—11:50",
			5: "14:30—15:15",
			6: "15:25—16:10",
			7: "16:20—17:05",
			8: "17:15—18:00",
		}
	}
	for i := range 8 {
		o.Image.DrawString(timeData[i+1], 45, float64(200+125*i))
	}
}

func (o *ImgOption) drawWeek() {
	days := []string{"星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"}
	dateNum := []float64{350, 600, 850, 1090, 1345, 1600, 1855}
	for i, text := range days {
		textWidth, _ := o.Image.MeasureString(text)
		x := dateNum[i] - textWidth/2
		o.Image.DrawString(text, x, 1165)
	}
}

func (o *ImgOption) drawCanvas() {
	o.Image.SetRGB(0.5, 0.5, 0.5) // 调至浅灰色
	// 格子线
	{
		for i := 0; i <= 6; i++ {
			o.Image.DrawLine(247+250*float64(i), 120, 247+250*float64(i), 1090)
		}
		for _, i := range []int{240, 360, 485, 620, 737, 855, 970} {
			o.Image.DrawLine(247, float64(i), 2000, float64(i))
		}
	}
	// 基础框线
	{
		o.Image.DrawLine(0, 120, 2000, 120)
		o.Image.DrawLine(0, 605, 2000, 605)
		o.Image.DrawLine(0, 1090, 2000, 1090)
	}
	o.Image.Stroke()
}

// CreateBasePhoto 创建基础模板图片
func (o *ImgOption) CreateBasePhoto(opt TemplateOption) *gg.Context {
	// 创建画布
	o.drawCanvas()
	// 顶部班级
	o.drawClassName(opt.Cname, opt.CreateTime)
	// 节次
	o.drawSelection()
	// 时间
	o.drawTime(opt.WeekTime)
	// 星期
	o.drawWeek()
	return o.Image
}
