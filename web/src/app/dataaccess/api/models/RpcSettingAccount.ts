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

import type { RpcAccountSetting } from './RpcAccountSetting';
import {
    instanceOfRpcAccountSetting,
    RpcAccountSettingFromJSON,
    RpcAccountSettingFromJSONTyped,
    RpcAccountSettingToJSON,
} from './RpcAccountSetting';

/**
 * @type RpcSettingAccount
 * 
 * @export
 */
export type RpcSettingAccount = RpcAccountSetting;

export function RpcSettingAccountFromJSON(json: any): RpcSettingAccount {
    return RpcSettingAccountFromJSONTyped(json, false);
}

export function RpcSettingAccountFromJSONTyped(json: any, ignoreDiscriminator: boolean): RpcSettingAccount {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return { ...RpcAccountSettingFromJSONTyped(json, true) };
}

export function RpcSettingAccountToJSON(value?: RpcSettingAccount | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }

    if (instanceOfRpcAccountSetting(value)) {
        return RpcAccountSettingToJSON(value as RpcAccountSetting);
    }

    return {};
}

