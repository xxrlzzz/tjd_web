/conf : 配置文件  
/middleware 中间件  
/pkg： 第三方包 封装  
/routers： api url相应 相当于controller  
/runtime： 运行时文件存放路径  
/service： service层  
/models： model层  

```
// 生成证书
openssl req -x509 -out localhost.crt -keyout localhost.key \
  -newkey rsa:2048 -nodes -sha256 \
  -subj '/CN=121.5.156.132' -extensions EXT -config <( \
   printf "[dn]\nCN=121.5.156.132\n[req]\ndistinguished_name = dn\n[EXT]\nsubjectAltName=DNS:tjd_service\nkeyUsage=digitalSignature\nextendedKeyUsage=serverAuth")
   ```



## 开发日志

11.3 下载项目,研究框架  
     明日计划: 完成用户逻辑初步设计,用户与登录模块, 查找加密模块
     
11.9 加入了百度的api框架, 完成了token的逻辑  
     明日计划: 完成加密模块, 完成日志逻辑