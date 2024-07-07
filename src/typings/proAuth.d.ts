declare namespace ProAuth{

    interface LoginReq {
        username: string;
        password: string;
    }
    
    interface GetStatusResp {
        proExpiration: string | null; // 时间戳字符串或 null
        loginStatus: boolean;
        username:string
        nickName:string
    }

}