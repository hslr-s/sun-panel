<script setup lang="ts">
import { NTag, NTooltip } from 'naive-ui'
import { SvgIcon } from '..'
import { useProAuthorizeStore } from '@/store'
defineProps<{
  showIcon?: boolean
  content?: string
}>()

const proAuthorizeStore = useProAuthorizeStore()
</script>

<template>
  <span v-if="!proAuthorizeStore.hideProBadge && showIcon">
    <NTooltip>
      <template #trigger>
        <NTag size="small" style="border-radius: 10px;font-weight: 800;" type="warning">
          <template #icon>
            <slot v-if="showIcon" name="icon">
              <SvgIcon
                icon="pajamas:question-o"
              />
            </slot>
          </template>
          PRO
        </NTag>
      </template>
      <slot>
        {{ content }}
      </slot>
    </NTooltip>
  </span>

  <NTag v-else-if="!proAuthorizeStore.hideProBadge" size="small" style="border-radius: 10px;font-weight: 800;" type="warning">
    <template #icon>
      <slot v-if="content && showIcon" name="icon">
        <SvgIcon
          icon="pajamas:question-o"
        />
      </slot>
    </template>
    PRO
  </NTag>
</template>
