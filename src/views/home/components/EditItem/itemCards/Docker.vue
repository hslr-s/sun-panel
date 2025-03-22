<script setup lang="ts">
import {
  NCheckbox, NCheckboxGroup,
  NFlex, NFormItem, NGrid, NGridItem, NInput, NSelect, NSpace, NSwitch, useMessage,
} from 'naive-ui'

import type { FormRules } from 'naive-ui'
import { onMounted, ref, watch } from 'vue'
import IconEditor from '../IconEditor.vue'
import UrlAddressInput from './components/UrlAddressInput.vue'
import type { SystemApp } from '@/utils/cmn/builtInApps'
import { t } from '@/locales'
import { AppLoader, ProAuthOverlay, ProTag, RoundCardModal, STip } from '@/components/common/'
import { getList as getContainerListApi } from '@/api/docker/container'
import { DockerContainerInfoState } from '@/enums/docker'
import { useProAuthorizeStore } from '@/store'

const props = defineProps<{
  itemInfo: Panel.Info
  rules: FormRules
}>()

const emit = defineEmits<{

  (e: 'update:itemInfo', info?: Panel.Info): void
  (e: 'update:rules', info?: FormRules): void
}>()

// const getIconLoading = ref([false, false])
const containerListLoading = ref(false)
const ms = useMessage()
const proAuthStore = useProAuthorizeStore()
const containers = ref<Docker.Container.Info[]>([])
async function getContainerList() {
  containerListLoading.value = true
  await getContainerListApi<Common.ListResponse<Docker.Container.Info[]>>().then((res) => {
    containers.value = res.data.list
  }).catch((res) => {
    if (res.code === -1) {
      ms.error(`${t('iconItem.dockerGetContainerListFailed')}: ${res.msg}`, {
        closable: true,
        duration: 10000,
        keepAliveOnHover: true,
      })
    }
  })
  containerListLoading.value = false
}

const item = ref(props.itemInfo)
const builtinAppModalShow = ref(false)

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
  expandParam: {
    containerId: {
      required: true,
      trigger: 'blur',
      type: 'string',
      message: t('form.required'),
    },
  },
}

const extendParam: Panel.ItemCardExpandParam.Docker = {
  containerId: '',
  containerName: '',
  openPageAndStart: false,
  cardType: '',
  showBadgeStatus: [],
  openPageAndStartDelaySecond: 0,
}

const model = ref(props.itemInfo)

watch(model.value, (newVal) => {
  emit('update:itemInfo', newVal)
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
  getContainerList()
  // 设置默认值
  if (!props.itemInfo.expandParam) {
    model.value.expandParam = { ...extendParam }
  }
  else {
    const ep = { ...extendParam }
    // 给子字段赋值，防止无用的字段
    ep.containerId = props.itemInfo.expandParam.containerId
    ep.openPageAndStart = props.itemInfo.expandParam.openPageAndStart
    ep.cardType = props.itemInfo.expandParam.cardType
    ep.showBadgeStatus = props.itemInfo.expandParam.showBadgeStatus
    ep.openPageAndStartDelaySecond = props.itemInfo.expandParam.openPageAndStartDelaySecond
    ep.containerName = props.itemInfo.expandParam.containerName

    model.value.expandParam = ep
  }
})

function handleUpdateshowBadgeStatusValue(value: (string | number)[]) {
  model.value.expandParam.showBadgeStatus = value
  // message.info(JSON.stringify(value))
  // console.log(JSON.stringify(value))
}

function handleChange(value: string, appInfo: SystemApp) {
  const containerName = appInfo.name
  if (model.value.title === '')
    model.value.title = containerName

  model.value.expandParam.containerName = containerName
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

    <!-- <NDivider class="!my-2" /> -->

    <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">
        <NFormItem
          path="expandParam.containerId"
        >
          <template #label>
            <span class="flex items-center font-bold">
              <span>
                {{ $t('iconItem.containerName') }}
              </span>
              <STip class="ml-1">
                {{ t("iconItem.containerFormTip") }}
              </STip>
            </span>
          </template>
          <NSelect
            v-model:value="model.expandParam.containerId"
            :loading="containerListLoading"
            filterable
            label-field="name"
            value-field="id"
            :options="containers"
            @update:value="handleChange"
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

    <NFormItem>
      <template #label>
        <span class="font-bold">
          {{ t('iconItem.showBadgeStatus') }}
        </span>
      </template>

      <NCheckboxGroup :value="model.expandParam.showBadgeStatus" @update:value="handleUpdateshowBadgeStatusValue">
        <NSpace item-style="display: flex;" align="center">
          <NCheckbox :value="DockerContainerInfoState.running" :label="t('iconItem.showBadgeStatusRunning')" />
          <NCheckbox :value="DockerContainerInfoState.restarting" :label="t('iconItem.showBadgeStatusRestarting')" />
          <NCheckbox :value="DockerContainerInfoState.exited" :label="t('iconItem.showBadgeStatusExited')" />
        </NSpace>
      </NCheckboxGroup>
    </NFormItem>

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
    </NFormItem> -->

    <!-- <NFormItem path="lanUrl" :label="$t('iconItem.lanUrl')">
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

    <!-- <NFormItem path="description" :label="$t('common.description')">
      <NInput v-model:value="model.description" type="text" show-count :maxlength="100" clearable />
    </NFormItem> -->
    <NFormItem>
      <template #label>
        <span class="font-bold">
          {{ $t('iconItem.openPageAndStartContainer') }}
        </span>
      </template>

      <ProAuthOverlay :show="proAuthStore.isExpired">
        <NSwitch v-model:value="model.expandParam.openPageAndStart" />
        <!-- <div v-if="model.expandParam.openPageAndStart" class="flex items-center">
        <span class="mr-2 ml-5">容器启动后延时打开页面</span>
        <NInputNumber
          v-model:value="model.expandParam.openPageAndStartDelaySecond"
          style="width: 80px;"
          :min="0"
          :max="60"
        />
        <span class="ml-2">s</span>
      </div> -->
        <ProTag class="ml-2" />
      </ProAuthOverlay>
    </NFormItem>

    <NFormItem path="openMethod">
      <template #label>
        <span class="font-bold">
          {{ $t('iconItem.openMethod') }}
        </span>
      </template>

      <NSelect v-model:value="model.openMethod" :options="options" />
    </NFormItem>
    <!-- <NGrid cols="2" :x-gap="10" item-responsive>
      <NGridItem span="2 500:1">

      </NGridItem>
      <NGridItem span="2 500:1">
        <NFormItem path="openMethod" :label="$t('iconItem.openMethod')">
          <NSelect v-model:value="model.openMethod" :options="options" />
        </NFormItem>
      </NGridItem>
    </NGrid> -->

    <RoundCardModal v-model:show="builtinAppModalShow">
      <AppLoader :component-name="item.expandParam.systemAppName" />
      <!--  -->
    </RoundCardModal>
  </div>
</template>
