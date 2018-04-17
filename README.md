pyxis
============

Pixiv notifications to Slack

## Deploy with CircleCI

Set following environment variables on CircleCI project settings

- `GCP_PROJECT_ID` - project id of Google Cloud Platform
- `GCP_SECRET_KEY` - base64 encoded JSON key which can deploy to Google App Engine
- `PYXIS_SESSION` - Pixiv PHPSESSID value
- `PYXIS_SLACK_WEBHOOK_URL` - Slack incoming webhook URL

Then you can deploy with `git push origin master` and access https://pyxis-dot-GCP_PROJECT_ID.appspot.com/
