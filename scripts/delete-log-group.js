const AWS = require('aws-sdk');

// eslint-disable-next-line no-unused-vars
const { argv } = require('yargs/yargs')(process.argv.slice(2))
  .command('region', 'The aws region')
  .command('log-group')
  .command('delete', 'Delete the log group if its found')
  .default('delete', false)
  .boolean('delete')
  .demandOption(['region', 'log-group']);

const cw = new AWS.CloudWatchLogs({ apiVersion: '2014-03-28', region: argv.region });

const checkLogGroup = async (name) => {
  const group = await cw.describeLogGroups({ logGroupNamePrefix: name })
    .promise()
    .catch((err) => { console.error(err); });

  if (group.logGroups.length === 1) {
    console.log(`Log group ${group.logGroups[0].logGroupName} found:\n${JSON.stringify(group.logGroups[0])}`);
    return true;
  } if (group.logGroups.length > 1) {
    console.log(`Multiple log groups found (may be truncated):\n${JSON.stringify(group.logGroups.slice(0, 5))}`);
    return false;
  } if (group.logGroups.length === 0) {
    console.log(`No log groups found:\n${JSON.stringify(group.logGroups)}`);
    return false;
  }

  return false;
};

(async () => {
  const name = argv['log-group'];
  const exists = await checkLogGroup(name)
    .catch((err) => {
      console.error(err, err.error);
      process.exit(1);
    });

  if ((exists) && (argv.delete)) {
    const del = await cw.deleteLogGroup({ logGroupName: name }).promise()
      .catch((err) => {
        console.error(err, err.error);
        process.exit(1);
      });
    console.log(`Log group deleted: \n${JSON.stringify(del)}`);
    process.exit(0);
  }

  process.exit(0);
})();
