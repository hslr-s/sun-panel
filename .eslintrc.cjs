module.exports = {
  root: true,
  extends: ['@antfu'],
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off', // 取消打印标红提醒
  },
}
