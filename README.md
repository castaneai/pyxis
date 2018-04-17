pyxis
============

Pixiv notifications to Slack

## Deploy with CircleCI

Set following environment variables on CircleCI project settings

- `GCP_PROJECT_ID` - project id of Google Cloud Platform
- `PYXIS_SESSION` - Pixiv PHPSESSID value
- `PYXIS_SLACK_WEBHOOK_URL` - Slack incoming webhook URL

Then you can deploy with `git push origin master` and access https://mikane-dot-GCP_PROJECT_ID.appspot.com/