const http = require('request-promise-native')
const BASE_URL = 'https://www.pixiv.net'

const createSessionCookieJar = (sessionId) => {
    const jar = http.jar();
    jar.setCookie(http.cookie(`PHPSESSID=${sessionId}`), 'http://www.pixiv.net');
    return jar;
};

const getNotifications = async (sessionId) => {
    const options = {jar: createSessionCookieJar(sessionId)}
    const res = JSON.parse(await http.get(`${BASE_URL}/ajax/notification`, options))
    if (res.error) {
        throw new Error(res.message);
    }
    return res.body.items
}

module.exports = {
    getNotifications,
}