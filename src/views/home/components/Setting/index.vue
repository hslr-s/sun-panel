<script setup lang="ts">
import { computed } from 'vue'
import { NTabPane, NTabs } from 'naive-ui'
import Style from './tabs/Style.vue'
import About from './tabs/About.vue'
import Users from './tabs/Users.vue'
import UserInfo from './tabs/UserInfo.vue'
import ItemGroupManage from './tabs/ItemGroupManage.vue'
import ImportExport from './tabs/ImportExport.vue'
import { useAuthStore } from '@/store'
import { AdminAuthRole } from '@/enums/admin'
import { RoundCardModal } from '@/components/common'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
}>()

const authStore = useAuthStore()

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => {
    emit('update:visible', visible)
  },
})
</script>

<template>
  <div>
    <RoundCardModal
      v-model:show="show" title="设置"
      class="fixed right-[10px] bottom-[10px]"
      style="max-height: 700px;max-width: 600px;"
    >
      <NTabs type="line" size="small" animated>
        <NTabPane name="style" tab="面板样式">
          <Style />
        </NTabPane>
        <NTabPane name="itemGroupManage" tab="分组管理">
          <ItemGroupManage />
        </NTabPane>
        <NTabPane name="userInfo" tab="账号信息">
          <UserInfo />
        </NTabPane>
        <NTabPane name="ImportExport" tab="导入导出">
          <ImportExport />
        </NTabPane>
        <NTabPane name="about" tab="关于">
          <About />
        </NTabPane>
        <NTabPane v-if="authStore.userInfo?.role === AdminAuthRole.admin" name="password" tab="账号管理">
          <Users />
        </NTabPane>
      </NTabs>
    </RoundCardModal>
  </div>
</template>

<style scoped>
.text-shadow {
  text-shadow: 0px 0px 5px gray;
}
</style>
