export function ping(url:string) {
    return new Promise((resolve, reject) => {
        var start = new Date().getTime();
        var xhr = new XMLHttpRequest();

        xhr.onreadystatechange = function () {
            if (xhr.readyState === XMLHttpRequest.DONE) {
                var end = new Date().getTime();
                var time = end - start;
                resolve(time);
            }
        };

        xhr.onerror = function () {
            reject(new Error('请求失败'));
        };

        xhr.open("GET", url, true);
        xhr.send();
    });
}
