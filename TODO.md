# 时频学习平台 - 后端代码

## todo

### 短期

**feat:**

+ 加入登录接口
+ 加入 CHANGELOG.md
+ 支持从环境变量中读取 jwt 相关配置
+ 加入 jwt 黑名单机制，在角色改变的时候，将 token 禁用掉
+ 写一个通过 token 判断是否是管理员的接口，传 ctx 作为参数，用于在接口中操作其他人的信息时校验。

**bug:**

**improve:**

+ ~~分页相关参数及返回，重构为使用 Page struct，方便统一管理~~
+ ~~数据库替换为 pg，学习下 go-pg/pg 用法~~
+ ~~helm 命令，timeout参数设置长一点，加上回滚参数~~
+ ~~配置下 iris 的超时时间，看看有没有 mode 选项~~
+ 日志前面加上带颜色的【级别】前缀
+ 生成环境关闭 swagger 文档
+ 开发环境加入数据库日志（失策，还是得把 mode 参数加回来）
+ 将 resp 封装到参数校验函数中去
+ 在 errcode.Error 中加入一个 errDebug 字段，用来放不对用户展示的 debug 错误信息，与 errDetails 字段区分开，未来可以考虑在前端默认展示 errMsg，弄个按钮，点击出个显示 errDetails
  的弹窗。
+ 日志库中使用的 color.New，优化一下，不要每次都 New

**docs:**

+ ~~更新用户信息接口，返回为空，需要补充下~~
+ readme.md 和 todo.md 区分开

### 中期

+ 重构为 iris 框架
+ 加入测试代码

### 长期

+ 看下 goframe 框架各个业务模块的代码，学习下经验
