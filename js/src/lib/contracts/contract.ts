import {SolidityEvent} from './event';
import {SolidityFunction, Handler} from './function';
import {HiveSmartChain} from '../hsc';
import {Function, Event} from 'solc';


export type FunctionOrEvent = Function | Event;

export type ABI = Array<FunctionOrEvent>

export type Bytecode = { bytecode: string, deployedBytecode?: string }


export interface Handlers {
  call?: Handler
  deploy?: Handler
}

const defaultHandlers: Handlers = {
  call: function (result) {
    return result.raw
  },
  deploy: function (result) {
    return result.contractAddress
  }
}

export class Contract {
  abi: ABI;
  address: string;
  code: Bytecode;
  hsc: HiveSmartChain;
  handlers: Handlers;

  _constructor: any;

  constructor(abi: Array<FunctionOrEvent>, code: string | Bytecode, address: string, hsc: HiveSmartChain, handlers?: Handlers) {
    handlers = Object.assign({}, defaultHandlers, handlers);

    this.address = address;
    this.abi = abi;
    this.code = typeof (code) === 'string' ? {bytecode: code} : code;
    this.hsc = hsc;
    this.handlers = handlers;
    addFunctionsToContract(this);
    addEventsToContract(this);
  }

}

const addFunctionsToContract = function (contract: Contract) {
  (contract.abi.filter(json => {
    return (json.type === 'function' || json.type === 'constructor');
  }) as Function[]).forEach(function (json) {
    let {displayName, typeName, call, encode, decode} = SolidityFunction(json, contract.hsc);

    if (json.type === 'constructor') {
      contract._constructor = call.bind(contract, false, contract.handlers.deploy, '');
    } else {
      // bind the function call to the contract, specify if call or transact is desired
      let execute = call.bind(contract, json.constant, contract.handlers.call, null);
      execute.sim = call.bind(contract, true, contract.handlers.call, null);
      // These allow the interface to be used for a generic contract of this type
      execute.at = call.bind(contract, json.constant, contract.handlers.call);
      execute.atSim = call.bind(contract, true, contract.handlers.call);

      execute.encode = encode.bind(contract);
      execute.decode = decode.bind(contract);

      // Attach to the contract object
      if (!contract[displayName]) {
        contract[displayName] = execute;
      }
      contract[displayName][typeName] = execute;
    }
  })

  // Not every abi has a constructor specification.
  // If it doesn't we force a _constructor with null abi
  if (!contract._constructor) {
    const {call} = SolidityFunction(null, contract.hsc);
    contract._constructor = call.bind(contract, false, contract.handlers.deploy, '');
  }
}

const addEventsToContract = function (contract: Contract) {
  (contract.abi.filter(json => {
    return json.type === 'event'
  }) as Event[]).forEach(json => {
    const {displayName, typeName, call} = SolidityEvent(json, contract.hsc);

    const execute = call.bind(contract, null);
    execute.once = call.bind(contract, null);
    execute.at = call.bind(contract);
    if (!contract[displayName]) {
      contract[displayName] = execute;
    }
    contract[displayName][typeName] = call.bind(contract);
  })
}
