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
import type { RpcCreateProblemResponseAuthor } from './RpcCreateProblemResponseAuthor';
import {
    RpcCreateProblemResponseAuthorFromJSON,
    RpcCreateProblemResponseAuthorFromJSONTyped,
    RpcCreateProblemResponseAuthorToJSON,
} from './RpcCreateProblemResponseAuthor';

/**
 * 
 * @export
 * @interface RpcGetUserResponse
 */
export interface RpcGetUserResponse {
    /**
     * 
     * @type {RpcCreateProblemResponseAuthor}
     * @memberof RpcGetUserResponse
     */
    user: RpcCreateProblemResponseAuthor;
}

/**
 * Check if a given object implements the RpcGetUserResponse interface.
 */
export function instanceOfRpcGetUserResponse(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "user" in value;

    return isInstance;
}

export function RpcGetUserResponseFromJSON(json: any): RpcGetUserResponse {
    return RpcGetUserResponseFromJSONTyped(json, false);
}

export function RpcGetUserResponseFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetUserResponse {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'user': RpcCreateProblemResponseAuthorFromJSON(json['User']),
    };
}

export function RpcGetUserResponseToJSON(value?: RpcGetUserResponse | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'User': RpcCreateProblemResponseAuthorToJSON(value.user),
    };
}

