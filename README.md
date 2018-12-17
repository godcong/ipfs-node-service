# go-ffmpeg

1. 使用工具：

        openssl
        ffmpeg

2. 生成apiDoc二进制文件(不生成doc文档跳过次步骤)：
>    
    配置：
        apidocjs:              
                npm install apidoc -g    
        statik:
                go get github.com/rakyll/statik
                
    运行：
        go generate

3. 运行
    go run

PS: 
    api接口文档地址：http://localhost:8080/doc
    转换后视频文件地址：http://localhost:8080/transfer/:id