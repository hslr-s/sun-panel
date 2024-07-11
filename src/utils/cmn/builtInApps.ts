import { setLocale, t } from '@/locales'
import { useAppStore } from '@/store'

import UserInfoSvg from '@/assets/svg-icons/material-symbols-person-edit-outline-rounded.svg'
import StyleSvg from '@/assets/svg-icons/ion-color-palette-outline.svg'
import ItemGroupManageSvg from '@/assets/svg-icons/material-symbols-ad-group-outline-rounded.svg'
import UploadFileManagerSvg from '@/assets/svg-icons/tabler-file-upload.svg'
import ImportExportSvg from '@/assets/svg-icons/icon-park-outline-import-and-export.svg'
import OpenAPISvg from '@/assets/svg-icons/tabler-api.svg'
import AboutSvg from '@/assets/svg-icons/lucide-info.svg'

import ProAuthSvg from '@/assets/svg-icons/tabler-award.svg'

// [未知原因]，需要深究（猜测无法监听）-此处需要单独设置一下语言否则不生效 // 临时解决
const appStore = useAppStore()
setLocale(appStore.language)

export interface SystemApp {
  name: string
  componentName: string
  icon: string
  auth?: number
  iconSrc?: string
}

export const systemApps: SystemApp[] = [
  {
    name: t('apps.userInfo.appName'),
    componentName: 'UserInfo',
    icon: 'material-symbols:person-edit-outline-rounded',
    iconSrc: UserInfoSvg,
  },
  {
    name: t('apps.baseSettings.appName'),
    componentName: 'Style',
    icon: 'ion:color-palette-outline',
    iconSrc: StyleSvg,
  },
  {
    name: t('apps.itemGroupManage.appName'),
    componentName: 'ItemGroupManage',
    icon: 'material-symbols:ad-group-outline-rounded',
    iconSrc: ItemGroupManageSvg,
  },
  {
    name: 'Docker',
    componentName: 'DockerManage',
    icon: 'tabler-brand-docker',
    auth: 1,
  },
  {
    name: t('apps.uploadsFileManager.appName'),
    componentName: 'UploadFileManager',
    icon: 'tabler:file-upload',
    iconSrc: UploadFileManagerSvg,
  },
  {
    name: t('apps.exportImport.appName'),
    componentName: 'ImportExport',
    icon: 'icon-park-outline:import-and-export',
    iconSrc: ImportExportSvg,
  },
  {
    name: 'Open API',
    componentName: 'OpenApi',
    icon: 'tabler:api',
    iconSrc: OpenAPISvg,
  },

  {
    name: t('apps.about.appName'),
    componentName: 'About',
    icon: 'lucide:info',
    iconSrc: AboutSvg,
  },
]

// 管理账号的应用
export const adminSystemApps: SystemApp[] = [
  {
    name: t('adminSettingUsers.appName'),
    componentName: 'Users',
    icon: 'lucide:users',
    auth: 1,
    iconSrc: OpenAPISvg,
  },
  {
    name: t('apps.globalSetting.appName'),
    componentName: 'GlobalSetting',
    icon: 'ic:outline-build-circle',
    auth: 1,
    iconSrc: OpenAPISvg,

  },
  {
    name: t('proAuth.appName'),
    componentName: 'ProAuth',
    icon: 'tabler:award',
    auth: 1,
    iconSrc: ProAuthSvg,
  },
]
