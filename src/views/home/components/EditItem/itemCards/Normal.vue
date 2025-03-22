<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  NFlex, NFormItem, NGrid, NGridItem,
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

    <!-- <NDivider class="!my-2" /> -->

    <NGrid cols="2" :x-gap="10" item-responsive>
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
      <NGridItem span="2 500:1">
        <NFormItem path="description">
          <template #label>
            <span class="font-bold">
              {{ $t('common.description') }}
            </span>
          </template>
          <NInput v-model:value="model.description" type="text" show-count :maxlength="100" clearable />
        </NFormItem>
      </NGridItem>
    </NGrid>

    <!-- <NDivider class="!my-2" /> -->

    <NFormItem path="url">
      <template #label>
        <span class="font-bold">
          {{ $t('iconItem.defaultUrl') }}
        </span>
      </template>
      <div class="w-full">
        <UrlAddressInput v-model:url="model.url" v-model:icon-info="model.icon" placeholder="http(s)://" />
      </div>
    </NFormItem>

    <NFormItem path="lanUrl">
      <template #label>
        <span class="font-bold">
          {{ $t('iconItem.lanUrl') }}
        </span>
      </template>
      <div class="w-full">
        <UrlAddressInput v-model:url="model.lanUrl" v-model:icon-info="model.icon" :placeholder="$t('iconItem.lanUrlInputPlaceholder')" />
      </div>
    </NFormItem>

    <!-- <NFormItem path="url" :label="$t('iconItem.defaultUrl')">
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

    <NFormItem path="openMethod">
      <template #label>
        <span class="font-bold">
          {{ $t('iconItem.openMethod') }}
        </span>
      </template>
      <NSelect v-model:value="model.openMethod" :options="options" />
    </NFormItem>
  </div>
</template>
