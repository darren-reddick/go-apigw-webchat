const AWS = require('aws-sdk');
const aws4 = require('aws4');
const { URL } = require('url');

const { argv } = require('yargs/yargs')(process.argv.slice(2))
  .command('url', 'The url to sign')
  .demandOption(['url'])
  .alias('u', 'url');

if (!process.env.AWS_PROFILE) {
  // eslint-disable-next-line no-console
  console.error('ERROR: AWS_PROFILE environment variable not set');
  process.exit(1);
}

const url = new URL(argv.url);

const options = {
  host: url.host,
  path: url.pathname,
  signQuery: true,
};

const credentials = {
  accessKeyId: AWS.config.credentials.accessKeyId,
  secretAccessKey: AWS.config.credentials.secretAccessKey,
};

const sign = aws4.sign(options, credentials);

// eslint-disable-next-line no-console
console.info(`wss://${sign.host}${sign.path}`);
