import { Injectable } from '@angular/core';

const KB_IN_BYTE = 1024;
const MB_IN_BYTE = KB_IN_BYTE * 1024;
const GB_IN_BYTE = MB_IN_BYTE * 1024;

@Injectable({
  providedIn: 'root',
})
export class UnitService {
  public timeLimitToValueAndUnit(timeLimitInMillisecond: number): {
    value: number;
    unit: string;
  } {
    if (timeLimitInMillisecond % 1000 === 0) {
      return { value: timeLimitInMillisecond / 1000, unit: 'second' };
    }
    return { value: timeLimitInMillisecond, unit: 'ms' };
  }

  public timeValueAndUnitToLimit(value: number, unit: string): number {
    if (unit === 'second') {
      return value * 1000;
    }
    return value;
  }

  public memoryLimitToValueAndUnit(memoryLimitInByte: number): {
    value: number;
    unit: string;
  } {
    if (memoryLimitInByte % GB_IN_BYTE === 0) {
      return { value: memoryLimitInByte / GB_IN_BYTE, unit: 'GB' };
    }
    if (memoryLimitInByte % MB_IN_BYTE === 0) {
      return { value: memoryLimitInByte / MB_IN_BYTE, unit: 'MB' };
    }
    if (memoryLimitInByte % KB_IN_BYTE === 0) {
      return { value: memoryLimitInByte / KB_IN_BYTE, unit: 'kB' };
    }
    return { value: memoryLimitInByte, unit: 'byte' };
  }

  public memoryValueAndUnitToLimit(value: number, unit: string): number {
    if (unit === 'GB') {
      return value * GB_IN_BYTE;
    }
    if (unit === 'MB') {
      return value * MB_IN_BYTE;
    }
    if (unit === 'kB') {
      return value * KB_IN_BYTE;
    }
    return value;
  }
}
