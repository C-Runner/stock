# Stock Analysis System

股票分析系统，提供股票数据查询和技术分析功能。

## 技术栈

### 前端
- Vue 3
- Naive UI (UI组件库)
- TypeScript
- Vite

### 后端
- Go
- Gin (Web框架)
- MongoDB

## 项目结构

```
stock/
├── frontend/           # 前端项目
│   ├── src/           # Vue 源代码
│   ├── public/        # 静态资源
│   └── package.json   # 前端依赖
├── backend/           # 后端项目
│   ├── handlers/      # HTTP 处理器
│   ├── services/      # 业务逻辑
│   ├── models/        # 数据模型
│   ├── config/        # 配置
│   └── main.go        # 入口文件
└── docker-compose.yml # Docker 配置
```

## 快速开始

### 前置要求
- Node.js 18+
- Go 1.25+
- Docker 和 Docker Compose

### 1. 启动数据库

```bash
docker-compose up -d
```

这将启动 MongoDB 服务，监听端口 27017。

### 2. 启动后端

```bash
cd backend
go run main.go
```

后端服务将运行在 `http://localhost:8080`

### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端服务将运行在 `http://localhost:5173`

## API 端点

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /api/stocks | 获取股票列表 |
| GET | /api/stocks/:symbol | 获取股票详情 |
| GET | /api/analysis/:symbol | 获取股票分析 |

## 开发

### 前端开发
```bash
cd frontend
npm run dev      # 开发模式
npm run build    # 构建生产版本
npm run preview  # 预览生产版本
```

### 后端开发
```bash
cd backend
go run main.go   # 启动服务
```
