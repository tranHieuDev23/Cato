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
/**
 * REQUIRED on error. This member MUST NOT exist if there was no error triggered during invocation.
 * @export
 * @interface RpcError
 */
export interface RpcError {
    /**
     * A Number that indicates the error type that occurred.
     * @type {number}
     * @memberof RpcError
     */
    code?: number;
    /**
     * A String providing a short description of the error.
     * @type {string}
     * @memberof RpcError
     */
    message?: string;
    /**
     * A Primitive or Structured value that contains additional information about the error.
     * @type {object}
     * @memberof RpcError
     */
    data?: object;
}

/**
 * Check if a given object implements the RpcError interface.
 */
export function instanceOfRpcError(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function RpcErrorFromJSON(json: any): RpcError {
    return RpcErrorFromJSONTyped(json, false);
}

export function RpcErrorFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcError {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'code': !exists(json, 'code') ? undefined : json['code'],
        'message': !exists(json, 'message') ? undefined : json['message'],
        'data': !exists(json, 'data') ? undefined : json['data'],
    };
}

export function RpcErrorToJSON(value?: RpcError | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'code': value.code,
        'message': value.message,
        'data': value.data,
    };
}

