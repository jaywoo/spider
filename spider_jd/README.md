spider_jd_file.go 将基础信息爬下来写入文件

productIntoDb.go  将信息写入到数据库

util/picasso.go   图床，根据平台不同需要自己重写

result.log  爬取结果实例

执行：./spider_jd_file 或 ./spider_jd_file spider_jd_test_url.log

JD的图站点都是直接渲染的，爬了一下基础信息，品牌，型号等

为了熟悉go的使用而写，主要用到http，file，json，等基础方法，未使用协程。



