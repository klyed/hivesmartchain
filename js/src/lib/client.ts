import { HiveSmartChain } from './hsc'
import * as convert from './utils/convert';
import * as coder from 'ethereumjs-abi';
import { Readable } from 'stream';
import { TxInput, CallTx } from '../../proto/payload_pb'
import { TxExecution } from '../../proto/exec_pb';

export type Interceptor = (result: TxExecution) => Promise<TxExecution>;

export class Client extends HiveSmartChain {
    interceptor: Interceptor;
    
    constructor(url: string, account: string) {
        super(url, account);

        this.interceptor = async (data) => data;
    }

    deploy(msg: CallTx, callback: (err: Error, addr: Uint8Array) => void) {
        this.pipe.transact(msg, (err, exec) => {
            if (err) callback(err, null);
            else if (exec.hasException()) callback(new Error(exec.getException().getException()), null);
            else callback(null, exec.getReceipt().getContractaddress_asU8());
        })
    }

    call(msg: CallTx, callback: (err: Error, exec: Uint8Array) => void) {
        this.pipe.transact(msg, (err, exec) => {
            if (err) callback(err, null);
            else if (exec.hasException()) callback(new Error(exec.getException().getException()), null);
            else this.interceptor(exec).then(exec => callback(null, exec.getResult().getReturn_asU8()));
        })
    }

    callSim(msg: CallTx, callback: (err: Error, exec: Uint8Array) => void) {
        this.pipe.call(msg, (err, exec) => {
            if (err) callback(err, null);
            else if (exec.hasException()) callback(new Error(exec.getException().getException()), null);
            else this.interceptor(exec).then(exec => callback(null, exec.getResult().getReturn_asU8()));
        })
    }

    listen(signature: string, address: string, callback: (err: Error, event: any) => void): Readable {
        return this.events.subContractEvents(address, signature, callback)
    }

    payload(data: string, address?: string): CallTx {
        const input = new TxInput();
        input.setAddress(Buffer.from(this.account, 'hex'));
        input.setAmount(0);
      
        const payload = new CallTx();
        payload.setInput(input);
        if (address) payload.setAddress(Buffer.from(address, 'hex'));
        payload.setGaslimit(1000000);
        payload.setFee(0);
        payload.setData(Buffer.from(data, 'hex'));
      
        return payload
    }

    encode(name: string, inputs: string[], ...args: any[]): string {
        args = convert.hscToAbi(inputs, args);
        return name + convert.bytesTB(coder.rawEncode(inputs, args));
    }

    decode(data: Uint8Array, outputs: string[]): any {
        return convert.abiToHiveSmartChain(outputs, coder.rawDecode(outputs, Buffer.from(data)));
    }
}