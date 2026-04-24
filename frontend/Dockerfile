# 使用 Node.js 20 作为基础镜像
FROM node:20-alpine

# 设置工作目录
WORKDIR /app

# 设置npm使用国内镜像源（可选，提高依赖安装速度）
RUN npm config set registry https://registry.npmmirror.com

# 复制 package.json 和 package-lock.json
COPY package*.json ./

# 安装依赖
RUN npm install

# 复制项目文件
COPY . .

# 构建生产版本
RUN npm run build

# 安装静态文件服务器
RUN npm install -g serve

# 暴露端口
EXPOSE 3000

# 启动服务
CMD ["serve", "-s", "dist"]
