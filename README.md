# @gkd-kit/inspect

- <https://i.gkd.li>

GKD 网页端审查工具, fork 自 [inspect](https://github.com/gkd-kit/inspect) 项目, 进行了一系列离线化改造，同时改造它更易于普通人使用。

## 编译运行流程

先安装依赖

```sh
pnpm install
```

然后编译项目

```sh
pnpm run build
```

最后`dist`目录里面就是编译好的成品，需要由golang后端启动才能正常使用。