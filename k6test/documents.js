import encoding from 'k6/encoding';
import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';
import chai, { describe, expect } from 'https://jslib.k6.io/k6chaijs/4.3.4.3/index.js';
import { ServerUrl, EHRSystemID } from "./consts.js";
import exp from 'constants';


// r.POST("access/document", a.DocAccess.Set)
export function set_document(ctx, document) {
    //TODO - implement
}

// r.GET("access/document/", a.DocAccess.List)
export function get_document_list(ctx) {
    //TODO - implement
}
