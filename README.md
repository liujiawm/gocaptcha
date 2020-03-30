# gocaptcha

图片验证码,因需求原因,在[https://github.com/steambap/captcha](https://github.com/steambap/captcha "https://github.com/steambap/captcha")的基础上简化修改而来,因使用要求与其有较大的不同,所以特此新建

字体使用go-bindata,字体在网上下载的两种字体,其版权问题不清楚


## panic "integer divide by zero"
 如果在github.com/golang/freetype/truetype/hint.go的函数mulDiv出除数不能为0的panic是因为字体造成,请换其他字体

## 示例

```
	data,_ := gocaptcha.New(&gocaptcha.Options{
		CharPreset:"0123456789", // 数字作基数
		Curve:2,                 // 两条弧线
		Length:4,                // 长度为4的验证码
		Width:80,                // 图片宽
		Height:33,               // 图片高
	})

	// TODO 将data.Text保存用于验证

	// 图片显示data.Text,也可用data.EncodeB64string()返回base64
	// c为gin的*gin.Context,也可以用其他的io.Writer
	data.WriteImage(c.Writer)

```

## 效果
![](https://github.com/liujiawm/gocaptcha/blob/master/test3.png?raw=true)
![](https://github.com/liujiawm/gocaptcha/blob/master/test4.png?raw=true)