import { HttpRequest } from "@aws-sdk/protocol-http";
import { S3RequestPresigner } from "@aws-sdk/s3-request-presigner";
import { parseUrl } from "@aws-sdk/url-parser";
import {Sha256} from '@aws-crypto/sha256-js';
import { ACCESS_KEY_ID, Cookies, SECRET_ACCESS_KEY, SESSION_TOKEN } from "@/utils/cookies";
import { PreSignModel } from "@/models/PreSignModel";

const presignV4 = async (sign:PreSignModel) => {
    const AccessKeyId = Cookies.getKey(ACCESS_KEY_ID);
    const SecretAccessKey = Cookies.getKey(SECRET_ACCESS_KEY);
    const SessionToken = Cookies.getKey(SESSION_TOKEN);
    const s3ObjectUrl = parseUrl("http://"+process.env['REACT_APP_HOST'] + sign.path);

    const preSigner = new S3RequestPresigner({
        credentials: {
            accessKeyId: sign.accessKeyId??AccessKeyId,
            secretAccessKey: sign.secretAccessKey??SecretAccessKey,
            sessionToken: sign.sessionToken??SessionToken,
        },
        region: sign.region,
        sha256: Sha256,
    });
    const request = new HttpRequest(s3ObjectUrl);

    const options = {
        expiresIn: sign.expiresIn,
    };

    const url = await preSigner.presign(request,options);
    return url;
};

export default presignV4;
