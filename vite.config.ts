import path from 'path'
import type { PluginOption } from 'vite'
import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { createSvgIconsPlugin } from 'vite-plugin-svg-icons'

function setupPlugins(env: ImportMetaEnv): PluginOption[] {
  return [
    vue(),

    env.VITE_GLOB_APP_PWA === 'true' && VitePWA({
      registerType: 'prompt',
      manifest: false,
      injectRegister: false,
      workbox: {
        cleanupOutdatedCaches: true,
        cacheId: env.VITE_APP_VERSION,
        clientsClaim: true,
        globPatterns: ['**/*.{js,css,html,ico,png,jpg,svg}'],
        globIgnores: ['**/custom/**'], // 排除 /custom/ 目录
        navigateFallback: '/index.html', // 回退页面
        navigateFallbackDenylist: [/^\/custom\//], // 排除 /custom/ 路径
        runtimeCaching: [
          {
            urlPattern: /^\/custom\//, // 匹配 /custom/ 路径
            handler: 'NetworkFirst', // 优先从网络获取
            options: {
              cacheName: 'custom-cache',
            },
          },
        ],
      },
      devOptions: {
        enabled: false,
      },
    }),

    createSvgIconsPlugin({
      iconDirs: [path.resolve(process.cwd(), 'src/assets/svg-icons')],
      symbolId: '[name]',
    }),
  ]
}

export default defineConfig((env) => {
  const viteEnv = loadEnv(env.mode, process.cwd()) as unknown as ImportMetaEnv

  return {
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
      },
    },
    plugins: setupPlugins(viteEnv),
    server: {
      host: '0.0.0.0',
      port: 1002,
      open: false,
      proxy: {
        '/api': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/api/', '/api/'),
        },
        '/uploads': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/uploads/', '/uploads/'),
        },
        '/custom': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/custom/', '/custom/'),
        },
        '/favicon.ico': {
          target: viteEnv.VITE_APP_API_BASE_URL,
          changeOrigin: true, // 允许跨域
          rewrite: path => path.replace('/favicon.ico', '/favicon.ico'),
        },

      },
    },
    build: {
      reportCompressedSize: false,
      sourcemap: false,
      commonjsOptions: {
        ignoreTryCatch: false,
      },
      terserOptions: {
        compress: {
          drop_console: true,
        },
      },
    },
  }
})
