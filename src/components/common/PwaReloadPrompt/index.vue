<script setup lang="ts">
import { useRegisterSW } from 'virtual:pwa-register/vue'
import { NButton, useNotification } from 'naive-ui'

import { h, watch } from 'vue'
import { t } from '@/locales'

const { needRefresh, updateServiceWorker } = useRegisterSW()

const notification = useNotification()

watch(() => needRefresh.value, (v: boolean) => {
  if (v)
    openNotification()
})

function openNotification() {
  const n = notification.warning({
    content: () => h('div', t('components.pwaReloadPrompt.promptContent')),
    action: () =>
      [
        h(
          NButton,
          {
            text: true,
            size: 'small',
            type: 'primary',
            onClick: () => {
              updateServiceWorker()
              n.destroy()
            },
          },
          {
            default: () => t('common.reload'),
          },
        ),
        h(
          NButton,
          {
            style: 'margin-left: 8px',
            text: true,
            size: 'small',
            onClick: () => {
              n.destroy()
            },
          },
          {
            default: () => t('common.cancel'),
          },
        ),
      ],
    onClose: () => {
      needRefresh.value = false
    },
  })
}
// onMounted(() => openNotification()) // 测试弹窗
</script>
