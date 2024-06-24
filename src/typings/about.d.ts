declare namespace About {


    interface AutoCheckNewVersionConfig {
        beta: boolean
        release: boolean
    }

    interface NewVersionInfo {
        releaseTime: string
        type: string
        version: string
        downloadURL: string
        pageUrl: string
    }


}