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
import type { RpcError } from './RpcError';
import {
    RpcErrorFromJSON,
    RpcErrorFromJSONTyped,
    RpcErrorToJSON,
} from './RpcError';

/**
 * 
 * @export
 * @interface ResponseBodyOfTheCreateTestCaseListMethod
 */
export interface ResponseBodyOfTheCreateTestCaseListMethod {
    /**
     * A String specifying the version of the JSON-RPC protocol. MUST be exactly "2.0".
     * @type {string}
     * @memberof ResponseBodyOfTheCreateTestCaseListMethod
     */
    jsonrpc?: ResponseBodyOfTheCreateTestCaseListMethodJsonrpcEnum;
    /**
     * It MUST be the same as the value of the id member in the Request.
     * @type {string}
     * @memberof ResponseBodyOfTheCreateTestCaseListMethod
     */
    id?: string;
    /**
     * 
     * @type {RpcError}
     * @memberof ResponseBodyOfTheCreateTestCaseListMethod
     */
    error?: RpcError;
    /**
     * 
     * @type {object}
     * @memberof ResponseBodyOfTheCreateTestCaseListMethod
     */
    result?: object;
}


/**
 * @export
 */
export const ResponseBodyOfTheCreateTestCaseListMethodJsonrpcEnum = {
    _20: '2.0'
} as const;
export type ResponseBodyOfTheCreateTestCaseListMethodJsonrpcEnum = typeof ResponseBodyOfTheCreateTestCaseListMethodJsonrpcEnum[keyof typeof ResponseBodyOfTheCreateTestCaseListMethodJsonrpcEnum];


/**
 * Check if a given object implements the ResponseBodyOfTheCreateTestCaseListMethod interface.
 */
export function instanceOfResponseBodyOfTheCreateTestCaseListMethod(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ResponseBodyOfTheCreateTestCaseListMethodFromJSON(json: any): ResponseBodyOfTheCreateTestCaseListMethod {
    return ResponseBodyOfTheCreateTestCaseListMethodFromJSONTyped(json, false);
}

export function ResponseBodyOfTheCreateTestCaseListMethodFromJSONTyped(json: any, ignoreDiscriminator: boolean): ResponseBodyOfTheCreateTestCaseListMethod {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'jsonrpc': !exists(json, 'jsonrpc') ? undefined : json['jsonrpc'],
        'id': !exists(json, 'id') ? undefined : json['id'],
        'error': !exists(json, 'error') ? undefined : RpcErrorFromJSON(json['error']),
        'result': !exists(json, 'result') ? undefined : json['result'],
    };
}

export function ResponseBodyOfTheCreateTestCaseListMethodToJSON(value?: ResponseBodyOfTheCreateTestCaseListMethod | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'jsonrpc': value.jsonrpc,
        'id': value.id,
        'error': RpcErrorToJSON(value.error),
        'result': value.result,
    };
}

