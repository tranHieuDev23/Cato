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
 * @interface RpcTestCaseSnippet
 */
export interface RpcTestCaseSnippet {
    /**
     * 
     * @type {string}
     * @memberof RpcTestCaseSnippet
     */
    uUID: string;
    /**
     * 
     * @type {string}
     * @memberof RpcTestCaseSnippet
     */
    input: string;
    /**
     * 
     * @type {string}
     * @memberof RpcTestCaseSnippet
     */
    output: string;
    /**
     * 
     * @type {boolean}
     * @memberof RpcTestCaseSnippet
     */
    isHidden: boolean;
}

/**
 * Check if a given object implements the RpcTestCaseSnippet interface.
 */
export function instanceOfRpcTestCaseSnippet(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "uUID" in value;
    isInstance = isInstance && "input" in value;
    isInstance = isInstance && "output" in value;
    isInstance = isInstance && "isHidden" in value;

    return isInstance;
}

export function RpcTestCaseSnippetFromJSON(json: any): RpcTestCaseSnippet {
    return RpcTestCaseSnippetFromJSONTyped(json, false);
}

export function RpcTestCaseSnippetFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcTestCaseSnippet {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'uUID': json['UUID'],
        'input': json['Input'],
        'output': json['Output'],
        'isHidden': json['IsHidden'],
    };
}

export function RpcTestCaseSnippetToJSON(value?: RpcTestCaseSnippet | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'UUID': value.uUID,
        'Input': value.input,
        'Output': value.output,
        'IsHidden': value.isHidden,
    };
}

