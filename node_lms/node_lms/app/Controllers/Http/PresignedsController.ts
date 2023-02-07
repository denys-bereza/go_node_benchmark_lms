import type { HttpContextContract } from '@ioc:Adonis/Core/HttpContext'

import { getSignedUrl } from "@aws-sdk/s3-request-presigner";
import { S3Client, GetObjectCommand } from "@aws-sdk/client-s3";

const config = {
  credentials: {
    accessKeyId: "",
    secretAccessKey: "",
  },
  region: "eu-central-1",
};

const client = new S3Client(config);

async function getSignedFileUrl(fileName, bucket, expiresIn) {
  // Instantiate the GetObject command,
  // a.k.a. specific the bucket and key
  const command = new GetObjectCommand({
    Bucket: bucket,
    Key: fileName,
  });

  // await the signed URL and return it
  return await getSignedUrl(client, command, { expiresIn });
}

export default class PresignedsController {

public async index(ctx: HttpContextContract) {
  const filename = ctx.params["*"].join("/");
  const link = await getSignedFileUrl(filename, "", 3600)

    ctx.response.redirect(link, false, 301);
  }
}
