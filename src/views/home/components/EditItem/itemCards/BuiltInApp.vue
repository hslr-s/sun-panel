<script setup lang="ts">
import { NFlex, NFormItem, NGrid, NGridItem, NInput, NSelect } from 'naive-ui'
import type { FormRules } from 'naive-ui'
import { onMounted, ref, watch } from 'vue'
import IconEditor from '../IconEditor.vue'
import { getAppList } from '../../common'
import type { SystemApp } from '@/utils/cmn/builtInApps'
import { systemApps } from '@/utils/cmn/builtInApps'
import { t } from '@/locales'
import { AppLoader, RoundCardModal, STip } from '@/components/common/'
import { useAuthStore } from '@/store'

const props = defineProps<{
  itemInfo: Panel.Info
  rules: FormRules
}>()

const emit = defineEmits<{

  (e: 'update:itemInfo', info?: Panel.Info): void
  (e: 'update:rules', info?: FormRules): void
}>()

const authStore = useAuthStore()
const item = ref(props.itemInfo)
const builtinAppModalShow = ref(false)
const apps = ref(systemApps)
const rulesExtends: FormRules = {
  url: [
    {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  ],
  expandParam: {
    systemAppName: {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  },
}

const extendParam: Panel.ItemCardExpandParam.SystemApp = {
  systemAppName: '',
}

const model = ref(props.itemInfo)

watch(model.value, () => {
  emit('update:itemInfo', model.value)
})

watch(
  () => props.rules,
  () => {
    const newRules = { ...props.rules, ...rulesExtends }
    // 防止重复 emit
    if (JSON.stringify(newRules) !== JSON.stringify(props.rules)) {
      emit('update:rules', newRules)
    }
  },
  { immediate: true },
)

onMounted(() => {
  apps.value = getAppList(systemApps, authStore.userInfo?.role || 0)

  // 设置默认值
  if (!props.itemInfo.expandParam) {
    model.value.expandParam = { ...extendParam }
  }
  else {
    const ep = { ...extendParam }
    // 给子字段赋值，防止无用的字段
    ep.systemAppName = props.itemInfo.expandParam.systemAppName
    model.value.expandParam = ep
  }
})

function handleSystemChange(value: string, appInfo: SystemApp) {
  model.value.title = appInfo.name
  model.value.icon = {
    text: appInfo.icon,
    itemType: 3,
  }
  model.value.description = t('iconItem.builtinApp')
}
</script>

<template>
  <div>
    <NFormItem path="icon">
      <template #label>
        <NFlex align="center">
          <span class="font-bold">
            {{ $t('iconItem.icon.iconStyle') }}
          </span>

          <STip class="cursor-pointer">
            <p>{{ t('iconItem.icon.iconSizeTip') }}</p>
            <p>{{ t('iconItem.icon.iconOnlineTip') }}</p>
          </STip>
        </NFlex>
      </template>
      <!-- <NCard size="small"> -->
      <IconEditor v-model:item-icon="model.icon" />
      <!-- </NCard> -->
    </NFormItem>

    <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">
        <NFormItem path="expandParam.systemAppName">
          <template #label>
            <span class="font-bold">
              {{ $t('iconItem.builtInAppOpen') }}
            </span>
          </template>

          <NSelect
            v-model:value="model.expandParam.systemAppName"
            label-field="name"
            value-field="componentName"
            :options="apps"
            filterable
            @update:value="handleSystemChange"
          />
        </NFormItem>
      </NGridItem>
      <NGridItem span="2 500:1">
        <NFormItem path="title">
          <template #label>
            <span class="font-bold">
              {{ $t('common.title') }}
            </span>
          </template>

          <NInput v-model:value="model.title" type="text" show-count :maxlength="20" clearable />
        </NFormItem>
      </NGridItem>
    </NGrid>

    <NFormItem path="description">
      <template #label>
        <span class="font-bold">
          {{ $t('common.description') }}
        </span>
      </template>
      <NInput v-model:value="model.description" type="text" show-count :maxlength="100" clearable />
    </NFormItem>

    <RoundCardModal v-model:show="builtinAppModalShow">
      <AppLoader :component-name="item.expandParam.systemAppName" />
      <!--  -->
    </RoundCardModal>
  </div>
</template>
