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
        cleanupOutdatedCaches: true, // 清理旧缓存
        cacheId: env.VITE_APP_VERSION,
        clientsClaim: true,
        // skipWaiting: true,
        globPatterns: ['**/*.{js,css,html,ico,png,jpg,svg}'], // 缓存相关静态资源
        // runtimeCaching: [
        //   // 静态文件缓存优先策略
        //   {
        //     urlPattern: ({ request }) =>
        //       request.destination === 'script'
        //       || request.destination === 'style'
        //       || request.destination === 'image'
        //       || request.destination === 'font'
        //       || request.destination === 'document',
        //     handler: 'CacheFirst', // 优先从缓存中读取静态文件
        //     options: {
        //       cacheName: `static-resources-${env.VITE_APP_VERSION}`, // 静态资源缓存
        //       expiration: {
        //         maxEntries: 100, // 最大缓存条目数
        //         maxAgeSeconds: 30 * 24 * 60 * 60, // 缓存有效期，30天
        //       },
        //     },
        //   },
        // ],
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
