# ShareVault

ShareVault 是一個基於 Go 語言開發的文件分享平台，提供用戶註冊、文件上傳下載以及 AWS S3 整合功能。

## 功能特色

- 用戶註冊與 JWT 身份驗證
- 文件上傳、下載與管理
- AWS S3 雲端儲存整合
- PostgreSQL 資料庫支援
- RESTful API 架構
- Docker 容器化部署

## 技術

- **後端框架**: Gin (Go)
- **資料庫**: PostgreSQL
- **雲端儲存**: AWS S3
- **身份驗證**: JWT
- **容器化**: Docker & Docker Compose
- **開發工具**: Air (熱更新)

## 快速開始

### 前置需求

- Go 1.24.5+
- Docker & Docker Compose
- PostgreSQL (或使用 Docker)
- AWS 帳戶

### 安裝步驟

1. **複製專案**
   ```bash
   git clone https://github.com/LeeChasel/shareVault.git
   cd shareVault
   ```

2. **安裝依賴**
   ```bash
   make deps
   ```

3. **設定環境變數**

   複製 `.env.example` 到 `.env` 並配置所需參數：
   ```bash
   cp .env.example .env
   ```

### 環境變數設定

編輯 `.env` 文件，設定以下參數：

```env
# JWT 密鑰 (請使用安全的隨機字串)
JWT_SECRET=your-secure-jwt-secret-key

# 資料庫設定
DB_HOST=localhost
DB_PORT=5332
DB_USER=sv_user
DB_PASSWORD=sv_pass
DB_NAME=shareVault

# AWS S3 設定
AWS_REGION=ap-east-1
AWS_S3_BUCKET=your-s3-bucket-name
```

#### 重要說明：

- **JWT_SECRET**: 請使用強密碼生成器產生至少 32 位元的隨機字串
- **資料庫設定**: 如果使用 Docker，請保持預設值；如果使用外部資料庫，請修改相應參數
- **AWS 設定**: 需要有效的 AWS 認證，可透過 AWS CLI 配置或 IAM 角色

### AWS 認證設定

ShareVault 需要 AWS 認證來存取 S3 服務。您可以選擇以下任一方式：

1. **使用 AWS CLI 配置**
   ```bash
   aws configure
   ```

2. **使用環境變數**
   ```bash
   export AWS_ACCESS_KEY_ID=your-access-key
   export AWS_SECRET_ACCESS_KEY=your-secret-key
   ```

3. **使用 IAM 角色** (推薦用於生產環境)

## 開發模式

### 使用 Docker Compose 啟動資料庫

```bash
make docker-up
```

這將啟動：
- PostgreSQL 資料庫 (端口: 5332)
- Adminer 資料庫管理介面 (端口: 8081)

### 啟動開發伺服器

```bash
make dev
```

應用程式將在 http://localhost:8083 上運行，並支援熱重載。

### 其他開發指令

```bash
# 建置應用程式
make build

# 運行應用程式 (無熱重載)
make run

# 運行測試
make test

# 程式碼格式化
make fmt

# 程式碼檢查
make vet

# 程式碼 Lint
make lint

# 清理建置檔案
make clean
```

## 生產部署

1. **準備生產環境變數**
   ```bash
   # 設定安全的 JWT 密鑰
   JWT_SECRET=your-production-jwt-secret

   # 配置生產資料庫
   DB_HOST=your-production-db-host
   DB_PORT=5432
   DB_USER=your-production-db-user
   DB_PASSWORD=your-production-db-password
   DB_NAME=your-production-db-name

   # 配置 AWS S3
   AWS_REGION=your-preferred-region
   AWS_S3_BUCKET=your-production-s3-bucket
   ```

2. **建置應用程式**
   ```bash
   make build
   ```

3. **啟動服務**
   ```bash
   ./bin/shareVault
   ```

## API 端點

應用程式提供 RESTful API，主要端點包括：

- 用戶認證相關端點
- 文件上傳/下載端點
- 文件管理端點

詳細 API 文件請參考 `internal/api/routes` 目錄。

## 專案結構

```
shareVault/
├── cmd/shareVault/          # 應用程式入口點
├── internal/
│   ├── api/                 # API 路由與處理器
│   ├── models/              # 資料模型
│   ├── repository/          # 資料存取層
│   ├── service/             # 業務邏輯層
│   ├── constants/           # 常數定義
│   └── utils/               # 工具函數
├── configs/                 # 配置檔案
├── deployments/             # 部署相關檔案
├── .env.example             # 環境變數範例
├── Makefile                 # 建置指令
└── go.mod                   # Go 模組定義
```

## 資料庫管理

### 使用 Adminer (推薦)

當 Docker 服務運行時，您可以透過 http://localhost:8081 存取 Adminer：

- 系統: PostgreSQL
- 伺服器: db
- 使用者名稱: sv_user
- 密碼: sv_pass
- 資料庫: shareVault

### 直接連線

```bash
psql -h localhost -p 5332 -U sv_user -d shareVault
```

## 故障排除

### 日誌檢查

```bash
# 查看 Docker 服務日誌
make docker-logs

# 查看應用程式建置錯誤
cat build-errors.log
```