import { IperhApi } from "./ipehrApi";
import { initializeStorageWithDefaults, getStorageData } from './storage';

const interceptable: string[] = [
  'https://sandbox.better.care/studio/api/rest/v1/query',
  'https://sandbox.better.care/ehr/rest/v1/ehr',
  'https://sandbox.better.care/studio/api/rest/v1/composition'
]

const uuidRegexp = /[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}/gm;
const compositionUidRegexp = /[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}::\w*::\d+/;
const EHR_CREATE_URL = 'https://sandbox.better.care/ehr/rest/v1/ehr';
const COMPOSITION_CREATE_URL = 'https://sandbox.better.care/studio/api/rest/v1/composition';
const BETTER_STATE_URL = 'https://sandbox.better.care/#state=';
const BETTER_SILENT_REFRESH_URL = 'https://sandbox.better.care/silent-refresh.html#state=';
//
chrome.runtime.onInstalled.addListener(async () => {
    await initializeStorageWithDefaults({
        url: "https://gateway.ipehr.org/v1",
        username: "p1",
        password: "p1",
        systemId: "sc_bsn_si"
    });


    console.log('Extension successfully installed!');
});

chrome.webRequest.onCompleted.addListener(function (details) {
    const url = details.url;
    const method = details.method;

    if (method === 'GET' && 
        (url.startsWith(BETTER_STATE_URL) || url.startsWith(BETTER_SILENT_REFRESH_URL))
    ) {
        const urlParams = new URLSearchParams(details.url);
        const access_token = urlParams.get('access_token');

        if (access_token) {
            console.log('access token', access_token)
            chrome.storage.sync.set(
                { token: access_token },
                () => {
                    console.log('Better token saved')
                }
            );
        }
    }

    if (method === 'POST' && url.startsWith(EHR_CREATE_URL)) {
        let ehrId = '';
        details.responseHeaders.forEach((h, index) => {
            if (h.name === 'link') {
                ehrId = decodeURIComponent(h.value).match(uuidRegexp)
            }
        })

        if (ehrId === '') {
            console.log(new Error('Could not get the value of ehrId'));
            return;
        }

        createEhr(url, ehrId);
    }

    if (method === 'POST' && url.startsWith(COMPOSITION_CREATE_URL)) {
        const searchParams = new URLSearchParams(url.split("?")[1]);
        const ehrId = searchParams.get('ehrId');

        if (ehrId === '') {
            console.log(new Error('Could not get the value of EhrId'));
            return;
        }

        let compositionUid = '';
        details.responseHeaders.forEach((h, index) => {
            if (h.name === 'link') {
                compositionUid = decodeURIComponent(h.value).match(compositionUidRegexp)
            }
        })

        if (compositionUid === '') {
            console.log(new Error('Could not get the value of composition UID'));
            return;
        }

        createComposition(ehrId, compositionUid);
    }

}, {urls: ['https://sandbox.better.care/*']}, ["responseHeaders"]);


async function ipehrApiInit(): Promise<IperhApi> {
    const data = await getStorageData();

    return new IperhApi({
        apiUrl: data.url,
        username: data.username,
        password: data.password,
        extensionId: '', 
        betterToken: data.token,
        systemId: data.systemId
    })
}

async function createEhr(url: string, ehrId: string) {
    const ipehrApi = await ipehrApiInit(); 
    ipehrApi.createEhr(url, ehrId);

    chrome.storage.sync.set(
        { ehrId: ehrId },
        () => {
            console.log('EhrId saved:', ehrId)
        }
    );
}

async function createComposition(ehrId: string, compositionUid: string) {
    const ipehrApi = await ipehrApiInit(); 
    ipehrApi.createComposition(ehrId, compositionUid);
}
