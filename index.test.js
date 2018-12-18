const index = require('./index')

test('get notifications from pixiv', async () => {
    const sessionId = process.env.PYXIS_SESSION_ID
    if (!sessionId) {
        throw new Error('env: PYXIS_SESSION_ID not found')
    }
    const items = await index.getNotifications(sessionId)
    expect(items).toBeInstanceOf(Array)
})