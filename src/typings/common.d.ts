declare namespace Common {
    interface ListResponse<T> {	
        list:T
        count:number
    }

    interface ListRequest{	
        limit:number
        page:number
        keyword?:string
    }

    interface InfoBase{	
        createTime?:string
        updateTime?:string
        id?:number
    }

    // 请求-带有弹窗验证数据结构
    interface VerificationRequest{
        codeId?:string
        vCode?:string
    }

    // 响应-带有弹窗验证数据结构
    interface VerificationResponse{
        codeId?:string
        result?:boolean
        message?:string
    }

    interface SortItemRequest{
        id:number
        sort:number
    }
}