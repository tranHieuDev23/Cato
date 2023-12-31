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
import type { RpcCreateAccountResponseAccount } from './RpcCreateAccountResponseAccount';
import {
    RpcCreateAccountResponseAccountFromJSON,
    RpcCreateAccountResponseAccountFromJSONTyped,
    RpcCreateAccountResponseAccountToJSON,
} from './RpcCreateAccountResponseAccount';

/**
 * 
 * @export
 * @interface RpcGetAccountResponse
 */
export interface RpcGetAccountResponse {
    /**
     * 
     * @type {RpcCreateAccountResponseAccount}
     * @memberof RpcGetAccountResponse
     */
    account: RpcCreateAccountResponseAccount;
}

/**
 * Check if a given object implements the RpcGetAccountResponse interface.
 */
export function instanceOfRpcGetAccountResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "account" in value;

    return isInstance;
}

export function RpcGetAccountResponseFromJSON(json: any): RpcGetAccountResponse {
    return RpcGetAccountResponseFromJSONTyped(json, false);
}

export function RpcGetAccountResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetAccountResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'account': RpcCreateAccountResponseAccountFromJSON(json['Account']),
    };
}

export function RpcGetAccountResponseToJSON(value?: RpcGetAccountResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'Account': RpcCreateAccountResponseAccountToJSON(value.account),
    };
}

