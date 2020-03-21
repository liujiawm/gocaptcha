# gocaptcha

图片验证码,因需求原因,在[https://github.com/steambap/captcha](https://github.com/steambap/captcha "https://github.com/steambap/captcha")的基础上简化修改而来,因使用要求与其有较大的不同,所以特此新建

字体使用go-bindata,字体在网上下载的两种字体,其版权问题不清楚


## panic "integer divide by zero"
 如果在github.com/golang/freetype/truetype/hint.go的函数mulDiv出除数不能为0的panic是因为字体造成,请换其他字体