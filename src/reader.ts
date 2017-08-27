import {forRead, isEqual, startsWith, toString, toType} from './typedarray';
const MAGIC = Uint8Array.of(0x40, 0x15, 0x00, 0x00, 0x44, 0x43, 0x53, 0x01, 0x01, 0x00, 0x00, 0x00);
const PY_MAGIC = Uint8Array.of(0x9D, 0x01, 0x00, 0x00);
const OFFSET_PINGYIN = 0x1540;
const OFFSET_CHINESE = 0x2628;

export interface ScelInfo {
    name: string;
    type: string;
    description: string;
    example: string;
}

export class ScelReader {
    private ua: Uint8Array;
    private _info?: ScelInfo;
    private _pinyins?: Array<[number, string]>;
    private _words?: Array<[number[], string]>;
    private _pinyin?: Map<number, string>;

    get info() {
        return this._info || (this._info = readInfo(this.ua));
    }

    get pinyins() {
        return this._pinyins || (this._pinyins = readPinyin(this.ua));
    }

    get words() {
        return this._words || (this._words = readWord(this.ua));
    }

    get pinyin() {
        if (this._pinyin != null) {
            return this._pinyin;
        }
        let py: Map<number, string> = this._pinyin = new Map;
        this.pinyins.forEach(([i, v]) => py.set(i, v))
        return py;
    }

    constructor(ua: Uint8Array) {
        this.ua = ua;
    }

    getPinyin(indexes: number[]): (string | undefined)[] {
        return indexes.map(v => this.pinyin.get(v));
    }
}


function isMagicMatch(buf: ArrayBufferView) {
    return isEqual(buf, MAGIC, 0, MAGIC.length)
}

function readString(buf: ArrayBufferView): string {
    let ua = toType(buf, Uint8Array);
    let i = 0;
    for (; i < ua.byteLength; i += 2) {
        if (ua[i] === 0 && ua[i + 1] === 0) {
            break;
        }
    }
    if (i > 0) {
        // console.log('LEN',i,'CON', ua.subarray(0, i).join(','));
        return toString(ua.subarray(0, i), 'utf-16le')
    }
    if (i < 0) {
        return toString(ua, 'utf-16le')
    }
    return '';
}

function readName(buf: Uint8Array): string {
    return readString(buf.subarray(0x130, 0x338));
}
function readType(buf: Uint8Array): string {
    return readString(buf.subarray(0x338, 0x540));
}
function readDescription(buf: Uint8Array): string {
    return readString(buf.subarray(0x540, 0xD40));
}
function readExample(buf: Uint8Array): string {
    return readString(buf.subarray(0xD40, OFFSET_PINGYIN));
}
function readInfo(buf: Uint8Array) {
    return {
        name: readName(buf),
        type: readType(buf),
        description: readDescription(buf),
        example: readExample(buf),
    };
}

function readPinyin(buf: Uint8Array): Array<[number, string]> {
    let ua = buf.subarray(OFFSET_PINGYIN, OFFSET_CHINESE);
    if (!startsWith(ua, PY_MAGIC)) {
        throw new Error('PinYin Magic not match');
    }
    let dvr = forRead(ua);
    dvr.littleEndian = true;
    dvr.position = PY_MAGIC.length;
    let items: Array<[number, string]> = [];
    while (dvr.hasRemaining()) {
        let index = dvr.getUint16();
        let len = dvr.getUint16();
        let s = readString(ua.subarray(dvr.position, dvr.position + len));
        dvr.position += len;
        items.push([index, s]);
    }
    return items;
}

function readWord(buf: Uint8Array): Array<[number[], string]> {
    let ua = buf.subarray(OFFSET_CHINESE);
    let dvr = forRead(ua);
    dvr.littleEndian = true;

    let words: Array<[number[], string]> = [];
    while (dvr.hasRemaining()) {
        // 同音字
        let same = dvr.getUint16();
        let pyLen = dvr.getUint16() / 2;// 2 byte / uint16
        let pys = [];
        for (let i = 0; i < pyLen; i++) {
            pys.push(dvr.getUint16());
        }

        for (let i = 0; i < same; i++) {
            let wordLen = dvr.getUint16();
            let word = readString(ua.subarray(dvr.position, dvr.position + wordLen));
            dvr.position += wordLen;
            let extLen = dvr.getUint16();
            // 目前不知道扩展是什么用
            // ua.slice(dvr.position, dvr.position + extLen);
            dvr.position += extLen;

            words.push([pys, word]);
        }
        if (words.length > 10) {
            break
        }
    }
    return words;
}
