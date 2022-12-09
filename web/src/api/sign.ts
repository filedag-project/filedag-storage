import {Sha256} from '@aws-crypto/sha256-js';
import { SignatureV4} from '@aws-sdk/signature-v4';
import {HttpRequest} from '@aws-sdk/protocol-http';
import { SignModel } from '@/models/SignModel';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from '@/utils/cookies';

const signV4 = async (sign:SignModel) => {
    const AccessKeyId = Cookies.getKey(ACCESS_KEY_ID);
    const SecretAccessKey = Cookies.getKey(SECRET_ACCESS_KEY);
    const SessionToken = Cookies.getKey(SESSION_TOKEN);
    const signer = new SignatureV4({
        service: sign.service, 
        region: sign.region,
        sha256: Sha256,
        applyChecksum: sign.applyChecksum,
        credentials: {
            accessKeyId: sign.accessKeyId??AccessKeyId,
            secretAccessKey: sign.secretAccessKey??SecretAccessKey,
            sessionToken: sign.sessionToken??SessionToken,
        },
        uriEscapePath: true
    });
    const minimalRequest = new HttpRequest({
        method: sign.method,
        protocol: sign.protocol,
        path: sign.path,
        port: 9985,
        query:sign.query??{},
        headers: {
            host: process.env['REACT_APP_HOST']??'',
            'Content-Type': sign.contentType??'',
        },
        hostname: process.env['REACT_APP_HOST_NAME']??'',
        body: sign.body,
    });
    const request = await signer.sign(minimalRequest);
    return request;
};

export default signV4;
