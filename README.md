# go-ffmpeg

一. 工具一（所有执行文件必须添加到系统变量PATH下）：
>
        openssl
        ffmpeg
二. 生成apiDoc二进制文件(不需要重新生成doc文档跳过次步骤)：
>    
    工具二（所有执行文件必须添加到系统变量PATH下）：
        apidocjs:              
                npm install apidoc -g    
        statik:
                go get github.com/rakyll/statik
                
    运行：
        go generate
三. 运行（go version > 1.11，使用Redis db 1）
> 
    
    go run

PS: 
>
    api接口文档地址：
        http://localhost:8080/doc
    转换后视频文件地址：
        http://localhost:8080/transfer/:id