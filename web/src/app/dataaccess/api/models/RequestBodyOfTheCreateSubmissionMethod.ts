/* tslint:disable */
/* eslint-disable */
/**
 * API
 * Generated by genpjrpc: v0.4.0
 *
 * The version of the OpenAPI document: v0.0.0-unknown
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { RpcCreateSubmissionRequest } from './RpcCreateSubmissionRequest';
import {
    RpcCreateSubmissionRequestFromJSON,
    RpcCreateSubmissionRequestFromJSONTyped,
    RpcCreateSubmissionRequestToJSON,
} from './RpcCreateSubmissionRequest';

/**
 * 
 * @export
 * @interface RequestBodyOfTheCreateSubmissionMethod
 */
export interface RequestBodyOfTheCreateSubmissionMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof RequestBodyOfTheCreateSubmissionMethod
     */
    jsonrpc?: RequestBodyOfTheCreateSubmissionMethodJsonrpcEnum;
    /**
     * An identifier established by the Client.
     * @type {string}
     * @memberof RequestBodyOfTheCreateSubmissionMethod
     */
    id?: string;
    /**
     * A String containing the name of the method to be invoked.
     * @type {string}
     * @memberof RequestBodyOfTheCreateSubmissionMethod
     */
    method?: RequestBodyOfTheCreateSubmissionMethodMethodEnum;
    /**
     * 
     * @type {RpcCreateSubmissionRequest}
     * @memberof RequestBodyOfTheCreateSubmissionMethod
     */
    params?: RpcCreateSubmissionRequest;
}


/**
 * @export
 */
export const RequestBodyOfTheCreateSubmissionMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type RequestBodyOfTheCreateSubmissionMethodJsonrpcEnum = typeof RequestBodyOfTheCreateSubmissionMethodJsonrpcEnum[keyof typeof RequestBodyOfTheCreateSubmissionMethodJsonrpcEnum];

/**
 * @export
 */
export const RequestBodyOfTheCreateSubmissionMethodMethodEnum = {
    CreateSubmission: 'create_submission'
} as const;
export type RequestBodyOfTheCreateSubmissionMethodMethodEnum = typeof RequestBodyOfTheCreateSubmissionMethodMethodEnum[keyof typeof RequestBodyOfTheCreateSubmissionMethodMethodEnum];


/**
 * Check if a given object implements the RequestBodyOfTheCreateSubmissionMethod interface.
 */
export function instanceOfRequestBodyOfTheCreateSubmissionMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RequestBodyOfTheCreateSubmissionMethodFromJSON(json: any): RequestBodyOfTheCreateSubmissionMethod {
    return RequestBodyOfTheCreateSubmissionMethodFromJSONTyped(json, false);
}

export function RequestBodyOfTheCreateSubmissionMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestBodyOfTheCreateSubmissionMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'method': !exists(json, 'method') ? undefined : json['method'],
        'params': !exists(json, 'params') ? undefined : RpcCreateSubmissionRequestFromJSON(json['params']),
    };
}

export function RequestBodyOfTheCreateSubmissionMethodToJSON(value?: RequestBodyOfTheCreateSubmissionMethod | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'jsonrpc': value.jsonrpc,
        'id': value.id,
        'method': value.method,
        'params': RpcCreateSubmissionRequestToJSON(value.params),
    };
}

