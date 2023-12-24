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
 * 
 * @export
 * @interface RpcUpdateAccountRequest
 */
export interface RpcUpdateAccountRequest {
    /**
     * 
     * @type {number}
     * @memberof RpcUpdateAccountRequest
     */
    iD: number;
    /**
     * 
     * @type {string}
     * @memberof RpcUpdateAccountRequest
     */
    displayName: string;
    /**
     * 
     * @type {string}
     * @memberof RpcUpdateAccountRequest
     */
    role: string;
    /**
     * 
     * @type {string}
     * @memberof RpcUpdateAccountRequest
     */
    password: string;
}

/**
 * Check if a given object implements the RpcUpdateAccountRequest interface.
 */
export function instanceOfRpcUpdateAccountRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "iD" in value;
    isInstance = isInstance && "displayName" in value;
    isInstance = isInstance && "role" in value;
    isInstance = isInstance && "password" in value;

    return isInstance;
}

export function RpcUpdateAccountRequestFromJSON(json: any): RpcUpdateAccountRequest {
    return RpcUpdateAccountRequestFromJSONTyped(json, false);
}

export function RpcUpdateAccountRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcUpdateAccountRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'iD': json['ID'],
        'displayName': json['DisplayName'],
        'role': json['Role'],
        'password': json['Password'],
    };
}

export function RpcUpdateAccountRequestToJSON(value?: RpcUpdateAccountRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'ID': value.iD,
        'DisplayName': value.displayName,
        'Role': value.role,
        'Password': value.password,
    };
}

