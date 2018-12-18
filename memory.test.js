const Memory = require('./memory')

test('get latest notification', async () => {
    const memory = new Memory('Notifications')
    const latest = await memory.getLatest()
    console.log(latest)
})