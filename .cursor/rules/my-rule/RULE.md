---
globs:
alwaysApply: false
---

## 必须支持的规则

- 自动把我发出的请求文字记录到cursor_chat.txt，你自己的操作日志不要添加。格式：第一行：时间，第二行：请求内容。每次聊天请求间空一行
- 响应json int64数据，统一以字符串输出。
- logic*.go不要直接写sql语句，把访问db的统一抽取到相关model类中。model包中*_gen.go是自动生成的不要修改，如如sysdict只在sysdicttypemodel.go中修改
- model类统一使用snake case风格，如sys_client_model.go、sys_client_model_gen.go。以后goctl model mysql需指定--style go_zero
- Total 总记录数 不转为string, 直接用int64类型。改一下*.api文件
- int64 json输出要用string类型，别手动改types.go文件

## goctl限制 
- goctl api go 统一使用snake case风格，迁移当前文件名为snake case风格。使用命令：goctl api go -api api/main.api -dir . --style go_zero。将所有handler和logic文件从camelCase重命名为snake_case风格（如configaddlogic.go -> config_add_logic.go），文件名已统一为snake_case，函数名和类型名保持PascalCase（Go标准）
- 使用命令：goctl api go -api api/main.api -dir . --style go_zero，别把main.api改成system.api。main.api是入口文件。
- 使用命令：goctl model 也指定 --style go_zero
- 使用命令：goctl model不忽略任何字段：goctl model -i="" ... (传空值).create_time、update_time由前端输入


## *.api限制
- 响应json int64数据，统一以字符串输出，但Total除外。
- 所有*.api文件不使用tag的options=X|Y 限定输入