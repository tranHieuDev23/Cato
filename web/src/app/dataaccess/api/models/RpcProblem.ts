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
import type { RpcProblemExample } from './RpcProblemExample';
import {
    RpcProblemExampleFromJSON,
    RpcProblemExampleFromJSONTyped,
    RpcProblemExampleToJSON,
} from './RpcProblemExample';

/**
 * 
 * @export
 * @interface RpcProblem
 */
export interface RpcProblem {
    /**
     * 
     * @type {string}
     * @memberof RpcProblem
     */
    uUID: string;
    /**
     * 
     * @type {string}
     * @memberof RpcProblem
     */
    displayName: string;
    /**
     * 
     * @type {RpcCreateAccountResponseAccount}
     * @memberof RpcProblem
     */
    author: RpcCreateAccountResponseAccount;
    /**
     * 
     * @type {string}
     * @memberof RpcProblem
     */
    description: string;
    /**
     * 
     * @type {number}
     * @memberof RpcProblem
     */
    timeLimitInMillisecond: number;
    /**
     * 
     * @type {number}
     * @memberof RpcProblem
     */
    memoryLimitInByte: number;
    /**
     * 
     * @type {Array<RpcProblemExample>}
     * @memberof RpcProblem
     */
    exampleList: Array<RpcProblemExample>;
    /**
     * 
     * @type {number}
     * @memberof RpcProblem
     */
    createdTime: number;
    /**
     * 
     * @type {number}
     * @memberof RpcProblem
     */
    updatedTime: number;
}

/**
 * Check if a given object implements the RpcProblem interface.
 */
export function instanceOfRpcProblem(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "uUID" in value;
    isInstance = isInstance && "displayName" in value;
    isInstance = isInstance && "author" in value;
    isInstance = isInstance && "description" in value;
    isInstance = isInstance && "timeLimitInMillisecond" in value;
    isInstance = isInstance && "memoryLimitInByte" in value;
    isInstance = isInstance && "exampleList" in value;
    isInstance = isInstance && "createdTime" in value;
    isInstance = isInstance && "updatedTime" in value;

    return isInstance;
}

export function RpcProblemFromJSON(json: any): RpcProblem {
    return RpcProblemFromJSONTyped(json, false);
}

export function RpcProblemFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcProblem {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'uUID': json['UUID'],
        'displayName': json['DisplayName'],
        'author': RpcCreateAccountResponseAccountFromJSON(json['Author']),
        'description': json['Description'],
        'timeLimitInMillisecond': json['TimeLimitInMillisecond'],
        'memoryLimitInByte': json['MemoryLimitInByte'],
        'exampleList': ((json['ExampleList'] as Array<any>).map(RpcProblemExampleFromJSON)),
        'createdTime': json['CreatedTime'],
        'updatedTime': json['UpdatedTime'],
    };
}

export function RpcProblemToJSON(value?: RpcProblem | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'UUID': value.uUID,
        'DisplayName': value.displayName,
        'Author': RpcCreateAccountResponseAccountToJSON(value.author),
        'Description': value.description,
        'TimeLimitInMillisecond': value.timeLimitInMillisecond,
        'MemoryLimitInByte': value.memoryLimitInByte,
        'ExampleList': ((value.exampleList as Array<any>).map(RpcProblemExampleToJSON)),
        'CreatedTime': value.createdTime,
        'UpdatedTime': value.updatedTime,
    };
}

