declare namespace GlobalSetting {

    interface SiteSetting {
        favicoUrl: string
        title: string
        loginBackgroundUrl: string   
        isEnableLoginCaptcha:boolean
        isDisableItemCardCache:boolean
    }

    interface GlobalSettingCustomJsAndCss {
        jsContent: string;
        cssContent: string;
    }

    interface GlobalSettingGetCustomPathReq {
        customPath: string;
    }
    

}
