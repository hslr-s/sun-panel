module.exports = {
  root: true,
  extends: ['@antfu'],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'warn', // 取消打印标红提醒
    'curly': 'off', // 这样配置确保在所有的if和else语句中都使用大括号
  },
}
