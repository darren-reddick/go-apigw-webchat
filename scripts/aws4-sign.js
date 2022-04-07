const AWS = require('aws-sdk');
const aws4 = require('aws4');
const { URL } = require('url');

const { argv } = require('yargs/yargs')(process.argv.slice(2))
  .command('url', 'The url to sign')
  .demandOption(['url'])
  .alias('u', 'url');

let credentials;

if ((process.env.AWS_ACCESS_KEY_ID) && (process.env.AWS_SECRET_ACCESS_KEY)) {
  console.log('AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY detected in the environment - will use these for signing');
  credentials = {
    accessKeyId: process.env.AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
  };
} else if (process.env.AWS_PROFILE) {
  console.log('Found AWS_PROFILE in environment - will use this for signing');
  credentials = {
    accessKeyId: AWS.config.credentials.accessKeyId,
    secretAccessKey: AWS.config.credentials.secretAccessKey,
  };
} else {
  console.error('ERROR: No credentials for signing detected in the environment');
  process.exit(1);
}

const url = new URL(argv.url);

const options = {
  host: url.host,
  path: url.pathname,
  signQuery: true,
};

const sign = aws4.sign(options, credentials);

// eslint-disable-next-line no-console
console.info(`wss://${sign.host}${sign.path}`);
