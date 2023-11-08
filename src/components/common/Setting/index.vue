<script setup lang='ts'>
import { computed, ref } from 'vue'
import { NModal, NTabPane, NTabs } from 'naive-ui'
import UserInfo from './UserInfo.vue'
// import Advanced from './Advanced.vue'
import About from './About.vue'
import InviteUsers from './InviteUsers.vue'

// import { useAuthStore } from '@/store'
import { SvgIcon } from '@/components/common'

interface Props {
  visible: boolean
}

interface Emit {
  (e: 'update:visible', visible: boolean): void
}

const props = defineProps<Props>()

const emit = defineEmits<Emit>()

// const authStore = useAuthStore()

// const isChatGPTAPI = computed<boolean>(() => !!authStore.isChatGPTAPI)

const active = ref('General')

const show = computed({
  get() {
    return props.visible
  },
  set(visible: boolean) {
    emit('update:visible', visible)
  },
})
</script>

<template>
  <NModal v-model:show="show" :auto-focus="false" preset="card" style="width: 95%; max-width: 640px">
    <div>
      <NTabs v-model:value="active" type="line" animated>
        <NTabPane name="General" tab="General">
          <template #tab>
            <SvgIcon class="text-lg" icon="ri:file-user-line" />
            <span class="ml-2">设置</span>
          </template>
          <div class="min-h-[100px]">
            <UserInfo />
          </div>
        </NTabPane>
        <!-- <NTabPane v-if="isChatGPTAPI" name="Advanced" tab="Advanced">
          <template #tab>
            <SvgIcon class="text-lg" icon="ri:equalizer-line" />
            <span class="ml-2">{{ $t('setting.advanced') }}</span>
          </template>
          <div class="min-h-[100px]">
            <Advanced />
          </div>
        </NTabPane> -->
        <NTabPane name="InviteUsers" tab="InviteUsers">
          <template #tab>
            <SvgIcon class="text-lg" icon="mdi:invite" />
            <span class="ml-2">邀请新用户</span>
          </template>
          <InviteUsers />
        </NTabPane>

        <NTabPane name="About" tab="about">
          <template #tab>
            <SvgIcon class="text-lg" icon="mdi:about-circle-outline" />
            <span class="ml-2">关于</span>
          </template>
          <About />
        </NTabPane>
      </NTabs>
    </div>
  </NModal>
</template>
