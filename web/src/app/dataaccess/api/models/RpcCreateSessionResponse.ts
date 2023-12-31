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
 * @interface RpcCreateSessionResponse
 */
export interface RpcCreateSessionResponse {
    /**
     * 
     * @type {RpcCreateAccountResponseAccount}
     * @memberof RpcCreateSessionResponse
     */
    account: RpcCreateAccountResponseAccount;
    /**
     * 
     * @type {string}
     * @memberof RpcCreateSessionResponse
     */
    token: string;
}

/**
 * Check if a given object implements the RpcCreateSessionResponse interface.
 */
export function instanceOfRpcCreateSessionResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "account" in value;
    isInstance = isInstance && "token" in value;

    return isInstance;
}

export function RpcCreateSessionResponseFromJSON(json: any): RpcCreateSessionResponse {
    return RpcCreateSessionResponseFromJSONTyped(json, false);
}

export function RpcCreateSessionResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcCreateSessionResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'account': RpcCreateAccountResponseAccountFromJSON(json['Account']),
        'token': json['Token'],
    };
}

export function RpcCreateSessionResponseToJSON(value?: RpcCreateSessionResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'Account': RpcCreateAccountResponseAccountToJSON(value.account),
        'Token': value.token,
    };
}

