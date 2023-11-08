<script setup lang='ts'>
import { NLayout, NLayoutContent, NLayoutSider } from 'naive-ui'
import { computed } from 'vue'
// import { Header } from './components'
// import { LeftSider, RightSider } from './layout'
import Index from './index.vue'
import { usePanelState } from '@/store/modules'

const panelState = usePanelState()

const leftSiderCollapsed = computed(() => panelState.leftSiderCollapsed)
const rightSiderCollapsed = computed(() => panelState.rightSiderCollapsed)
</script>

<template>
  <div class="h-full">
    <NLayout has-sider class="h-full">
      <NLayoutSider
        v-model:collapsed="leftSiderCollapsed"
        collapse-mode="transform"
        :collapsed-width="0"
        :width="240"
        bordered
      >
        <LeftSider />
      </NLayoutSider>
      <NLayoutContent>
        <NLayout has-sider sider-placement="right" class="h-full">
          <NLayoutContent class="h-full">
            <!-- 内容 -->
            <!-- <Header /> -->
            <div>
              <Index />
            </div>
          </NLayoutContent>
          <NLayoutSider
            v-model:collapsed="rightSiderCollapsed"
            collapse-mode="transform"
            :collapsed-width="0"
            :width="280"
            content-style="padding: 20px;"
            bordered
          >
            <RightSider />
          </NLayoutSider>
        </NLayout>
      </NLayoutContent>
    </NLayout>
  </div>
</template>
