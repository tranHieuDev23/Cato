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

import type { RpcSetting } from './RpcSetting';
import {
    instanceOfRpcSetting,
    RpcSettingFromJSON,
    RpcSettingFromJSONTyped,
    RpcSettingToJSON,
} from './RpcSetting';

/**
 * @type RpcGetServerInfoResponseSetting
 * 
 * @export
 */
export type RpcGetServerInfoResponseSetting = RpcSetting;

export function RpcGetServerInfoResponseSettingFromJSON(json: any): RpcGetServerInfoResponseSetting {
    return RpcGetServerInfoResponseSettingFromJSONTyped(json, false);
}

export function RpcGetServerInfoResponseSettingFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcGetServerInfoResponseSetting {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return { ...RpcSettingFromJSONTyped(json, true) };
}

export function RpcGetServerInfoResponseSettingToJSON(value?: RpcGetServerInfoResponseSetting | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }

    if (instanceOfRpcSetting(value)) {
        return RpcSettingToJSON(value as RpcSetting);
    }

    return {};
}
