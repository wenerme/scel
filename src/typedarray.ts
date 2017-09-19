/// <reference path="../node_modules/@types/text-encoding/index.d.ts" />

import {TextDecoder, TextEncoder} from 'text-encoding';

if (typeof window === 'undefined') {
    require('text-encoding');
}
export interface ArrayType<T extends ArrayBufferView> extends Function {
    new (...args: any[]): T;
}

export function toType<T extends ArrayBufferView>(buf: ArrayBufferView, type: ArrayType<T>): T {
    if (buf.constructor === type) {
        return buf as T;
    }

    return new type(buf.buffer, buf.byteOffset);
}

export function startsWith(buf: ArrayBufferView, prefix: ArrayBufferView) {
    return isEqual(buf, prefix, 0, prefix.byteLength);
}
export function isEqual(buf1: ArrayBufferView, buf2: ArrayBufferView, offset = 0, length = -1) {
    if (buf1 === buf2) {
        return true;
    }
    if (length < 0) {
        length = Math.max(buf1.byteLength, buf2.byteLength)
    }
    if (buf1.byteLength < offset + length || buf2.byteLength < offset + length) {
        return false;
    }

    let a = toType(buf1, Uint8Array);
    let b = toType(buf2, Uint8Array);

    let i = length + offset;
    while (i--) {
        if (a[i] !== b[i]) {
            return false;
        }
    }
    return true;
}

export function fromString(s: string, encoding = 'utf8'): Uint8Array {
    return new TextEncoder(encoding).encode(s);
}

export function toString(ua: Uint8Array, encoding = 'utf8'): string {
    return new TextDecoder(encoding).decode(ua);
}

export function forRead(buf: ArrayBufferView | DataView): DataViewReader {

    function extracted(ctx: { position: number, littleEndian: number }, n: number, byteOffset?: number, littleEndian?: boolean | any): [number, boolean] {
        if (byteOffset == null) {
            byteOffset = ctx.position;
            ctx.position += n;
        }
        if (littleEndian == null) {
            littleEndian = ctx.littleEndian;
        }
        return [byteOffset, littleEndian];
    }

    let dv: DataView;
    if (!(buf instanceof DataView)) {
        dv = new DataView(buf.buffer, buf.byteOffset, buf.byteLength);
    } else {
        dv = buf;
    }

    let proxy = new Proxy(Object.assign(() => 0,
        {
            view: dv,
            position: 0,
            littleEndian: false,
            __m: '',
            hasRemaining(this: DataViewReader){
                return this.position < this.byteLength
            }
        } as any),
        {
            set(target, key, value, receiver){
                if (key in target) {
                    target[key] = value;
                    return true;
                }
                return false;
            },
            get(target, key, receiver){
                if (key in target) {
                    if (typeof target[key] === 'function') {
                        target.__m = key;
                        return receiver;
                    }
                    return target[key];
                }
                switch (key) {
                    case 'buffer':
                    case 'byteLength':
                    case 'byteOffset':
                        return target.view[key];

                }
                target.__m = key;
                return receiver;
            },
            apply(target, self, args){
                // console.log('Apply', target.__m, args);
                if (target.__m in target) {
                    return target[target.__m].apply(self, args)
                }
                let match = target.__m.match(/\d+$/);
                if (match == null) {
                    throw new Error('Invalid delegate');
                }
                let n = +match[0] / 8;
                return ((target.view as any)[target.__m] as () => number).apply(target.view, extracted(target, n, args[0], args[1]));
            }
        });

    return proxy as any as DataViewReader;
}

export interface DataViewReader {
    readonly buffer: ArrayBuffer;
    readonly byteLength: number;
    readonly byteOffset: number;

    position: number;
    littleEndian: boolean;

    getFloat32(byteOffset?: number, littleEndian?: boolean): number;
    getFloat64(byteOffset?: number, littleEndian?: boolean): number;
    getInt8(byteOffset?: number): number;
    getInt16(byteOffset?: number, littleEndian?: boolean): number;
    getInt32(byteOffset?: number, littleEndian?: boolean): number;
    getUint8(byteOffset?: number): number;
    getUint16(byteOffset?: number, littleEndian?: boolean): number;
    getUint32(byteOffset?: number, littleEndian?: boolean): number;

    hasRemaining(): boolean;
}

/**
 // TS2556: Expected 1-2 arguments, but got a minimum of 0.

 export class DataViewReader {

    view: DataView;
    position: number = 0;
    littleEndian: boolean;

    get buffer(): ArrayBuffer {
        return this.view.buffer
    }

    get byteLength(): number {
        return this.view.byteLength;
    }

    get byteOffset(): number {
        return this.view.byteOffset;
    }

    constructor(view: DataView) {
        this.view = view;
    }

    hasRemaining(): boolean {
        return this.position < this.byteLength
    }

    getFloat32(byteOffset?: number, littleEndian?: boolean | undefined): number {
        return this.view.getFloat32(...this.extracted(4, byteOffset, littleEndian));
    }

    getFloat64(byteOffset: number, littleEndian?: boolean | undefined): number {
        return this.view.getFloat64(...this.extracted(8, byteOffset, littleEndian));
    }

    getInt8(byteOffset?: number): number {
        return this.view.getInt8(byteOffset || this.position++);
    }

    getInt16(byteOffset?: number, littleEndian?: boolean | undefined): number {
        return this.view.getInt16(...this.extracted(2, byteOffset, littleEndian));
    }

    getInt32(byteOffset?: number, littleEndian?: boolean | undefined): number {
        return this.view.getInt32(...this.extracted(4, byteOffset, littleEndian));
    }

    getUint8(byteOffset?: number): number {
        return this.view.getUint8(byteOffset || this.position++);
    }

    getUint16(byteOffset?: number, littleEndian?: boolean | undefined): number {
        return this.view.getUint16(...this.extracted(2, byteOffset, littleEndian));
    }

    private extracted(n: number, byteOffset?: number, littleEndian?: boolean | any): [number, boolean] {
        if (byteOffset == null) {
            byteOffset = this.position;
            this.position += n;
        }
        if (littleEndian == null) {
            littleEndian = this.littleEndian;
        }
        return [byteOffset, littleEndian];
    }

    getUint32(byteOffset?: number, littleEndian?: boolean | undefined): number {
        return this.view.getUint32(...this.extracted(4, byteOffset, littleEndian) as any[number]);
    }
}
 */
