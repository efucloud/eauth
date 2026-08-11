
import CryptoJS from 'crypto-js';
 
//  * @name AES-加密
//  * @param raw 待加密字段
//  * @param AESKey AES Key
//  * @return {string} 返回加密字段
 
export const aesEncrypt = (raw: any, AESKey: string) => {
    const cypherKey = CryptoJS.enc.Utf8.parse(AESKey);
    CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);

    const iv = CryptoJS.SHA256(AESKey).toString();
    const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };
    return CryptoJS.AES.encrypt(raw, cypherKey, cfg).toString();
}

//    * @name AES-解密
//    * @param raw 待解密数据
//    * @param AESKey 解密 key
//    * @returns {string} 返回解密字符串
export const aesDecrypt = (raw: string, AESKey: string) => {
    const cypherKey = CryptoJS.enc.Utf8.parse(AESKey);
    CryptoJS.pad.ZeroPadding.pad(cypherKey, 4);
    const iv = CryptoJS.SHA256(AESKey).toString();
    const cfg = { iv: CryptoJS.enc.Utf8.parse(iv) };

    const decrypt = CryptoJS.AES.decrypt(raw, cypherKey, cfg);
    return CryptoJS.enc.Utf8.stringify(decrypt).toString();
}


export const  encode=(arraybuffer) =>{
    let s = String.fromCharCode.apply(null, new Uint8Array(arraybuffer))
    return window.btoa(s).replace(/\+/g, '-').replace(/\//g, '_');
  };
  
export const decode=(base64string: string) =>{
    let s = window.atob(base64string.replace(/-/g, '+').replace(/_/g, '/'))
    let bytes = Uint8Array.from(s, c => c.charCodeAt(0))
    return bytes.buffer
  };
export const randomString=()=> {    
    var t = "ABCDEFGHJKMNPQRSTWXYZabcdefhijkmnprstwxyz2345678",
    a = t.length,
    n = "";
    for (let i = 0; i < 16; i++) n += t.charAt(Math.floor(Math.random() * a));
    return n
}