初始化项目，使用go-zero框架，使用mysql和redis
生成README.md文件
配置CacheRedis改名为Redis
redis只有一个，不是集群

项目需求是将java项目的http接口转化为go语言接口，原java项目是:docs/java_src/RuoYi-Vue-Plus，docs/api是其http接口。plus-ui是vue写的web端。
迁移分两部分：common和sys。common是通用接口，不需要验证登录。sys需要验证登录后jwt token。
第一步：先实现common的接口：/auth/code 生成验证码、/auth/login 登录、/auth/tenant/list 租户列表
/auth/code 暂只支持生成4位数字验证码

我建了main.api文件，包含"auth.api"和"system.api"，以后只需要goctl api go -api api/main.api -dir .

auth.api和system.api的service名统一为admin-api

调用goctl model生成表模板代码：sys_client ,sys_config ,sys_dept ,sys_dict_data ,sys_dict_type ,sys_logininfor ,sys_menu ,sys_notice ,sys_oper_log ,sys_oss ,sys_oss_config ,sys_post ,sys_role ,sys_role_dept ,sys_role_menu ,sys_social ,sys_tenant ,sys_tenant_package ,sys_user ,sys_user_post ,sys_user_role