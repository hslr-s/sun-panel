<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import {
  NCard, NDivider, NFlex, NFormItem, NGrid, NGridItem,
  NInput, NSelect,
} from 'naive-ui'
import type { FormRules } from 'naive-ui'
import IconEditor from '../IconEditor.vue'
import UrlAddressInput from './components/UrlAddressInput.vue'
import { t } from '@/locales'
import { STip } from '@/components/common'

const props = defineProps<{
  itemInfo: Panel.Info
  rules: FormRules
}>()
const emit = defineEmits<{

  (e: 'update:itemInfo', info?: Panel.Info): void
  (e: 'update:rules', info?: FormRules): void
}>()
// const getIconLoading = ref([false, false])
// const ms = useMessage()
const model = ref(props.itemInfo)

// 打开方式
const options = [
  {
    default: true,
    label: t('iconItem.currentPageOpen'),
    value: 1,
  },
  {
    label: t('iconItem.newWindowOpen'),
    value: 2,
  },
  {
    label: t('iconItem.currentPageLayerOpen'),
    value: 3,
  },
]

const rulesExtends: FormRules = {
  url: [
    {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  ],
}

// async function getIconByUrl(url: string, loadingIndex: number) {
//   getIconLoading.value[loadingIndex] = true
//   try {
//     const { code, data } = await getSiteFavicon<{ iconUrl: string }>(url)
//     if (code === 0) {
//       model.value.icon = {
//         itemType: 2,
//         src: data.iconUrl,
//       }
//     }
//     else {
//       ms.error(t('iconItem.geticonFail'))
//     }
//   }
//   catch (error) {
//     ms.error(t('iconItem.geticonFail'))
//   }
//   getIconLoading.value[loadingIndex] = false
// }

watch(model.value, () => {
  // console.log('models 跟新', model.value)
  emit('update:itemInfo', model.value)
})

onMounted(() => {
  const r = { ...props.rules, ...rulesExtends }
  emit('update:rules', r)
})
</script>

<template>
  <div>
    <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">
        <NFormItem path="title" :label="$t('common.title')">
          <NInput v-model:value="model.title" type="text" show-count :maxlength="20" clearable />
        </NFormItem>
      </NGridItem>
      <NGridItem span="2 500:1">
        <NFormItem path="description" :label="$t('common.description')">
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" clearable />
        </NFormItem>
      </NGridItem>
    </NGrid>

    <NFormItem path="icon">
      <template #label>
        <NFlex align="center">
          {{ $t('iconItem.icon.iconStyle') }}
          <STip class="cursor-pointer">
            <p>{{ t('iconItem.icon.iconSizeTip') }}</p>
            <p>{{ t('iconItem.icon.iconOnlineTip') }}</p>
          </STip>
        </NFlex>
      </template>
      <NCard size="small">
        <IconEditor v-model:item-icon="model.icon" />
      </NCard>
    </NFormItem>

    <NDivider class="!my-2" />

    <NFormItem path="url" :label="`${$t('iconItem.url')}`">
      <div class="w-full">
        <UrlAddressInput v-model:url="model.url" v-model:icon-info="model.icon" />
      </div>
    </NFormItem>

    <NFormItem path="lanUrl" :label="`${$t('iconItem.lanUrl')}`">
      <div class="w-full">
        <UrlAddressInput v-model:url="model.lanUrl" v-model:icon-info="model.icon" />
      </div>
    </NFormItem>

    <!-- <NFormItem path="url" :label="$t('iconItem.url')">
      <div class="w-full">
        <NInputGroup>
          <NInput v-model:value="model.url" type="text" :maxlength="1000" placeholder="http(s)://" />
          <NButton :disabled="!model.url" :loading="getIconLoading[0]" @click="getIconByUrl(model.url, 0)">
            {{ $t('iconItem.getIcon') }}
          </NButton>
        </NInputGroup>
        <span v-if="!verifyUrlHttp(model.url)" class="text-xs text-orange-400">
          {{ t("iconItem.urlNoHttpStartWarn") }}
        </span>
      </div>
    </NFormItem>

    <NFormItem path="lanUrl" :label="$t('iconItem.lanUrl')">
      <div class="w-full">
        <NInputGroup>
          <NInput v-model:value="model.lanUrl" type="text" :maxlength="1000" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
          <NButton :disabled="!model.lanUrl" :loading="getIconLoading[1]" @click="getIconByUrl(model.lanUrl || '', 1)">
            {{ $t('iconItem.getIcon') }}
          </NButton>
        </NInputGroup>
        <span v-if="!verifyUrlHttp(model.lanUrl || '')" class="text-xs text-orange-400">
          {{ t("iconItem.urlNoHttpStartWarn") }}
        </span>
      </div>
    </NFormItem> -->

    <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
      <NSelect v-model:value="model.openMethod" :options="options" />
    </NFormItem>
  </div>
</template>
