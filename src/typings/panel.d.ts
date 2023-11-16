declare namespace Panel {

    interface Info extends ItemInfo {

    }

    interface ItemInfo extends Common.InfoBase {
        icon: ItemIcon |null
        title: string
        url: string
        lanUrl?: string
        description?: string
        openMethod: number
    }

    interface ItemIconGroup extends Common.InfoBase {
        icon?: string
        title?: string
        sort?:number
    }

    interface ItemIcon {
        itemType: number
        src ?: string
        text ?: string
        bgColor ?: string
    }

    interface State {
        rightSiderCollapsed: boolean
        leftSiderCollapsed: boolean
        networkMode:PanelStateNetworkModeEnum | null
        panelConfig:panelConfig
    }

    interface panelConfig{
        backgroundImageSrc?:string
        backgroundBlur?:number
        iconStyle?:PanelPanelConfigStyleEnum
        iconTextColor?:string
        logoText?:string
        logoImageSrc?:string
        clockShowSecond?:boolean
        clockColor?:string

    }

    interface userConfig{
        panel:panelConfig
        searchEngine?:any
    }
}

