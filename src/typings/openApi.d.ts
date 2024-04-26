declare namespace OpenApi{

    interface TokenInfo {
        token: string;
    }
    
    interface GetTokenStatusInfoResp extends TokenInfo {
        status: boolean;
    }
    

}