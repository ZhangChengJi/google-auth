# google-auth
google-authenticator 双因子认证和zabbix结合 实现 双认证(等保要求)



### 使用步骤

1. 修改redis连接，修改db.go文件，修改成自己的redis地址和密码

   ```go
   // 初始化连接
   func initClient() (err error) {
   	rdb = redis.NewClient(&redis.Options{
   		Addr:     "172.16.1.3:6379",      
   		Password: "djs@12316", // no password set
   		DB:       0,           // use default DB
   		PoolSize: 10,          // 连接池大小
   	})
   
   ```

2. 执行交叉编译打包，这里写的是在liunx环境编译

   ```sh
   $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
   ```

3. 把编译后的`google-auth`二进制文件上传到liunx服务器，在通缉目录建立image文件夹

   ```sh
   //比如我在liunx root目录进行创建一个叫google-auth的文件夹
   $ makir google-auth  
   //上传`google-auth`二进制文件到google-auth文件夹内
   $ cd google-auth
   $ rz google-auth //上传二进制文件
   $ chmod a+x google-auth //授权
   $ makir image //创建image文件用于存放二维码图片
   $ nohup ./google-auth & //后台启动
   ```

4. 启动`google-auth`测试注册和验证码验证接口

   ```sh
   // 用户注册 //比如 http://localhost:8082/createCode?issuer=chengzhenyuan
   $ http://[google-authd地址+端口]/createCode?issuer=[用户名]
    //返回 {"code":1,"msg":"chengzhenyuan 用户注册成功"}  这个时候就在image目录生成了图片
   ```

5. 下载 Authenticator APP ,扫描image 文件下的二维码图片，就会显示当前这个人的口令码了

   ![app](https://gitee.com/zhangchengji/pic/raw/master/uPic/app.png)

 

6. 测试口令是否生效

   ```sh
   $ curl -i -X GET \
    'http://localhost:8082/verifyCode?issuer=chengzhenyuan&code=269761'
    // 返回 {"code":0,"msg":"Google验证码验证成功"}
   
   ```

   

7. 修改zabbix php替换文件，主要有两个文件`general.login.php`和 `index.php`

   在当前项目下面下有这两个文件 直接去你的zabbix 进行替换

   替换路径为`include/views/general.login.php`

   Zabbix 根目录`index.php`

   替换前注意index.php里面第73行google-auth地址改成自己部署的地址

   ```php
           $authflag=file_get_contents("http://[改成自己地址+端口]/verifyCode?issuer=".getRequest('name', ZBX_GUEST_USER)."&code=".getRequest('code', ''));
   
   ```

   

8. 登陆测试

   ![xf6BOJ](https://gitee.com/zhangchengji/pic/raw/master/uPic/xf6BOJ.png)
