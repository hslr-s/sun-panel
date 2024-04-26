const fs = require('fs')
// const { execSync } = require('child_process')
const moment = require('moment')

// git 最新标签
// const latestTag = execSync('git describe --tags --abbrev=0').toString().trim()

// 设置默认时区为 'Asia/Shanghai'
const packDate = moment().utc().format('YYYYMMDD')

// 要追加的内容
const contentToAppend = `\nVITE_APP_VERSION=${packDate}`
// 读取文件原始内容
const envFilePath = '.env'
let envContent = fs.readFileSync(envFilePath, 'utf-8')

const versionRegex = /^VITE_APP_VERSION=.*$/m
if (versionRegex.test(envContent)) {
  // 使用正则表达式查找并替换 VITE_APP_VERSION=* 这一行
  envContent = envContent.replace(versionRegex, contentToAppend)
}
else {
  // 追加内容
  envContent = envContent + contentToAppend
}

// 将新内容写回 .env 文件
fs.writeFileSync(envFilePath, envContent)

console.log('update to .env file.', contentToAppend)
