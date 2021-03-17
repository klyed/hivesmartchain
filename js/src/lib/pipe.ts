import {ServiceError} from "@grpc/grpc-js";
import { ITransactClient } from '../../proto/rpctransact_grpc_pb';
import { CallTx } from '../../proto/payload_pb';
import { TxExecution } from '../../proto/exec_pb';
import { Events } from './events';
import { LogEvent } from '../../proto/exec_pb';
import * as grpc from '@grpc/grpc-js';

export type TxCallback = grpc.requestCallback<TxExecution>;

export class Pipe {
  hsc: ITransactClient;
  events: Events;

  constructor(hsc: ITransactClient, events: Events) {
    this.hsc = hsc;
    this.events = events;
  }

  transact(payload: CallTx, callback: TxCallback) {
    return this.hsc.callTxSync(payload, callback)
  }

  call(payload: CallTx, callback: TxCallback) {
    this.hsc.callTxSim(payload, callback)
  }

  eventSub(accountAddress: string, signature: string, callback: (err: ServiceError, log: LogEvent) => void) {
    return this.events.subContractEvents(accountAddress, signature, callback)
  }
}

