const datastore = require('@google-cloud/datastore')();

class Memory {
    constructor(memoryName) {
        this.memoryName = memoryName
    }

    getLatest() {
        const query = datastore.createQuery(this.memoryName)
            .order('notifiedAt', {descending: true})
            .limit(1)
        return datastore.runQuery(query)
    }
}

module.exports = Memory